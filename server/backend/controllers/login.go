package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"forum/server/backend/lib"
	DB "forum/server/database"
)

func LoginHandler(w http.ResponseWriter, r *http.Request, repo *DB.Repo) {
	if r.Method == http.MethodGet {
		lib.Help.LoginGET(w)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not alloweds", http.StatusMethodNotAllowed)
		lib.Help.ErrLogin(w, "err", "method not alloweds", "login.html", http.StatusMethodNotAllowed)
		return
	}
	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))

	if username == "" || password == "" {
		lib.Help.ErrLogin(w, "err", "method not alloweds", "login.html", http.StatusBadRequest)

		return
	}
	user, err := repo.GetUserByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			lib.Help.ErrLogin(w, "err", "erro user", "login.html", http.StatusUnauthorized)

			return
		}
		fmt.Println("error getting user:", err)
		lib.Help.InternalServerError(w)
		return
	}
	err = lib.Help.VerifyPassword(password, user.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create DB session and set cookie
	sessionID, err := repo.CreateSessionforUser(user.Id)
	if err != nil {
		// lib.Help.ErrorPage(w, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteStrictMode,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
