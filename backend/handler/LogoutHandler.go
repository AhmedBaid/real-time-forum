package handler

import (
	"net/http"
	"time"

	"real_time/backend/config"
	"real_time/backend/helpers"
)

var LoggedOut bool = false

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		config.ResponseJSON(w, http.StatusUnauthorized, map[string]any{
			"message": "Unauthorized. Invalid session.",
			"status":  http.StatusUnauthorized,
		})
		return
	}
	// Check if user session exists
	exist, session := helpers.SessionChecked(w, r)
	if !exist {
		// Clear session cookie if not logged in
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour),
			HttpOnly: true,
			Path:     "/",
		})
		config.ResponseJSON(w, http.StatusUnauthorized, map[string]any{
			"message": "User is not logged in.",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	// Remove session from database
	_, err := config.Db.Exec("UPDATE users SET session = NULL WHERE session = ?", session)
	if err != nil {
		config.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Internal server error while logging out.",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	// Respond with success
	config.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "User logged out successfully.",
		"status":  http.StatusOK,
	})
	LoggedOut = true
}
