package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"real_time/backend/config"

	"github.com/google/uuid"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// check the session
	if r.Method == http.MethodPost {
		var user config.Users
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			http.Error(w, "error in decoder", http.StatusInternalServerError)
			return
		}

		errMsg, ok := Isvalid(user.Username, user.Email, user.FirstName, user.LastName, user.Password, user.Gender, user.Age)
		if !ok {
			config.ResponseJSON(w, http.StatusBadRequest, map[string]any{
				"message": errMsg,
				"status":  http.StatusBadRequest,
			})
			return
		}
		var userId int
		query1 := `select id from users where username = ? or email = ?`
		err = config.Db.QueryRow(query1, user.Username, user.Email).Scan(&userId)
		if err != sql.ErrNoRows {
			config.ResponseJSON(w, http.StatusBadRequest, map[string]any{
				"message": "Username or email already exists",
				"status":  http.StatusBadRequest,
			})
			return
		}

		session := uuid.New().String()
		query2 := `INSERT INTO users (username, firstname, lastname,email, password,gender, age,session) VALUES (?, ?, ?, ?, ?, ?, ?,?)`
		_, err = config.Db.Exec(query2, user.Username, user.FirstName, user.LastName, user.Email, user.Password, user.Gender, user.Age, session)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "error in query", http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    session,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   3600,
		})
		config.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "Registration successful",
			"status":  http.StatusOK,
		})
	} else {
		http.ServeFile(w, r, "frontend/main.html")
	}
}

func Isvalid(username, email, firstName, lastName, password, gender string, age int) (string, bool) {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if username == "" || firstName == "" || lastName == "" || email == "" || password == "" || age == 0 || gender == "" {
		return "All fields are required", false
	} else if len(email) < 10 || len(email) > 50 {
		return "Email must be between 10 and 50 characters", false
	} else if match, _ := regexp.MatchString(emailRegex, email); !match {
		return "Invalid email format", false
	} else if len(username) < 3 || len(username) > 15 {
		return "Username must be between 3 and 15 characters", false
	} else if len(password) < 6 || len(password) > 15 {
		return "Password must be between 6 and 15 characters", false
	}
	return "", true
}
