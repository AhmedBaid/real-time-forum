package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"real_time/backend/config"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.ServeFile(w, r, "frontend/main.html")
		return
	}

	// Decode user credentials from request body
	var user config.Users
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		config.ResponseJSON(w, http.StatusBadRequest, map[string]any{
			"message": "Invalid request body.",
			"status":  http.StatusBadRequest,
		})
		return
	}

	// Check for required fields
	if user.Username == "" || user.Password == "" {
		config.ResponseJSON(w, http.StatusBadRequest, map[string]any{
			"message": "Username and password are required.",
			"status":  http.StatusBadRequest,
		})
		return
	}

	// Retrieve hashed password from database
	var password string
	query := `SELECT password FROM users WHERE username = ? OR email = ?`
	err := config.Db.QueryRow(query, user.Username, user.Username).Scan(&password)
	if err != nil {
		if err == sql.ErrNoRows {
			config.ResponseJSON(w, http.StatusNotFound, map[string]any{
				"message": "Incorrect Username or Password",
				"status":  http.StatusNotFound,
			})
		} else {
			config.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
				"message": "Internal server error while retrieving user.",
				"status":  http.StatusInternalServerError,
			})
		}
		return
	}

	// Compare provided password with stored hash
	if bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password)) != nil {
		config.ResponseJSON(w, http.StatusUnauthorized, map[string]any{
			"message": "Incorrect Username or Password",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	// Generate new session ID and update user session in database
	sessionID := uuid.New().String()
	updateQuery := `UPDATE users SET session = ? WHERE username = ? OR email = ?`
	_, err = config.Db.Exec(updateQuery, sessionID, user.Username, user.Username)
	if err != nil {
		config.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Internal server error while updating session.",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   3600,
	})

	// Respond with success
	config.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "Login successful.",
		"status":  http.StatusOK,
		"data":    user,
	})
}
