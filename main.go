package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"forum/database" // (تحسين 1: اسم الـ package تبدل)
	"html/template"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// ما بقيناش محتاجين متغير عام
// var db *sql.DB

func main() {
	// كنتكونيكطاو مرة وحدة فالبداية
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close() // هادي هي البلاصة الصحيحة ديالها

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./frontend/css"))
	mux.Handle("/frontend/css/", http.StripPrefix("/frontend/css/", fs))

	mux.HandleFunc("/", Home)
	// (تحسين 2: كندوزو الكونيكسيون 'db' نيشان للـ handler)
	mux.HandleFunc("/register", RegisterHandler(db))
	mux.HandleFunc("/login", Login)

	port := ":3000"
	cert := "cert.pem"
	key := "key.pem"
	// Root handler — simple welcome route
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	server := &http.Server{
		Addr:      port,
		TLSConfig: tlsConfig,
		// Handler:   rl.Middleware(middlwares.Compression(middlwares.Rate_time(middlwares.Cors(mux)))),
		Handler: mux,
	}
	fmt.Println("https://localhost:3000")
	err = server.ListenAndServeTLS(cert, key)
	// Start the HTTP server
	// fmt.Println("server is running on port:", port)
	// err := http.ListenAndServe(port, nil)
	if err != nil {
		// Fatal if server fails to start
		log.Fatalln("Error starting server:", err)
	}
}

// ... الفونكشن Home كتبقى كيفما هي ...
func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	tmpl, _ := template.ParseFiles("frontend/html/login.html")
	tmpl.Execute(w, nil)
}

// هاد الفونكشن دابا كترجع لينا http.HandlerFunc
func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			RegisterGET(w, r)
		case http.MethodPost:
			RegisterPost(db, w, r) // <-- كندوزو 'db'
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// ... الفونكشن RegisterGET كتبقى كيفما هي ...
func RegisterGET(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("./frontend/html/signUp.html")
	tmpl.Execute(w, nil)
}

type User struct {
	Username string
	Password string
	Email    string
}

func RegisterPost(db *sql.DB, w http.ResponseWriter, r *http.Request) { // <-- زدنا 'db' هنا
	r.ParseForm()

	user := &User{
		Username: strings.TrimSpace(r.FormValue("username")), // (تحسين 3: كنزيدو TrimSpace باش نحيدو الفراغات)
		Email:    strings.TrimSpace(r.FormValue("email")),
		Password: r.FormValue("password"),
	}
	confirmPassword := r.FormValue("confirm")

	if user.Username == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "all fields are required", http.StatusBadRequest)
		return
	}
	if user.Password != confirmPassword {
		http.Error(w, "passwords do not match", http.StatusBadRequest)
		return
	}

	hashed, err := HashPassword(user.Password)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	stmt, err := db.Prepare("INSERT INTO users (username, password, email) VALUES (?, ?, ?)")
	if err != nil {
		http.Error(w, "db error (prepare)", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Username, hashed, user.Email)
	if err != nil {
		// هنا ممكن يكون الخطأ أن الايميل أو اسم المستخدم ديجا كاين
		http.Error(w, "failed to insert user (maybe username or email exists)", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ... الفونكشن HashPassword كتبقى كيفما هي ...
func HashPassword(pass string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(pass), 12)
	return string(hashedBytes), err
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowd", http.StatusMethodNotAllowed)
		return
	}
	user := &User{
		Username: strings.TrimSpace(r.FormValue("username")), // (تحسين 3: كنزيدو TrimSpace باش نحيدو الفراغات)
		Password: r.FormValue("password"),
	}
	tmpl, err := template.ParseFiles("./frontend/html/home.html")
	if err != nil {
		http.Error(w, "method not allowd", http.StatusMethodNotAllowed)
		return
	}
	tmpl.Execute(w, user)
}
