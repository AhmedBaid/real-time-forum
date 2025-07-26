package handler

import (
	"net/http"
	"os"

	"real_time/backend/config"
)

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		config.ResponseJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"message": "Method not allowrd",
			"status":  http.StatusMethodNotAllowed,
		})
		return
	}
	file, err := os.Stat(r.URL.Path[1:])
	if err != nil || file.IsDir() {
		config.ResponseJSON(w, http.StatusNotFound, map[string]any{
			"message": "Page not found",
			"status":  http.StatusNotFound,
		})
		return
	}

	http.ServeFile(w, r, r.URL.Path[1:])
}
