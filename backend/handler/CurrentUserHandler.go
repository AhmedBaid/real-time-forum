package handler

import (
	"net/http"

	"real_time/backend/config"
	"real_time/backend/helpers"
)

func CurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	_, session := helpers.SessionChecked(w, r)
	var userID int
	var username string
	err := config.Db.QueryRow("SELECT id, username FROM users WHERE session = ?", session).Scan(&userID, &username)
	if err != nil {
		config.ResponseJSON(w, http.StatusUnauthorized, map[string]any{
			"message": "Unauthorized. Invalid session.",
			"status":  http.StatusUnauthorized,
		})
		return
	}
	config.ResponseJSON(w, http.StatusOK, map[string]any{
		"userId":   userID,
		"username": username,
		"status":   http.StatusOK,
	})
}
