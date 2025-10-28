package controllers

import (
	"net/http"

	"forum/backend/lib"
)

func RoutesHandle(mux *http.ServeMux) {
	lib.Help.StaticsHandler(mux, "media", "styles")

	routes := map[string]http.HandlerFunc{
		"/logout": Logout,

		"/interaction":   HandleLikes,
		"/deletepost":    DeletePost,
		"/createpost":    CreatePost,
		"/createcomment": CreateComment,
		"/deletecomment": DeleteComment,

		"/comments": CommentsHandler,
		"/home":     Home,
		"/login":    LoginHandler,
		"/register": RegisterHandler,
		"/":         GuestHandler,
	}
	for path, handler := range routes {
		mux.HandleFunc(path, handler)
	}
}
