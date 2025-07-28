package handler

import (
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/main.html")
}
