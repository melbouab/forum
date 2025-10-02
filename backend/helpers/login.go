package helpers

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
		http.Error(w, "method not allowd", http.StatusMethodNotAllowed)
		return
	}
	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))

	if username == "" || password == "" {
		http.Error(w, "Username and password are requerd", http.StatusBadRequest)
		return
	}
	user, err := GetUserByuserName(username)
	if err != nil {
		fmt.Println("error get user", err)
		return
	}
	err = VerifyPassword(password, user.Password)
	if err != nil {
		return
	}
	tokenString, err := SignToken(string(user.Id), username)
	if err != nil {
		http.Error(w, "Could not create login token", http.StatusInternalServerError)
		return
	}
	// Send token as a response or as a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "beare",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteStrictMode,
	})
	resp := struct {
		Token string `json:"token"`
	}{
		Token: tokenString,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("encode Execss resp:", err)
	}
}

func GetUserByuserName(username string) (*models.User, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err

		// return nil, utils.ErrorHandlar(err, "internal error")
	}
	defer db.Close()
	var user models.User
	err = db.QueryRow(`SELECT id, first_name, last_name, email, username, password FROM user WHERE username = ?`,
		username).Scan(&user.Id,
		&user.FirstName, &user.LastName,
		&user.Email, &user.UserName, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
			// return nil, utils.ErrorHandlar(err, "database error")
		}
		return nil, err

		// return nil, utils.ErrorHandlar(err, "internal error")
	}
	return &user, nil
}

func LoginGET(w http.ResponseWriter) {
	tmpl, _ := template.ParseFiles("./frontend/html/login.html")
	tmpl.Execute(w, nil)
}
