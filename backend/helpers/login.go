package helpers

import (
	"database/sql"
	"forum/backend/models"
	"forum/database"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		LoginGET(w)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))

	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}
	user, err := GetUserByuserName(username)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		log.Println("error getting user:", err)
		ErrorPage(w)
		return
	}
	err = VerifyPassword(password, user.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	tokenString, err := SignToken(string(user.Id), username)
	if err != nil {
		http.Error(w, "Could not create login token", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "bearer",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteStrictMode,
	})
	http.Redirect(w, r, "/home", http.StatusSeeOther)

}

func GetUserByuserName(username string) (*models.User, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var user models.User
	err = db.QueryRow(`SELECT id, email, username, password FROM users WHERE username = ?`,
		username).Scan(&user.Id, &user.Email, &user.UserName, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func LoginGET(w http.ResponseWriter) {
	tmpl, err := template.ParseFiles("./frontend/html/login.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func ErrorPage(w http.ResponseWriter) {
	tmpl, err := template.ParseFiles("./frontend/html/error.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
