package handler

import (
	"encoding/json"
	"net/http"

	"real_time/backend/config"
	"real_time/backend/helpers"
)

func UnreadMessagesHandler(w http.ResponseWriter, r *http.Request) {
	_, session := helpers.SessionChecked(w, r)
	var userID int
	err := config.Db.QueryRow("SELECT id FROM users WHERE session = ?", session).Scan(&userID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, _ := config.Db.Query(`
    SELECT sender_id, COUNT(*) as count
    FROM messages
    WHERE receiver_id = ? AND is_read = FALSE
    GROUP BY sender_id`, userID)
	defer rows.Close()

	notifications := []map[string]interface{}{}
	for rows.Next() {
		var senderID, count int
		rows.Scan(&senderID, &count)
		notifications = append(notifications, map[string]interface{}{
			"sender": senderID,
			"count":  count,
		})
	}
	json.NewEncoder(w).Encode(notifications)
}
