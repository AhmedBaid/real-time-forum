package handler

import (
	"fmt"
	"html"
	"net/http"

	"real_time/backend/config"
	"real_time/backend/helpers"
)

func GetPosts(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session")
	var sessValue string
	if err != nil {
		sessValue = ""
	} else {
		sessValue = session.Value
	}
	query := `select id ,  session from users where session = ?`
	var userId int
	sess := ""

	config.Db.QueryRow(query, sessValue).Scan(&userId, &sess)

	sessValue = sess
	// get comments
	categorMap, errcat := helpers.FetchCategories()

	if  errcat != nil {
		fmt.Println("err1efez", err)
		fmt.Println("err2fefe", errcat)

		config.ResponseJSON(w, config.ErrorInternalServerErr.Code, map[string]any{
			"message": "server Error  1",
			"status":  config.ErrorInternalServerErr.Code,
		})
		return
	}

	// !  get posts
	stmt := `SELECT 
				p.id, 
				p.username, 
				p.title, 
				p.description, 
				p.time, 
				COUNT(CASE WHEN l.value = 1 THEN 1 END) AS total_likes, 
				COUNT(CASE WHEN l.value = -1 THEN 1 END) AS total_dislikes,
				COALESCE((
					SELECT value FROM likes WHERE postID = p.id AND userID = ?
				), 0) AS user_reaction_pub
				FROM posts p
				LEFT JOIN likes l ON p.id = l.postID
				GROUP BY p.id
				ORDER BY p.time DESC;	
	`



	rows, err := config.Db.Query(stmt, userId)
	if err != nil {
		config.ResponseJSON(w, config.ErrorInternalServerErr.Code, map[string]any{
			"message": "server Error 2 ",
			"status":  config.ErrorInternalServerErr.Code,
		})
		return
	}
	var posts []config.Posts
	var post config.Posts

	 var totalcmnts = `SELECT  count(*) as totalcomments FROM comments c 
	 INNER JOIN posts p  on p.id = c.postID
	 WHERE c.postID = ?
	`
	var totalLikes, totalDislikes, user_reaction_pub int
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&post.Id, &post.Username, &post.Title, &post.Description, &post.Time, &totalLikes, &totalDislikes, &user_reaction_pub)
		if err != nil {
			config.ResponseJSON(w, config.ErrorInternalServerErr.Code, map[string]any{
				"message": "server Error 3",
				"status":  config.ErrorInternalServerErr.Code,
			})
			return
		}
		var Total int 
post.Username =  html.EscapeString(post.Username)
post.Title  =  html.EscapeString(post.Title)
post.Description =  html.EscapeString(post.Description)

		config.Db.QueryRow(totalcmnts, post.Id).Scan(&Total)
		post.Categories = categorMap[post.Id]
		post.TotalLikes = totalLikes
		post.TotalDislikes = totalDislikes
		post.TotalComments = Total
		post.UserReactionPosts = user_reaction_pub
		posts = append(posts, post)

	}

	// !  end get posts

	variables := struct {
		Session string
		Posts   []config.Posts
	}{
		Session: sessValue,
		Posts:   posts,
	}

	config.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "posts successful",
		"status":  http.StatusOK,
		"data":    variables,
	})
}
