package router

import (
	"net/http"

	"real_time/backend/handler"
	"real_time/backend/middleware"
)

func Router() {
	http.HandleFunc("/api/current-user",  middleware.Authorisation((handler.CurrentUserHandler)))
	http.HandleFunc("/ws",  middleware.Authorisation(handler.WsHandler))
	http.HandleFunc("/messages",  middleware.Authorisation(handler.GetMessagesHandler))
	http.HandleFunc("/ReactionHandler", middleware.Authorisation(handler.ReactionHandler))
	http.HandleFunc("/getComments",  middleware.Authorisation(handler.GetComments))
	http.HandleFunc("/getUsers",  middleware.Authorisation(handler.GetUsers))
	http.HandleFunc("/createComment",  middleware.Authorisation(handler.CommentHandler))
	http.HandleFunc("/getPosts", middleware.Authorisation( handler.GetPosts))
	http.HandleFunc("/isloged", middleware.IsLogged)
	http.HandleFunc("/",( handler.HomeHandler))
	http.HandleFunc("/register", handler.RegisterHandler)
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/logout", handler.LogoutHandler)
	http.HandleFunc("/frontend/src/", handler.StaticHandler)
	http.HandleFunc("/createpost", middleware.Authorisation(handler.CreatePost))
}
