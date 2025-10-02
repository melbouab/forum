package helpers

import (
	"database/sql"
	"encoding/json"
	"forum/backend/models"
	"forum/database"
	"log"
	"net/http"
	"strings"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))

	if username == "" || password == "" {
		http.Error(w, "Username and password are requerd", http.StatusBadRequest)
		return
	}
	user, err := GetUserByuserName(username)
	if err != nil {
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
		// return nil, utils.ErrorHandlar(err, "internal error")
	}
	defer db.Close()
	var user models.User
	err = db.QueryRow(`SELECT id, first_name, last_name, email, username, password, inactive_status, role FROM user WHERE username = ?`, username).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.UserName, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {

			// return nil, utils.ErrorHandlar(err, "database error")
		}
		// return nil, utils.ErrorHandlar(err, "internal error")
	}
	return &user, nil
}
