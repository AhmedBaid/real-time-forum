package handler

import (
	"encoding/json"
	"net/http"
	"slices"
	"time"

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

	// Valid categories
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

	// Get user info
	var userId int
	var username string
	err := config.Db.QueryRow(`SELECT username, id FROM users WHERE session = ?`, session).Scan(&username, &userId)
	if err != nil {
		config.ResponseJSON(w, http.StatusUnauthorized, map[string]any{
			"message": "Unauthorized. Invalid session.",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	// Insert post
	res, err := config.Db.Exec(`INSERT INTO posts (title, description, username, userID) VALUES (?, ?, ?, ?)`,
		post.Title, post.Description, username, userId)
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

	// Assign categories
	catNames := []string{}
	for _, catID := range post.Categories {
		_, err := config.Db.Exec(`INSERT INTO categories_post (categoryID, postID) VALUES (?, ?)`, catID, postID)
		if err != nil {
			config.ResponseJSON(w, http.StatusInternalServerError, map[string]any{
				"message": "Failed to assign category to post.",
				"status":  http.StatusInternalServerError,
			})
			return
		}

		// Get category name
		var catName string
		err = config.Db.QueryRow(`SELECT name FROM categories WHERE id = ?`, catID).Scan(&catName)
		if err == nil {
			catNames = append(catNames, catName)
		}
	}

	config.ResponseJSON(w, http.StatusOK, map[string]any{
		"message": "Post created successfully.",
		"status":  http.StatusOK,
		"data": map[string]any{
			"id":                postID,
			"title":             post.Title,
			"description":       post.Description,
			"username":          username,
			"time":              time.Now(),
			"categories":        catNames,
			"totalLikes":        0,
			"totalDislikes":     0,
			"totalComments":     0,
			"userReactionPosts": 0,
		},
	})
}
