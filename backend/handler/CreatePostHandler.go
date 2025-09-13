package handler

import (
	"encoding/json"
	"net/http"
	"slices"

	"real_time/backend/config"
	"real_time/backend/helpers"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		config.ResponseJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"message": "Method not allowed. Only POST is accepted.",
			"status":  http.StatusMethodNotAllowed,
		})
		return
	}

	_, session := helpers.SessionChecked(w, r)

	var post struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Categories  []int  `json:"categories"`
	}
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		config.ResponseJSON(w, http.StatusBadRequest, map[string]any{
			"message": "Invalid JSON in request body.",
			"status":  http.StatusBadRequest,
		})
		return
	}

	validCategories := []int{1, 2, 3, 4, 5, 6, 7, 8}
	for _, v := range post.Categories {
		if !slices.Contains(validCategories, v) {
			config.ResponseJSON(w, http.StatusBadRequest, map[string]any{
				"message": "Invalid category ID.",
				"status":  http.StatusBadRequest,
			})
			return
		}
	}

	if post.Title == "" || post.Description == "" || len(post.Categories) == 0 {
		config.ResponseJSON(w, http.StatusBadRequest, map[string]any{
			"message": "Title, description, and at least one category are required.",
			"status":  http.StatusBadRequest,
		})
		return
	}
	if len(post.Title) < 3 || len(post.Title) > 30 ||
		len(post.Description) < 1 || len(post.Description) > 100 {
		config.ResponseJSON(w, http.StatusBadRequest, map[string]any{
			"message": "Title or description length is invalid.",
			"status":  http.StatusBadRequest,
		})
		return
	}

	var userId int
	var username string
	stmt2 := `SELECT username, id FROM users WHERE session = ?`
	err := config.Db.QueryRow(stmt2, session).Scan(&username, &userId)
	if err != nil {
		config.ResponseJSON(w, http.StatusUnauthorized, map[string]any{
			"message": "Unauthorized. Invalid session.",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	stmt := `INSERT INTO posts (title, description, username, userID) VALUES (?, ?, ?, ?)`
	res, err := config.Db.Exec(stmt, post.Title, post.Description, username, userId)
	if err != nil {
		config.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Failed to create post.",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	postID, err := res.LastInsertId()
	if err != nil {
		config.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
			"message": "Failed to retrieve post ID.",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	stmtcat := `INSERT INTO categories_post (categoryID, postID) VALUES (?, ?)`
	for _, v := range post.Categories {
		_, err := config.Db.Exec(stmtcat, v, postID)
		if err != nil {
			config.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
				"message": "Failed to assign category to post.",
				"status":  http.StatusInternalServerError,
			})
			return
		}
	}

	config.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "Post created successfully.",
		"status":  http.StatusOK,
	})
}
