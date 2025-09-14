package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"real_time/backend/config"
	"real_time/backend/helpers"
)

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		config.ResponseJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"message": "Method not allowed. Use POST.",
			"status":  http.StatusMethodNotAllowed,
		})
		return
	}

	_, session := helpers.SessionChecked(w, r)
	stmt2 := `select username from users where session = ?`
	query := config.Db.QueryRow(stmt2, session)

	var username string
	if err := query.Scan(&username); err != nil {
		config.ResponseJSON(w, http.StatusUnauthorized, map[string]any{
			"message": "Unauthorized. Invalid session.",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	var comment config.Comments
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		config.ResponseJSON(w, http.StatusBadRequest, map[string]any{
			"message": "Invalid request body.",
			"status":  http.StatusBadRequest,
		})
		return
	}

	stmt3 := `select id from posts where id = ?`
	query3 := config.Db.QueryRow(stmt3, comment.PostID)

	var postID2 int
	if err := query3.Scan(&postID2); err != nil {
		config.ResponseJSON(w, http.StatusNotFound, map[string]any{
			"message": "Post not found.",
			"status":  http.StatusNotFound,
		})
		return
	}

	stmt := `insert into comments (postID, comment, username) values (?, ?, ?)`
	res, err := config.Db.Exec(stmt, postID2, comment.Comment, username)
	if err != nil {
		config.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Failed to create comment.",
			"status":  http.StatusInternalServerError,
		})
		return
	}
	commentID, _ := res.LastInsertId()

	comment.Username = username
	comment.Id = int(commentID)
	comment.Time = time.Now()

	config.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "Comment created successfully.",
		"status":  http.StatusOK,
		"data":    comment,
	})
}
