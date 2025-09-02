package handler

import (
	"fmt"
	"net/http"

	"real_time/backend/config"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")
	var sessValue string
	if err != nil {
		sessValue = ""
	} else {
		sessValue = session.Value
	}
	var userId int
	query := `select id  from users where session = ?`
	config.Db.QueryRow(query, sessValue).Scan(&userId)

	userQuery := `
    SELECT u.username, u.id, MAX(m.created_at) as last_message_time
    FROM users u
    LEFT JOIN messages m 
        ON u.id = m.sender_id OR u.id = m.receiver_id
    WHERE u.id != ?
    GROUP BY u.id, u.username
`
	rows, err := config.Db.Query(userQuery, userId)
	if err != nil {
		fmt.Println(err)
		config.ResponseJSON(w, config.ErrorInternalServerErr.Code, map[string]any{
			"message": "database Error 1",
			"status":  config.ErrorInternalServerErr.Code,
		})
		return
	}

	defer rows.Close()

	var Users []config.UserStatus
	var user config.UserStatus
	for rows.Next() {
		err := rows.Scan(&user.Username, &user.Id, &user.LastMessageTime)
		if err != nil {
			fmt.Println(err)
			config.ResponseJSON(w, config.ErrorInternalServerErr.Code, map[string]any{
				"message": "database Error 2",
				"status":  config.ErrorInternalServerErr.Code,
			})
			return
		}
		user.Status = "offline"
		Users = append(Users, user)
	}

	config.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "users   successful",
		"status":  http.StatusOK,
		"data":    Users,
	})
}
