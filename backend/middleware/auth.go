package middleware

import (
	"net/http"

	"real_time/backend/config"
	"real_time/backend/helpers"
)

func IsLogged(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		config.ResponseJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"message": "Method not allowrd",
			"status":  http.StatusMethodNotAllowed,
		})
		return
	}
	exist, session := helpers.SessionChecked(w, r)
	if !exist {
		config.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "user is not logged",
			"status":  http.StatusInternalServerError,
		})
		return
	}
	var UserId int
	var Username string
	query := `select id,username from users where session = ?`
	err := config.Db.QueryRow(query, session).Scan(&UserId, &Username)
	if err != nil {
		config.ResponseJSON(w, http.StatusUnauthorized, map[string]any{
			"message": "you are not authorized",
			"status":  http.StatusInternalServerError,
		})
		return
	}
	config.ResponseJSON(w, http.StatusOK, map[string]any{
		"message":  "user is logged",
		"status":   http.StatusOK,
		"id":       UserId,
		"username": Username,
	})
}
