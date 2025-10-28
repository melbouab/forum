package controllers

import (
	"fmt"
	"net/http"
	"time"

	"forum/backend/lib"
	"forum/database"
	"forum/models"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	err := models.Error{}
	if r.Method == http.MethodPost {

		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirmPw")

		err = database.DB.ValidSignUp(username, email, password, confirmPassword)
		if err.Exist {
			if err.Type == "DB" {
				lib.Help.ErrorHandler(w, models.Error{}.InternalServerErr())
				fmt.Println(err.Message, "here")
				return
			}
		} else {
			er := database.DB.CreateNewUser(username, email, password)
			if er.Exist {
				lib.Help.ErrorHandler(w, models.Error{}.InternalServerErr())
				fmt.Println(err.Message)
				return
			}

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}
	Rerr := lib.Help.RenderPage(w, "templates/register.html", err)
	if Rerr.Exist {
		lib.Help.ErrorHandler(w, models.Error{}.InternalServerErr())
		fmt.Println(Rerr.Message)
		return
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	err := models.Error{}

	if r.Method == http.MethodPost {

		email := r.FormValue("email")
		password := r.FormValue("password")

		err := database.DB.ValidSignIn(email, password)
		if err.Type == "DB" {

			lib.Help.ErrorHandler(w, models.Error{}.InternalServerErr())
			fmt.Println(err.Message)
			return

		} else {
			err := database.DB.SetSessionForUser(w, email)
			if err.Exist {
				lib.Help.ErrorHandler(w, models.Error{}.InternalServerErr())
				fmt.Println(err.Message)
				return
			}
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}
	}
	if r.Method != http.MethodGet {
		lib.Help.ErrorHandler(w, models.Error{}.MethodNotAllowed())
	}
	Rerr := lib.Help.RenderPage(w, "templates/login.html", err)
	if Rerr.Exist {
		lib.Help.ErrorHandler(w, models.Error{}.InternalServerErr())
		fmt.Println(Rerr.Message)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}

	c, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	sessionToken := c.Value
	database.DB.Db.Exec("DELETE FROM sessions WHERE session_id = ?", sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   "",
		Expires: time.Now(),
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
