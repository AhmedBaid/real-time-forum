package handler

import (
	"fmt"
	"html"
	"net/http"

	"real_time/backend/config"
	"real_time/backend/helpers"
)

func GetPosts(w http.ResponseWriter, r *http.Request) {
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
	var sess string
	query := `SELECT id, session FROM users WHERE session = ?`
	config.Db.QueryRow(query, sessValue).Scan(&userId, &sess)
	sessValue = sess

	// Fetch categories map
	categorMap, errcat := helpers.FetchCategories()
	if errcat != nil {
		fmt.Println("FetchCategories error:", errcat)
		config.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Internal server error while fetching categories.",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	// Query to get posts with like/dislike counts and user reaction
	stmt := `
		SELECT 
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
		config.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Internal server error while retrieving posts.",
			"status":  http.StatusInternalServerError,
		})
		return
	}
	defer rows.Close()

	var posts []config.Posts
	var post config.Posts

	// Query to get total comments for a post
	totalcmnts := `
		SELECT COUNT(*) as totalcomments FROM comments c
		INNER JOIN posts p ON p.id = c.postID
		WHERE c.postID = ?
	`

	for rows.Next() {
		var totalLikes, totalDislikes, userReaction int
		err = rows.Scan(&post.Id, &post.Username, &post.Title, &post.Description, &post.Time, &totalLikes, &totalDislikes, &userReaction)
		if err != nil {
			config.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
				"message": "Internal server error while processing posts.",
				"status":  http.StatusInternalServerError,
			})
			return
		}

		// Escape HTML for safety
		post.Username = html.EscapeString(post.Username)
		post.Title = html.EscapeString(post.Title)
		post.Description = html.EscapeString(post.Description)

		// Get total comments for this post
		var totalComments int
		config.Db.QueryRow(totalcmnts, post.Id).Scan(&totalComments)

		// Assign post details
		post.Categories = categorMap[post.Id]
		post.TotalLikes = totalLikes
		post.TotalDislikes = totalDislikes
		post.TotalComments = totalComments
		post.UserReactionPosts = userReaction

		posts = append(posts, post)
	}

	// Prepare response data
	responseData := struct {
		Session string
		Posts   []config.Posts
	}{
		Session: sessValue,
		Posts:   posts,
	}

	// Send successful response
	config.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "Posts retrieved successfully.",
		"status":  http.StatusOK,
		"data":    responseData,
	})
}
