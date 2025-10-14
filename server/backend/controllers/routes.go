package controllers

import (
	"net/http"

	"forum/server/backend/midemware"
	"forum/server/database"
)

func RoutesHandle(mux *http.ServeMux, repo *database.Repo) {
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./web/uploads"))))
	mux.Handle("/web/css/", http.StripPrefix("/web/css/", http.FileServer(http.Dir("./web/css"))))

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		LoginHandler(w, r, repo)
	})
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		RegisterHandler(w, r, repo)
	})
	mux.HandleFunc("POST /logout", func(w http.ResponseWriter, r *http.Request) {
		Logout(w, r, repo)
	})
	mux.HandleFunc("POST /delete", func(w http.ResponseWriter, r *http.Request) {
		DeletePost(w, r, repo)
	})
	mux.HandleFunc("/", midemware.HomeMiddleware(Home, repo))
}
