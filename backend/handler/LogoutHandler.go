package handler

import (
	"net/http"
	"time"

	"real_time/backend/config"
	"real_time/backend/helpers"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		exist, session := helpers.SessionChecked(w, r)
		if !exist {
			config.ResponseJSON(w, config.ErrorUnauthorized.Code, map[string]any{
				"message": "user is not logged",
				"status":  config.ErrorUnauthorized.Code,
			})
			return
		}
		_, err := config.Db.Exec("Update users set session = ? where session = ?", "Null", session)
		if err != nil {
			http.Error(w, "error in updating session", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour),
			HttpOnly: true,
			Path:     "/",
		})
	} else {
		
	}
}
