package handler

import (
	"database/sql"
	"encoding/json"
	"html"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"real_time/backend/config"
	"real_time/backend/helpers"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var (
	users   = make(map[int][]*websocket.Conn)
	usersMu sync.Mutex
)

// addUserConn adds a connection for a user (thread-safe).
func addUserConn(userID int, conn *websocket.Conn) {
	usersMu.Lock()
	defer usersMu.Unlock()
	users[userID] = append(users[userID], conn)
}

// removeUserConn removes one connection for a user. If no connections left,
// it deletes the user entry and broadcasts offline status to friends/all users.
func removeUserConn(userID int, conn *websocket.Conn) {
	usersMu.Lock()
	defer usersMu.Unlock()

	conns, exists := users[userID]
	if !exists {
		return
	}
	for i, c := range conns {
		if c == conn {
			// remove that connection
			users[userID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
	if len(users[userID]) == 0 {
		delete(users, userID)
		// broadcast offline synchronously to all connections
		msg := map[string]interface{}{
			"type":   "offline",
			"userId": userID,
			"time":   time.Now().Format(time.RFC3339),
		}
		// send to everyone
		for _, conns := range users {
			for _, c := range conns {
				// ignore write errors here
				_ = c.WriteJSON(msg)
			}
		}
	}
}

// writeToConnSafe writes JSON to a single connection and returns write error.
func writeToConnSafe(c *websocket.Conn, v interface{}) error {
	if c == nil {
		return nil
	}
	if err := c.WriteJSON(v); err != nil {
		// on error, try to close
		_ = c.Close()
		return err
	}
	return nil
}

// sendToUser sends data to all connections of a given userID.
// If a connection write fails it removes that connection.
func sendToUser(userID int, data map[string]interface{}) {
	usersMu.Lock()
	conns := users[userID]
	usersMu.Unlock()

	for _, c := range conns {
		if err := writeToConnSafe(c, data); err != nil {
			// remove this failed connection
			removeUserConn(userID, c)
		}
	}
}

// broadcastToAll sends data to every connected user (synchronously).
func broadcastToAll(data map[string]interface{}) {
	usersMu.Lock()
	all := make([]*websocket.Conn, 0)
	for _, conns := range users {
		all = append(all, conns...)
	}
	usersMu.Unlock()

	for _, c := range all {
		if err := writeToConnSafe(c, data); err != nil {
			// best-effort removal: find which user had this conn and remove
			// (acquire lock inside removeUserConn)
			// We can't easily know owner here; just attempt remove by scanning users map:
			usersMu.Lock()
			for uid, conns := range users {
				for i, cc := range conns {
					if cc == c {
						users[uid] = append(conns[:i], conns[i+1:]...)
						break
					}
				}
				if len(users[uid]) == 0 {
					delete(users, uid)
					// also notify others about offline
					offlineMsg := map[string]interface{}{
						"type":   "offline",
						"userId": uid,
						"time":   time.Now().Format(time.RFC3339),
					}
					// send notification to remaining conns (best-effort)
					for _, remaining := range users {
						for _, rc := range remaining {
							_ = rc.WriteJSON(offlineMsg)
						}
					}
				}
			}
			usersMu.Unlock()
		}
	}
}

// sendUnreadMessages sends unread DB messages to the provided conn for userID.
func sendUnreadMessages(userID int, conn *websocket.Conn, db *sql.DB) {
	if db == nil || conn == nil {
		return
	}
	rows, err := db.Query(`
		SELECT m.id, m.sender_id, m.message, m.created_at, u.username
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.receiver_id = ? AND m.is_read = FALSE
		ORDER BY m.created_at ASC`, userID)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var msgID, senderID int
		var message, createdAt, senderUsername string
		if err := rows.Scan(&msgID, &senderID, &message, &createdAt, &senderUsername); err != nil {
			continue
		}
		message = html.EscapeString(message)
		senderUsername = html.EscapeString(senderUsername)
		data := map[string]interface{}{
			"type":           "message",
			"id":             msgID,
			"sender":         senderID,
			"receiver":       userID,
			"message":        message,
			"time":           createdAt,
			"senderUsername": senderUsername,
		}
		if err := writeToConnSafe(conn, data); err != nil {
			// if writing unread fails, stop trying further (conn probably dead)
			return
		}
		_, _ = db.Exec("UPDATE messages SET is_read = TRUE WHERE id = ?", msgID)
	}
}

// handleConnection reads & handles messages from a connection (synchronously).
// It does NOT spawn a goroutine — caller should call this inline.
func handleConnection(userID int, conn *websocket.Conn, db *sql.DB) {
	defer func() {
		removeUserConn(userID, conn)
		_ = conn.Close()
		// if no connections left for this user, an offline broadcast already handled in removeUserConn
	}()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			// read error: client disconnected or protocol error
			return
		}
		var msg map[string]interface{}
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}
		msgType, ok := msg["type"].(string)
		if !ok {
			continue
		}

		switch msgType {
		case "message":
			// receiver expected numeric (float64 from JSON)
			receiverF, ok := msg["receiver"].(float64)
			if !ok {
				continue
			}
			receiver := int(receiverF)
			content, ok := msg["message"].(string)
			if !ok {
				continue
			}
			// sanitize
			content = html.EscapeString(content)

			// get sender username
			var senderUsername string
			_ = db.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&senderUsername)

			// insert into DB
			res, err := db.Exec(`
				INSERT INTO messages (sender_id, receiver_id, message, is_read)
				VALUES (?, ?, ?, FALSE)`, userID, receiver, content)
			if err != nil {
				continue
			}
			msgID64, _ := res.LastInsertId()
			msgID := int(msgID64)

			out := map[string]interface{}{
				"type":           "message",
				"id":             msgID,
				"sender":         userID,
				"receiver":       receiver,
				"message":        content,
				"time":           time.Now().Format(time.RFC3339),
				"senderUsername": senderUsername,
			}

			// send to receiver(s)
			sendToUser(receiver, out)

			// send back to sender connections (confirm)
			sendToUser(userID, map[string]interface{}{
				"type":           "message_sent",
				"id":             msgID,
				"sender":         userID,
				"receiver":       receiver,
				"message":        content,
				"time":           time.Now().Format(time.RFC3339),
				"senderUsername": senderUsername,
			})

			// send a notification to receiver (synchronously)
			notification := map[string]interface{}{
				"type":     "notification",
				"receiver": receiver,
				"from":     userID,
				"message":  "New message",
				"time":     time.Now().Format(time.RFC3339),
			}
			sendToUser(receiver, notification)

		case "typing", "stopTyping":
			receiverF, ok := msg["receiver"].(float64)
			if !ok {
				continue
			}
			receiver := int(receiverF)
			typingMsg := map[string]interface{}{
				"type":           msgType,
				"senderId":       userID,
				"senderUsername": msg["senderUsername"],
				"time":           time.Now().Format(time.RFC3339),
			}
			sendToUser(receiver, typingMsg)

		default:
			// unknown type - optionally broadcast to all
			// broadcast synchronously
			broadcastToAll(msg)
		}
	}
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	_, session := helpers.SessionChecked(w, r)
	var userID int
	err := config.Db.QueryRow("SELECT id FROM users WHERE session = ?", session).Scan(&userID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// register connection
	addUserConn(userID, conn)

	// broadcast online to everyone (synchronously)
	onlineMsg := map[string]interface{}{
		"type":   "online",
		"userId": userID,
		"time":   time.Now().Format(time.RFC3339),
	}
	broadcastToAll(onlineMsg)

	// send current online list to this connection
	usersMu.Lock()
	onlineUsers := []int{}
	for id := range users {
		onlineUsers = append(onlineUsers, id)
	}
	usersMu.Unlock()
	_ = conn.WriteJSON(map[string]interface{}{
		"type":  "online_list",
		"users": onlineUsers,
		"time":  time.Now().Format(time.RFC3339),
	})

	// send unread messages to this connection
	sendUnreadMessages(userID, conn, config.Db)

	// handle this connection in the same goroutine (no go routine)
	handleConnection(userID, conn, config.Db)
}

// MarkReadHandler - unchanged logic, marks message read
func MarkReadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	msgID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	config.Db.Exec("UPDATE messages SET is_read = TRUE WHERE id = ?", msgID)
	w.WriteHeader(http.StatusOK)
}

// GetMessagesHandler - unchanged core logic
func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	receiverIDStr := r.URL.Query().Get("receiver")
	offsetStr := r.URL.Query().Get("offset")

	offsetId, err := strconv.Atoi(offsetStr)
	if err != nil {
		http.Error(w, "Invalid offset", http.StatusBadRequest)
		return
	}

	receiverID, err := strconv.Atoi(receiverIDStr)
	if err != nil {
		http.Error(w, "Invalid receiver ID", http.StatusBadRequest)
		return
	}

	_, session := helpers.SessionChecked(w, r)
	var senderID int
	err = config.Db.QueryRow("SELECT id FROM users WHERE session = ?", session).Scan(&senderID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := config.Db.Query(`
		SELECT m.id, m.sender_id, m.receiver_id, m.message, m.created_at, u.username
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE (m.sender_id = ? AND m.receiver_id = ?) OR (m.sender_id = ? AND m.receiver_id = ?)
		ORDER BY m.id Desc LIMIT 10 OFFSET ?`, senderID, receiverID, receiverID, senderID, offsetId)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	messages := []map[string]interface{}{}
	for rows.Next() {
		var m config.Messages
		var senderUsername string
		var createdAt string
		if err := rows.Scan(&m.Id, &m.Sender, &m.Reciever, &m.Message, &createdAt, &senderUsername); err != nil {
			log.Printf("Error scanning message: %v", err)
			continue
		}
		m.Message = html.EscapeString(m.Message)
		senderUsername = html.EscapeString(senderUsername)

		messages = append(messages, map[string]interface{}{
			"id":             m.Id,
			"sender":         m.Sender,
			"receiver":       m.Reciever,
			"message":        m.Message,
			"time":           createdAt,
			"senderUsername": senderUsername,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// MessageHandler - keep existing HTTP message send route but send synchronously to connections
func MessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		config.ResponseJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"message": "method not allowed",
			"status":  http.StatusMethodNotAllowed,
		})
		return
	}

	_, session := helpers.SessionChecked(w, r)
	var senderId int
	err := config.Db.QueryRow("SELECT id FROM users WHERE session=?", session).Scan(&senderId)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var message config.Messages
	err = json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		config.ResponseJSON(w, http.StatusBadRequest, map[string]any{
			"message": "Invalid request",
		})
		return
	}

	message.Sender = senderId

	stmt := `INSERT INTO messages (sender_id, receiver_id, message) VALUES (?, ?, ?)`
	res, err := config.Db.Exec(stmt, message.Sender, message.Reciever, message.Message)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	messageId, _ := res.LastInsertId()
	message.Id = int(messageId)

	var username string
	_ = config.Db.QueryRow("SELECT username FROM users WHERE id=?", senderId).Scan(&username)

	out := map[string]any{
		"type":           "message",
		"id":             message.Id,
		"sender":         message.Sender,
		"reciever":       message.Reciever,
		"message":        message.Message,
		"time":           time.Now().Format(time.RFC3339),
		"senderUsername": username,
	}

	// send directly to recipient and sender synchronously
	sendToUser(message.Reciever, out)
	sendToUser(message.Sender, map[string]any{
		"type":           "message_sent",
		"id":             message.Id,
		"sender":         message.Sender,
		"reciever":       message.Reciever,
		"message":        message.Message,
		"time":           time.Now().Format(time.RFC3339),
		"senderUsername": username,
	})

	// notification
	notif := map[string]any{
		"type":     "notification",
		"reciever": message.Reciever,
		"from":     message.Sender,
		"message":  "new Message",
		"time":     time.Now().Format(time.RFC3339),
	}
	sendToUser(message.Reciever, notif)

	config.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "Message sent",
		"data":    out,
	})
}
