package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"real_time/backend/config"
	"real_time/backend/helpers"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		config.ResponseJSON(w, config.ErrorMethodnotAll.Code, config.ErrorMethodnotAll)
		return
	}

	// Check session
	_, session := helpers.SessionChecked(w, r)

	// Decode JSON
	var post struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Categories  []string    `json:"categories"`
	}
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		fmt.Println(err)
		http.Error(w, "Error in decode", http.StatusBadRequest)
		return
	}

	// Validation
	if post.Title == "" || post.Description == "" || len(post.Categories) == 0 {
		fmt.Println("Title, description, or categories are empty")
		config.ResponseJSON(w, http.StatusBadRequest, map[string]any{
			"message": "Title, description, and categories are required",
			"status":  http.StatusBadRequest,
		})
		return
	}
	if len(post.Title) < 3 || len(post.Title) > 30 ||
		len(post.Description) < 1 || len(post.Description) > 100 {
		config.ResponseJSON(w, http.StatusBadRequest, map[string]any{
			"message": "invalid data length",
			"status":  http.StatusBadRequest,
		})
		return
	}

	// Get user
	var userId int
	var username string
	stmt2 := `SELECT username, id FROM users WHERE session = ?`
	err := config.Db.QueryRow(stmt2, session).Scan(&username, &userId)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Insert post
	stmt := `INSERT INTO posts (title, description, username, userID) VALUES (?, ?, ?, ?)`
	res, err := config.Db.Exec(stmt, post.Title, post.Description, username, userId)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	postID, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "Error getting post ID", http.StatusInternalServerError)
		return
	}

	// Insert categories
	stmtcat := `INSERT INTO categories_post (categoryID, postID) VALUES (?, ?)`
	for _, v := range post.Categories {
		_, err := config.Db.Exec(stmtcat, v, postID)
		if err != nil {
			http.Error(w, "Error inserting category", http.StatusInternalServerError)
			return
		}
	}
	config.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "Post created successfully",
		"status":  http.StatusOK,
	})
}
