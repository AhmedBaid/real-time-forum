package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"real_time/backend/config"
	"real_time/backend/helpers"
)

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {

		config.ResponseJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"message": "method not allowd ",
			"status":  http.StatusMethodNotAllowed,
		})
		return

	}

	_, session := helpers.SessionChecked(w, r)

	stmt2 := `select username from users where session = ?`
	query := config.Db.QueryRow(stmt2, session)

	var username string
	errr := query.Scan(&username)
	if errr != nil {
		config.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Erron in database 1",
			"status":  http.StatusInternalServerError,
		})
		return

	}

	var comment config.Comments
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		fmt.Println(err)
		config.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Erron in decoding",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	fmt.Println(comment.PostID)
	stmt3 := `select id from posts where id = ?`
	query3 := config.Db.QueryRow(stmt3, comment.PostID)

	var postID2 int
	errrr := query3.Scan(&postID2)
	if errrr != nil {
		config.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Erron in database 2",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	stmt := `insert into comments (postID, comment, username ) values(?, ? ,?)`
	_, errrrr := config.Db.Exec(stmt, postID2, comment.Comment, username)
	if errrrr != nil {
		config.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Erron in database 3",
			"status":  http.StatusInternalServerError,
		})
		return
	}
	config.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "comments created  successful",
		"status":  http.StatusOK,
		"data":    comment,
	})
}
