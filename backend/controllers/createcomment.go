package controllers

import (
	"fmt"
	"net/http"

	"forum/backend/lib"
	"forum/database"
	"forum/models"
)

func CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		lib.Help.ErrorHandler(w, models.Error{}.MethodNotAllowed())
		return
	}
	userID, Err := lib.Help.CheckSession(r)
	if Err.Exist {
		if Err.Type == "CK" {
			http.Redirect(w, r, "/logout", http.StatusSeeOther)
			return
		}
		lib.Help.ErrorHandler(w, models.Error{}.InternalServerErr())
		fmt.Println(Err.Message)
		return
	}

	postID := r.FormValue("post-id")
	content := r.FormValue("content")

	if len(content) > 300 || len(content) == 0 {
		http.Redirect(w, r, "/comments?id="+postID, http.StatusSeeOther)
		return
	}

	err := database.DB.CreateComment(userID, postID, content)
	if err.Exist {
		lib.Help.ErrorHandler(w, models.Error{}.InternalServerErr())
		fmt.Println(err.Message)
		return
	}
	http.Redirect(w, r, "/comments?id="+postID, http.StatusSeeOther)
}
