package handler

import (
	"fmt"
	"html"
	"net/http"

	"real_time/backend/config"
)

func GetComments(w http.ResponseWriter, r *http.Request) {
	postid := r.URL.Query().Get("id")
	fmt.Println(postid)
	session, err := r.Cookie("session")
	var sessValue string
	if err != nil {
		sessValue = ""
	} else {
		sessValue = session.Value
	}

	query := `select id from users where session = ?`
	var userId int
	config.Db.QueryRow(query, sessValue).Scan(&userId)

	//!  get comments

	stmtcommnts := `
	SELECT 
    c.id,
    COALESCE(c.comment, '') AS comment,
    c.time AS time,
    COALESCE(c.username, '') AS username,
    c.postID , 
    count(*) as totalcomments 
	FROM comments c
	INNER JOIN posts p  on p.id = c.postID
	 WHERE c.postID = ?
	GROUP BY c.id
	ORDER BY c.time DESC;
	`

	rows2, err2 := config.Db.Query(stmtcommnts, postid)
	if err2 != nil {
		config.ResponseJSON(w, config.ErrorInternalServerErr.Code, map[string]any{
			"message": "server Error ",
			"status":  config.ErrorInternalServerErr.Code,
		})
		return
	}

	var comments []config.Comments
	defer rows2.Close()
	for rows2.Next() {
		var comment config.Comments
		err2 = rows2.Scan(&comment.Id, &comment.Comment, &comment.Time, &comment.Username, &comment.PostID, &comment.TotalComments)
		if err2 != nil {
			config.ResponseJSON(w, config.ErrorInternalServerErr.Code, map[string]any{
				"message": "server Error ",
				"status":  config.ErrorInternalServerErr.Code,
			})
			return
		}

		comment.Comment =  html.EscapeString(comment.Comment)
		comment.Username =  html.EscapeString(comment.Username)
		
		comments = append(comments, comment)
	}

	//  !  end get comments

	// !  add the communts to   map
	config.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "posts successful",
		"status":  http.StatusOK,
		"data":    comments,
	})

}
