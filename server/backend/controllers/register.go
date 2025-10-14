package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"forum/server/backend/lib"
	DB "forum/server/database"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request, repo *DB.Repo) {
	if r.Method == http.MethodGet {
		lib.Help.RegisterGET(w)
		return
	}
	if r.Method != http.MethodPost {
		lib.Help.ErrorPage(w, http.StatusMethodNotAllowed)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("password"))
	confirmPassword := strings.TrimSpace(r.FormValue("confirm"))

	isexist, err := repo.IsUserExistInDB(username, email)
	if err != nil {
		fmt.Println("error at checking user existance: ", err)
		lib.Help.ErrorPage(w, http.StatusInternalServerError)
		return
	}

	iscorrect, msg := lib.Help.IsUserCridentialCorrect(isexist, password, confirmPassword, email, username)
	if !iscorrect {
		lib.Help.ErrLogin(w, "Username or email already exists", msg, "signup.html", http.StatusBadRequest)
		return
	}

	err = repo.CreateNewUser(username, email, password)
	if err != nil {
		fmt.Println("Error creating user: ", err)
		lib.Help.ErrorPage(w, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
