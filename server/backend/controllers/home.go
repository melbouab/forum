package controllers

import (
	"fmt"
	"net/http"
	"text/template"

	"forum/server/backend/lib"
	"forum/server/database"
)

func Home(w http.ResponseWriter, r *http.Request, repo *database.Repo, userID int) {
	// create post
	if r.Method == http.MethodPost {
		userID, isValid := lib.Help.CheckSession(r, repo)
		if !isValid {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		lib.Help.KeepUserCreatePost(r, w, repo, userID)
		return
	}

	// fetch all posts
	posts, err := repo.GetAllPosts()
	if err != nil {
		fmt.Println("error from 'server/backend/controllers/home.go', Error fetching posts: ", err)
		posts = []database.Post{}
	}

	tmpl, err := template.ParseFiles("web/html/home.html")
	if err != nil {
		fmt.Println("error from 'server/backend/controllers/home.go', can't parse home.html ")
		lib.Help.InternalServerError(w)
		return
	}
	tmpl.Execute(w, posts)
}
