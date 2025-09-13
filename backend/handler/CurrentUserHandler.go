package handler

import (
	"encoding/json"
	"net/http"

	"real_time/backend/config"
	"real_time/backend/helpers"
)

func CurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	_, session := helpers.SessionChecked(w, r)
	var userID int
	var username string
	err := config.Db.QueryRow("SELECT id,  username FROM users WHERE session = ?", session).Scan(&userID, &username)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"userId": userID, "username": username})
}
