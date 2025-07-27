package handler

import (
	"net/http"
	// "real_time/backend/helpers"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/main.html")
}
