package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"real_time/backend/config"
	"real_time/backend/helpers"
)

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		config.ResponseJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"message": "method not allowd ",
			"status":  http.StatusMethodNotAllowed,
		})
		return
	}

	_, session := helpers.SessionChecked(w, r)

	stmt2 := `select username , id  from users where session = ?`
	query := config.Db.QueryRow(stmt2, session)

	var username string
	var userId int
	errr := query.Scan(&username, &userId)
	if errr != nil {
		config.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Erron in database 1",
			"status":  http.StatusInternalServerError,
		})
		return

	}

	var message config.Messages
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		fmt.Println(err)
		config.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Erron in decoding",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	stmt := `INSERT INTO users (sender_id, receiver_id, message) VALUES (?, ?, ?, ?)`
	res, err := config.Db.Exec(stmt, message.Sender, message.Reciever, message.Message)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	messageId, _ := res.LastInsertId()
	message.Id =  int(messageId)
	config.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "comments created  successful",
		"status":  http.StatusOK,
		"data":    message,
	})
}
