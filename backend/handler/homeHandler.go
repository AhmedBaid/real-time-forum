package handler

import (
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "not found ajmi ",http.StatusNotFound)
		return 
	}
	http.ServeFile(w, r, "frontend/main.html")
}
