package router

import (
	"net/http"

	"real_time/backend/handler"
	"real_time/backend/middleware"
)

func Router() {
	http.HandleFunc("/", handler.HomeHandler)
	http.HandleFunc("/register", handler.RegisterHandler)
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/isloged", middleware.IsLogged)
	http.HandleFunc("/frontend/src/", handler.StaticHandler)
}
