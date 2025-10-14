package controllers

import (
	"net/http"

	"forum/server/backend/lib"
	DB "forum/server/database"
)

func DeletePost(w http.ResponseWriter, r *http.Request, repo *DB.Repo) {
	if r.Method != http.MethodPost {
		return
	}
	_, isValid := lib.Help.CheckSession(r, repo)
	if !isValid {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	postId := r.FormValue("post_id")
	userid := r.FormValue("user_id")
	er := repo.DeletePostfromDB(postId, userid)
	if er != nil {
		lib.Help.ErrorPage(w, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
