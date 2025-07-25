package router

import (
	"net/http"

	"real_time/backend/handler"
)

func Router() {
	http.HandleFunc("/", handler.HomeHandler)
	http.HandleFunc("/register", handler.RegisterHandler)
	http.HandleFunc("/frontend/src/", handler.StaticController)
}
