package controllers

import (
	"fmt"
	"net/http"

	"forum/backend/lib"
	"forum/database"
	"forum/models"
)

func DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		lib.Help.ErrorHandler(w, models.Error{}.MethodNotAllowed())
		return
	}
	userID, err := lib.Help.CheckSession(r)
	if err.Exist {
		if err.Type == "CK" {
			http.Redirect(w, r, "/logout", http.StatusSeeOther)
			return
		}
		lib.Help.ErrorHandler(w, models.Error{}.InternalServerErr())
		fmt.Println(err.Message)
		return
	}

	postID := r.URL.Query().Get("id")

	isPostCreator, err := lib.Help.IsPostCreator(userID, postID)
	if err.Exist {
		lib.Help.ErrorHandler(w, models.Error{}.InternalServerErr())
		fmt.Println(err.Message)
		return
	}

	if !isPostCreator {
		lib.Help.ErrorHandler(w, models.Error{}.Forbidden())
		return
	}

	er := database.DB.DeletePostfromDB(postID, userID)
	if er != nil {
		lib.Help.ErrorHandler(w, models.Error{}.InternalServerErr())
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
