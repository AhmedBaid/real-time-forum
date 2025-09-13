package handler

import (
	"html"
	"net/http"

	"real_time/backend/config"
)

func GetComments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		config.ResponseJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"message": "Method not allowed. Only GET is permitted.",
			"status":  http.StatusMethodNotAllowed,
		})
		return
	}

	postid := r.URL.Query().Get("id")
	session, err := r.Cookie("session")
	var sessValue string
	if err != nil {
		sessValue = ""
	} else {
		sessValue = session.Value
	}

	var userId int
	config.Db.QueryRow(`SELECT id FROM users WHERE session = ?`, sessValue).Scan(&userId)

	stmtComments := `
    SELECT 
        c.id,
        COALESCE(c.comment, '') AS comment,
        c.time AS time,
        COALESCE(c.username, '') AS username,
        c.postID, 
        COUNT(*) AS totalcomments 
    FROM comments c
    INNER JOIN posts p ON p.id = c.postID
    WHERE c.postID = ?
    GROUP BY c.id
    ORDER BY c.time DESC;
    `

	rows, err := config.Db.Query(stmtComments, postid)
	if err != nil {
		config.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Internal server error while retrieving comments.",
			"status":  http.StatusInternalServerError,
		})
		return
	}
	defer rows.Close()

	var comments []config.Comments
	for rows.Next() {
		var comment config.Comments
		if err := rows.Scan(&comment.Id, &comment.Comment, &comment.Time, &comment.Username, &comment.PostID, &comment.TotalComments); err != nil {
			config.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
				"message": "Internal server error while processing comments.",
				"status":  http.StatusInternalServerError,
			})
			return
		}
		comment.Comment = html.EscapeString(comment.Comment)
		comment.Username = html.EscapeString(comment.Username)
		comments = append(comments, comment)
	}

	config.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "Comments retrieved successfully.",
		"status":  http.StatusOK,
		"data":    comments,
	})
}
