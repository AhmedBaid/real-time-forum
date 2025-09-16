package handler

import (
	"net/http"

	"real_time/backend/config"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	// Allow only GET requests
	if r.Method != http.MethodGet {
		config.ResponseJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"message": "Method not allowed. Only GET is permitted.",
			"status":  http.StatusMethodNotAllowed,
		})
		return
	}

	// Retrieve session cookie if present
	session, err := r.Cookie("session")
	var sessValue string
	if err != nil {
		sessValue = ""
	} else {
		sessValue = session.Value
	}

	// Get user ID from session
	var userId int
	query := `SELECT id FROM users WHERE session = ?`
	config.Db.QueryRow(query, sessValue).Scan(&userId)

	// Query to get users and their last message time
	userQuery := `
        SELECT 
    u.username, 
    u.id, 
    COALESCE(MAX(m.created_at), '') as last_message_time
FROM users u
LEFT JOIN messages m 
    ON ( (u.id = m.sender_id AND m.receiver_id = ?) 
      OR (u.id = m.receiver_id AND m.sender_id = ?) )
WHERE u.id != ?
GROUP BY u.id, u.username

    `
	rows, err := config.Db.Query(userQuery, userId ,userId, userId)
	if err != nil {
		config.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Internal server error while retrieving users.",
			"status":  http.StatusInternalServerError,
		})
		return
	}
	defer rows.Close()

	var Users []config.UserStatus
	var user config.UserStatus
	for rows.Next() {
		err := rows.Scan(&user.Username, &user.Id, &user.LastMessageTime)
		if err != nil {
			config.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
				"message": "Internal server error while processing users.",
				"status":  http.StatusInternalServerError,
			})
			return
		}
		user.Status = "offline"
		Users = append(Users, user)
	}

	config.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "Users retrieved successfully.",
		"status":  http.StatusOK,
		"data":    Users,
	})
}
