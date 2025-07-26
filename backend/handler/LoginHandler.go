package handler

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"net/http"

	"real_time/backend/config"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		config.ResponseJSON(w, config.ErrorMethodnotAll.Code, config.ErrorMethodnotAll)
		var user config.Users
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "error in decoder", http.StatusInternalServerError)
			return
		}
		query1 := `SELECT password FROM users WHERE username = ? or email=?`
		var password string
		err = config.Db.QueryRow(query1, user.Username, user.Username).Scan(&password)
		if err != nil {
			if err == sql.ErrNoRows {
				config.ResponseJSON(w, config.ErrorNotFound.Code, map[string]any{
					"message": "User not found",
					"status":  config.ErrorNotFound.Code,
				})
			} else {
				http.Error(w, "error in query", http.StatusInternalServerError)
			}
			return
		}
		if bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password)) != nil {
			config.ResponseJSON(w, config.ErrorUnauthorized.Code, map[string]any{
				"message": "Invalid password",
				"status":  config.ErrorUnauthorized.Code,
			})
			return
		}
		sessionID := uuid.New().String()
		query2 := `UPDATE users SET session = ? WHERE username = ? or email = ?`
		_, err = config.Db.Exec(query2, sessionID, user.Username)
		if err != nil {
			http.Error(w, "error in updating session", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    sessionID,
			HttpOnly: true,
			Path:     "/",
			MaxAge:   3600,
		})
		config.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "Login successful",
			"status":  http.StatusOK,
			"data":    user,
		})
	} else {
		template, err := template.ParseFiles("frontend/main.html")
		if err != nil {
			http.Error(w, "error in parsing template", http.StatusInternalServerError)
			return
		}
		template.Execute(w, nil)
	}
}
