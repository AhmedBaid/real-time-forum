package handler

import (
	"log"
	"net/http"
	"real_time/backend/config"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	users     = make(map[*websocket.Conn]bool) // all connected clients
	broadcast = make(chan []config.Messages)   // broadcast channel
)

func handleBroadcast() {
	for {
		updateOnline := <-broadcast
		for client := range users {
			err := client.WriteJSON(updateOnline)
			if err != nil {
				log.Printf("WebSocket error: %v", err)
				client.Close()
				delete(users, client)
			}
		}
	}
}

func reader(conn *websocket.Conn) {
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Client disconnected:", err)
			delete(users, conn)
			conn.Close()
			break
		}
	}
}

func ws(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	users[conn] = true

	if err := conn.WriteJSON(config.Messages); err != nil {
		log.Println("Initial send error:", err)
	}
	go reader(conn)
}
