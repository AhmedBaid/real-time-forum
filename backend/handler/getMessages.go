package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"real_time/backend/config"
	"real_time/backend/helpers"
)

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	recieverIdStr := r.URL.Query().Get("reciever")
	recieverId, _ := strconv.Atoi(recieverIdStr)
	fmt.Println(recieverId)
	_, session := helpers.SessionChecked(w, r)
	var senderId int
	config.Db.QueryRow("SELECT id FROM users WHERE session=?", session).Scan(&senderId)

	rows, err := config.Db.Query(`
        SELECT id, sender_id, receiver_id, message, created_at 
        FROM messages 
        WHERE (sender_id=? AND receiver_id=?) 
        ORDER BY created_at ASC`,
		senderId, recieverId)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	messages := []config.Messages{}
	for rows.Next() {
		var m config.Messages
		var t string
		if err := rows.Scan(&m.Id, &m.Sender, &m.Reciever, &m.Message, &t); err != nil {
			continue
		}
		m.Time, _ = time.Parse(time.RFC3339, t)
		messages = append(messages, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
