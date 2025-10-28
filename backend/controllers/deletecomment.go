package controllers

import (
	"fmt"
	"net/http"

	"forum/backend/lib"
	"forum/database"
	"forum/models"
)

func DeleteComment(w http.ResponseWriter, r *http.Request) {
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

	commentId := r.FormValue("id")

	postId, err := database.DB.GetPostIdByCommentId(commentId)
	if err.Exist {
		lib.Help.ErrorHandler(w, models.Error{}.InternalServerErr())
		fmt.Println(Err.Message)
		return
	}

	if postId == "" {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	err = database.DB.DeleteCommentfromDB(commentId, userID)
	if err.Exist {
		lib.Help.ErrorHandler(w, models.Error{}.InternalServerErr())
		fmt.Println(Err.Message)
		return
	}

	http.Redirect(w, r, "/comments?id="+postId, http.StatusSeeOther)
}
