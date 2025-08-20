package handler

import (
	"log"
	"net/http"
	"real_time/backend/config"
	"real_time/backend/helpers"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var users = make(map[int]*websocket.Conn)

var broadcast = make(chan map[string]any)

func HandleBroadcast() {
	for {
		data := <-broadcast

		msgType, ok := data["type"].(string)
		if !ok {
			log.Println("⚠️ Broadcast  type:", data)
			continue
		}

		if reciever, ok := data["reciever"].(int); ok {
			if conn, ok := users[reciever]; ok {
				err := conn.WriteJSON(data)
				if err != nil {
					log.Printf("WebSocket send error: %v", err)
					conn.Close()
					delete(users, reciever)
				}
			}
		} else {
			for id, conn := range users {
				err := conn.WriteJSON(data)
				if err != nil {
					log.Printf("WebSocket send error: %v", err)
					conn.Close()
					delete(users, id)
				}
			}
		}

		log.Printf("📤 Sent [%s]: %+v\n", msgType, data)
	}
}

func reader(userId int, conn *websocket.Conn) {
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Client disconnected:", userId, err)
			delete(users, userId)
			conn.Close()

			broadcast <- map[string]any{
				"type":   "offline",
				"userId": userId,
				"time":   time.Now(),
			}
			break
		}
	}
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	_, session := helpers.SessionChecked(w, r)
	var userId int
	config.Db.QueryRow("SELECT id FROM users WHERE session=?", session).Scan(&userId)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	users[userId] = conn
	log.Println("User connected:", userId)

	// ⬇️ خبر الجميع أن هذا user راه online
	broadcast <- map[string]any{
		"type":   "online",
		"userId": userId,
		"time":   time.Now(),
	}

	// ⬇️ رجع للـ user الجديد لائحة الناس لي already online
	onlineUsers := []int{}
	for id := range users {
		if id != userId {
			onlineUsers = append(onlineUsers, id)
		}
	}

	err = conn.WriteJSON(map[string]any{
		"type":   "online_list",
		"users":  onlineUsers,
		"time":   time.Now(),
	})
	if err != nil {
		log.Println("Error sending online list:", err)
	}

	// Start reader
	go reader(userId, conn)
}

