// handler/websocket.go
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
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, 
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
	log.Printf("Added connection for user %d, total connections: %d", userID, len(users[userID]))
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
		log.Printf("User %d is offline, no remaining connections", userID)
	}
}

func sendUnreadMessages(userID int, conn *websocket.Conn, db *sql.DB) {
	if db == nil {
		log.Println("Error: Database connection is nil in sendUnreadMessages")
		return
	}
	rows, err := db.Query(`
        SELECT m.id, m.sender_id, m.message, m.created_at, u.username
        FROM messages m
        JOIN users u ON m.sender_id = u.id
        WHERE m.receiver_id = ? AND m.is_read = FALSE
        ORDER BY m.created_at ASC`, userID)
	if err != nil {
		log.Printf("Error fetching unread messages: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var msgID, senderID int
		var message, createdAt, senderUsername string
		if err := rows.Scan(&msgID, &senderID, &message, &createdAt, &senderUsername); err != nil {
			log.Printf("Error scanning unread message: %v", err)
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
		if err := conn.WriteJSON(data); err != nil {
			log.Printf("Error sending unread message: %v", err)
			continue
		}
		_, err = db.Exec("UPDATE messages SET is_read = TRUE WHERE id = ?", msgID)
		if err != nil {
			log.Printf("Error updating is_read: %v", err)
		}
	}
}

func HandleBroadcast(db *sql.DB) {
	if db == nil {
		log.Fatal("Error: Database connection is nil in HandleBroadcast")
	}
	for {
		data := <-broadcast

		msgType, ok := data["type"].(string)
		if !ok {
			log.Println("⚠️ Invalid broadcast type:", data)
			continue
		}

		usersMu.Lock()
		if receiver, ok := data["receiver"].(int); ok && msgType == "message" {
			if conns, ok := users[receiver]; ok {
				for _, conn := range conns {
					err := conn.WriteJSON(data)
					if err != nil {
						log.Printf("WebSocket send error for user %d: %v", receiver, err)
						conn.Close()
						removeUserConn(receiver, conn)
					} else if msgID, ok := data["id"].(int); ok {
						_, err = db.Exec("UPDATE messages SET is_read = TRUE WHERE id = ?", msgID)
						if err != nil {
							log.Printf("Error updating is_read: %v", err)
						}
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
			for userID, conns := range users {
				for _, conn := range conns {
					err := conn.WriteJSON(data)
					if err != nil {
						log.Printf("WebSocket send error for user %d: %v", userID, err)
						conn.Close()
						removeUserConn(userID, conn)
					}
				}
			}
		}
		usersMu.Unlock()
		log.Printf("📤 Sent [%s]: %+v\n", msgType, data)
	}
}

func reader(userID int, conn *websocket.Conn, db *sql.DB) {
	defer func() {
		removeUserConn(userID, conn)
		conn.Close()
	}()

	conn.SetReadLimit(512)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Ping error for user %d: %v", userID, err)
				return
			}
		default:
			_, data, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Client disconnected: %d, %v", userID, err)
				return
			}

			var msg map[string]interface{}
			if err := json.Unmarshal(data, &msg); err != nil {
				log.Printf("Invalid JSON from user %d: %v", userID, err)
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

				if db == nil {
					log.Println("Error: Database connection is nil in reader")
					continue
				}

				var senderUsername string
				err = db.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&senderUsername)
				if err != nil {
					log.Printf("Error fetching sender username: %v", err)
					continue
				}

				res, err := db.Exec(`
                    INSERT INTO messages (sender_id, receiver_id, message, is_read)
                    VALUES (?, ?, ?, FALSE)`, userID, int(receiver), content)
				if err != nil {
					log.Printf("Database error: %v", err)
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
		log.Println("Upgrade error:", err)
		return
	}

	if config.Db == nil {
		log.Println("Error: Database connection is nil in WsHandler")
		conn.Close()
		return
	}

	addUserConn(userID, conn)
	log.Println("User connected:", userID)

	broadcast <- map[string]interface{}{
		"type":   "online",
		"userId": userID,
		"time":   time.Now().Format(time.RFC3339),
	}

	usersMu.Lock()
	onlineUsers := []int{}
	for id := range users {
		if id != userID {
			onlineUsers = append(onlineUsers, id)
		}
	}
	usersMu.Unlock()

	err = conn.WriteJSON(map[string]interface{}{
		"type":  "online_list",
		"users": onlineUsers,
		"time":  time.Now().Format(time.RFC3339),
	})
	if err != nil {
		log.Println("Error sending online list:", err)
		conn.Close()
		return
	}

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

	if config.Db == nil {
		http.Error(w, "Database connection is nil", http.StatusInternalServerError)
		return
	}

	_, err = config.Db.Exec("UPDATE messages SET is_read = TRUE WHERE id = ?", msgID)
	if err != nil {
		log.Printf("Error updating is_read: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	receiverIDStr := r.URL.Query().Get("receiver")
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

	if config.Db == nil {
		http.Error(w, "Database connection is nil", http.StatusInternalServerError)
		return
	}

	rows, err := config.Db.Query(`
        SELECT m.id, m.sender_id, m.receiver_id, m.message, m.created_at, u.username
        FROM messages m
        JOIN users u ON m.sender_id = u.id
        WHERE (m.sender_id = ? AND m.receiver_id = ?) OR (m.sender_id = ? AND m.receiver_id = ?)
        ORDER BY m.created_at ASC`, senderID, receiverID, receiverID, senderID)
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	messages := []map[string]interface{}{}
	for rows.Next() {
		var m config.Messages
		var senderUsername string
		var t string
		if err := rows.Scan(&m.Id, &m.Sender, &m.Reciever, &m.Message, &t, &senderUsername); err != nil {
			log.Printf("Error scanning message: %v", err)
			continue
		}
		m.Time, _ = time.Parse(time.RFC3339, t)
		messages = append(messages, map[string]interface{}{
			"id":             m.Id,
			"sender":         m.Sender,
			"receiver":       m.Reciever,
			"message":        m.Message,
			"time":           t,
			"senderUsername": senderUsername,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
