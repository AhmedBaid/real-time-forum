package middleware

import (
	"net/http"

	"real_time/backend/config"
	"real_time/backend/helpers"
)

func IsLogged(w http.ResponseWriter, r *http.Request) {
	
	exist, _ := helpers.SessionChecked(w, r)
	if !exist {
		config.ResponseJSON(w, config.ErrorUnauthorized.Code, map[string]any{
			"message": "user is not logged",
			"status":  config.ErrorUnauthorized.Code,
		})
		return
	}
	
	config.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "user is logged",
	})
}

func Authorisation(HandlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil || cookie.Value == "" {
			config.ResponseJSON(w, config.ErrorUnauthorized.Code, map[string]any{
				"message": "you are not authorized",
				"status":  config.ErrorUnauthorized.Code,
			})
			return
		} else {
			// Check if the session is valid
			stmt := "SELECT id FROM users WHERE session = ?"
			var userID int
			err = config.Db.QueryRow(stmt, cookie.Value).Scan(&userID)
			if err != nil {
				config.ResponseJSON(w, config.ErrorUnauthorized.Code, map[string]any{
					"message": "you are not authorized",
					"status":  config.ErrorUnauthorized.Code,
				})
				return
			}
		}
		HandlerFunc(w, r)
	}
}
