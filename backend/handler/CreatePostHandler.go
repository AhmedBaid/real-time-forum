package handler

import (
	"encoding/json"
	"net/http"

	"real_time/backend/config"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var post config.Posts
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Error in docoding", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Post created successfully"}
	json.NewEncoder(w).Encode(response)
}
