package handler

import (
	"net/http"
	"regexp"

	"real_time/backend/config"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// check the session
	if r.Method != http.MethodPost {
		config.ResponseJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"message": "Method not allowrd",
			"status":  http.StatusMethodNotAllowed,
		})
		return
	}
	Username := r.FormValue("Username")
	FirstName := r.FormValue("FirstName")
	LastName := r.FormValue("LastName")
	Email := r.FormValue("Email")
	Password := r.FormValue("Password")
	Age := r.FormValue("Age")
	Gender := r.Form.Get("Gender")
	errMsg, ok := Isvalid(Username, FirstName, LastName, Email, Password, Age, Gender)
	if !ok {
		config.ResponseJSON(w, http.StatusBadRequest, map[string]any{
			"message": errMsg,
			"status":  http.StatusBadRequest,
		})
		return
	}
	

	config.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "Registration successful",
		"status":  http.StatusOK,
	})
}

func Isvalid(username, firstname, lastname, email, password, age, gender string) (string, bool) {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if username == "" || firstname == "" || lastname == "" || email == "" || password == "" || age == "" || gender == "" {
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
