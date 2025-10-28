package controllers

import (
	"fmt"
	"net/http"

	"forum/backend/lib"
	"forum/models"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		lib.Help.ErrorHandler(w, models.NoErrors.MethodNotAllowed())
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

	err := lib.Help.KeepUserCreatePost(w, r, userID)
	if err.Exist {
		if err.Type == "INV" {
			lib.Help.ErrorHandler(w, models.Error{}.BadRequest())
			return
		}
		lib.Help.ErrorHandler(w, models.Error{}.InternalServerErr())
		fmt.Println(err.Message)
		return
	}
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
