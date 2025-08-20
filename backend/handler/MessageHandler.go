package handler

import (
	"encoding/json"
	"net/http"

	"real_time/backend/config"
	"real_time/backend/helpers"
)

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
	config.Db.QueryRow("SELECT username FROM users WHERE id=?", senderId).Scan(&username)

	out := map[string]any{
		"type":           "message",
		"id":             message.Id,
		"sender":         message.Sender,
		"reciever":       message.Reciever,
		"message":        message.Message,
		"time":           message.Time, 
		"senderUsername": username,
	}

	broadcast <- out

	broadcast <- map[string]any{
		"type":     "notification",
		"reciever": message.Reciever,
		"from":     message.Sender,
		"message":  "new Message",
	}

	config.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "Message sent",
		"data":    out,
	})
}
