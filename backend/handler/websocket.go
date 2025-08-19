package handler

import (
	"log"
	"net/http"
	"real_time/backend/config"
	"real_time/backend/helpers"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var users = make(map[int]*websocket.Conn)

var broadcast = make(chan config.Messages)

func HandleBroadcast() {
	for {
		msg := <-broadcast
		if conn, ok := users[msg.Reciever]; ok {
			err := conn.WriteJSON(msg)
			if err != nil {
				log.Printf("WebSocket send error: %v", err)
				conn.Close()
				delete(users, msg.Reciever)
			}
		}
	}
}

func reader(userId int, conn *websocket.Conn) {
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Client disconnected:", userId, err)
			delete(users, userId)
			conn.Close()
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

	go reader(userId, conn)
}
