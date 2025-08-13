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
	query := `select id ,  session , username  from users where session = ?`
	var userId int
	var username string

	sess := ""

	config.Db.QueryRow(query, sessValue).Scan(&userId, &sess, &username)

	userQuery := ` select username from users where id != ?`
	rows, err := config.Db.Query(userQuery, userId)
	if err != nil {
		fmt.Println(err)
		config.ResponseJSON(w, config.ErrorInternalServerErr.Code, map[string]any{
			"message": "database  Error   1 ",
			"status":  config.ErrorInternalServerErr.Code,
		})
		return
	}

	defer rows.Close()

	var Users []string
	for rows.Next() {
		var user string
		err := rows.Scan(&user)
		if err != nil {
		fmt.Println(err)

			config.ResponseJSON(w, config.ErrorInternalServerErr.Code, map[string]any{
				"message": "database  Error  2 ",
				"status":  config.ErrorInternalServerErr.Code,
			})
			return
		}
		Users = append(Users, user)
	}

	config.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "users   successful",
		"status":  http.StatusOK,
		"data":    Users,
	})
}
