// package helpers - RegisterHandler and related (corrected)
package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"forum/backend/models"
	"forum/database"
	"net/http"
	"text/template"

	"golang.org/x/crypto/argon2"
)

func RegisterGET(w http.ResponseWriter) {
	tmpl, err := template.ParseFiles("./frontend/html/signUp.html")
	if err != nil {
		http.Error(w, "Could not load the page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		RegisterGET(w)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm") // Assuming corrected HTML form field name

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ? OR email = ?", username, email).Scan(&count)
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "Username or email already exists", http.StatusBadRequest)
		return
	}

	if password != confirmPassword {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	salt := make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)
	encodedHash := fmt.Sprintf("%s.%s", saltBase64, hashBase64)

	user := models.User{
		UserName: username,
		Email:    email,
		Password: encodedHash,
	}
	_, err = db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		user.UserName, user.Email, user.Password)
	fmt.Println("INSERTING succesffuly")
	if err != nil {
		// Check for unique constraint violation if needed, but SQLite will error on duplicate
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
