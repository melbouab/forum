package helpers

import (
	"fmt"
	"forum/backend/models"
	"forum/database"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
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

	// Validation بسيطة
	if password != confirmPassword {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	// تشفير الباسورد
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error while hashing password", http.StatusInternalServerError)
		return
	}

	// إنشاء اليوزر
	user := models.User{
		FirstName: firstName,
		LastName:  lastName,
		UserName:  username,
		Email:     email,
		Password:  string(hashedPassword),
	}

	// إدخال البيانات فالداتا بيز
	_, err = db.Exec("INSERT INTO users (firstname, lastname, username, email, password) VALUES (?, ?, ?, ?, ?)",
		user.FirstName, user.LastName, user.UserName, user.Email, user.Password)

	if err != nil {
		http.Error(w, "Error inserting user", http.StatusInternalServerError)
		fmt.Println("Insert error:", err)
		return
	}

	// Response للكلينت
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}
