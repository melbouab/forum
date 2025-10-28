package controllers

import (
	"net/http"

	"forum/backend/lib"
	"forum/models"
)

func GuestHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		lib.Help.ErrorHandler(w, models.Error{}.PageNotFound())
		return
	}
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
