package controllers

import (
	"net/http"

	"forum/server/database"
)

func Logout(w http.ResponseWriter, r *http.Request, repo *database.Repo) {
	if r.Method != http.MethodPost {
		return
	}
	
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
