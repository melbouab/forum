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
		http.Error(w, "Could not load the page.", http.StatusInternalServerError)
		fmt.Println("Template error:", err)
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
		fmt.Println("DB error:", err)
		return
	}
	defer db.Close()

	firstName := r.FormValue("firstname")
	lastName := r.FormValue("lastname")
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirempassword")
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ? OR email = ?", username, email).Scan(&count)
	if err != nil {
		http.Error(w, "Database  failed", http.StatusInternalServerError)
		fmt.Println("DB error:", err)
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
		fmt.Println(err)
		return
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)
	encodedHash := fmt.Sprintf("%s.%s", saltBase64, hashBase64)
	password = encodedHash

	user := models.User{
		FirstName: firstName,
		LastName:  lastName,
		UserName:  username,
		Email:     email,
		Password:  (password),
	}
	_, err = db.Exec("INSERT INTO users (firstname, lastname, username, email, password) VALUES (?, ?, ?, ?, ?)",
		user.FirstName, user.LastName, user.UserName, user.Email, user.Password)

	if err != nil {
		http.Error(w, "Error inserting user", http.StatusInternalServerError)
		fmt.Println("Insert error:", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
