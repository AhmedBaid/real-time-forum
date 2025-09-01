package handler

import (
	"database/sql"
	"encoding/json"
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
	users     = make(map[int][]*websocket.Conn)
	usersMu   sync.Mutex
	broadcast = make(chan map[string]interface{}, 256)
)

func addUserConn(userID int, conn *websocket.Conn) {
	usersMu.Lock()
	defer usersMu.Unlock()
	users[userID] = append(users[userID], conn)
}

func removeUserConn(userID int, conn *websocket.Conn) {
	usersMu.Lock()
	defer usersMu.Unlock()
	conns, exists := users[userID]
	if !exists {
		return
	}
	for i, c := range conns {
		if c == conn {
			users[userID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
	if len(users[userID]) == 0 {
		delete(users, userID)
		broadcast <- map[string]interface{}{
			"type":   "offline",
			"userId": userID,
			"time":   time.Now().Format(time.RFC3339),
		}
	}
}

func sendUnreadMessages(userID int, conn *websocket.Conn, db *sql.DB) {
	if db == nil {
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
		data := map[string]interface{}{
			"type":           "message",
			"id":             msgID,
			"sender":         senderID,
			"receiver":       userID,
			"message":        message,
			"time":           createdAt,
			"senderUsername": senderUsername,
		}
		conn.WriteJSON(data)
		db.Exec("UPDATE messages SET is_read = TRUE WHERE id = ?", msgID)
	}
}

func HandleBroadcast(db *sql.DB) {
	if db == nil {
		log.Fatal("Database connection is nil in HandleBroadcast")
	}
	for {
		data := <-broadcast
		msgType, ok := data["type"].(string)
		if !ok {
			continue
		}
		usersMu.Lock()
		if receiver, ok := data["receiver"].(int); ok && msgType == "message" {
			if conns, ok := users[receiver]; ok {
				for _, conn := range conns {
					conn.WriteJSON(data)
					if msgID, ok := data["id"].(int); ok {
						db.Exec("UPDATE messages SET is_read = TRUE WHERE id = ?", msgID)
					}
				}
			}
			if sender, ok := data["sender"].(int); ok {
				broadcast <- map[string]interface{}{
					"type":     "notification",
					"receiver": receiver,
					"from":     sender,
					"message":  "New message",
					"time":     time.Now().Format(time.RFC3339),
				}
			}
		} else {
			for _, conns := range users {
				for _, conn := range conns {
					conn.WriteJSON(data)
				}
			}
		}
		usersMu.Unlock()
	}
}

func reader(userID int, conn *websocket.Conn, db *sql.DB) {
	defer func() {
		removeUserConn(userID, conn)
		conn.Close()

		if LoggedOut {
			usersMu.Lock()
			_, exists := users[userID]
			if !exists {
				usersMu.Unlock()
				return
			}

			delete(users, userID)

			broadcast <- map[string]interface{}{
				"type":   "offline",
				"userId": userID,
				"time":   time.Now().Format(time.RFC3339),
			}

			usersMu.Unlock()
		}
	}()
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
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
		if msgType == "message" {
			receiver, ok := msg["receiver"].(float64)
			if !ok {
				continue
			}
			content, ok := msg["message"].(string)
			if !ok {
				continue
			}
			var senderUsername string
			db.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&senderUsername)
			res, err := db.Exec(`
				INSERT INTO messages (sender_id, receiver_id, message, is_read)
				VALUES (?, ?, ?, FALSE)`, userID, int(receiver), content)
			if err != nil {
				continue
			}
			msgID, _ := res.LastInsertId()
			broadcast <- map[string]interface{}{
				"type":           "message",
				"id":             int(msgID),
				"sender":         userID,
				"receiver":       int(receiver),
				"message":        content,
				"time":           time.Now().Format(time.RFC3339),
				"senderUsername": senderUsername,
			}
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

	addUserConn(userID, conn)
	broadcast <- map[string]interface{}{
		"type":   "online",
		"userId": userID,
		"time":   time.Now().Format(time.RFC3339),
	}
	usersMu.Lock()
	onlineUsers := []int{}
	for id := range users {
		onlineUsers = append(onlineUsers, id)
	}
	usersMu.Unlock()
	conn.WriteJSON(map[string]interface{}{
		"type":  "online_list",
		"users": onlineUsers,
		"time":  time.Now().Format(time.RFC3339),
	})
	sendUnreadMessages(userID, conn, config.Db)
	go reader(userID, conn, config.Db)
}

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
		ORDER BY m.created_at Desc LIMIT 10 OFFSET ?`, senderID, receiverID, receiverID, senderID, offsetId)
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
