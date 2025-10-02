package main

import (
	"crypto/tls"
	"fmt"
	"forum/backend/helpers"
	"forum/database" // (تحسين 1: اسم الـ package تبدل)
	"log"
	"net/http"
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

	mux.HandleFunc("/", helpers.Home)
	mux.HandleFunc("/register", helpers.RegisterHandler)
	mux.HandleFunc("/login", helpers.LoginHandler)

	port := ":3000"
	cert := "cert.pem"
	key := "key.pem"

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

// // هاد الفونكشن دابا كترجع لينا http.HandlerFunc

// // ... الفونكشن RegisterGET كتبقى كيفما هي ...

// type User struct {
// 	Username string
// 	Password string
// 	Email    string
// }

// func RegisterPost(w http.ResponseWriter, r *http.Request) { // <-- زدنا 'db' هنا

// 	user := &User{
// 		Username: strings.TrimSpace(r.FormValue("username")), // (تحسين 3: كنزيدو TrimSpace باش نحيدو الفراغات)
// 		Email:    strings.TrimSpace(r.FormValue("email")),
// 		Password: r.FormValue("password"),
// 	}
// 	confirmPassword := r.FormValue("confirm")

// 	if user.Username == "" || user.Email == "" || user.Password == "" {
// 		http.Error(w, "all fields are required", http.StatusBadRequest)
// 		return
// 	}
// 	if user.Password != confirmPassword {
// 		http.Error(w, "passwords do not match", http.StatusBadRequest)
// 		return
// 	}

// 	hashed, err := HashPassword(user.Password)
// 	if err != nil {
// 		http.Error(w, "failed to hash password", http.StatusInternalServerError)
// 		return
// 	}

// 	stmt, err := db.Prepare("INSERT INTO users (username, password, email) VALUES (?, ?, ?)")
// 	if err != nil {
// 		http.Error(w, "db error (prepare)", http.StatusInternalServerError)
// 		return
// 	}
// 	defer stmt.Close()

// 	_, err = stmt.Exec(user.Username, hashed, user.Email)
// 	if err != nil {
// 		// هنا ممكن يكون الخطأ أن الايميل أو اسم المستخدم ديجا كاين
// 		http.Error(w, "failed to insert user (maybe username or email exists)", http.StatusInternalServerError)
// 		return
// 	}

// 	http.Redirect(w, r, "/", http.StatusSeeOther)
// }

// // ... الفونكشن HashPassword كتبقى كيفما هي ...
// func HashPassword(pass string) (string, error) {
// 	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(pass), 12)
// 	return string(hashedBytes), err
// }

// func Login(w http.ResponseWriter, r *http.Request) {
// 	tmpl, err := template.ParseFiles("./frontend/html/signUp.html")
// 	if err != nil {
// 		http.Error(w, "method not allowd", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	if r.Method == http.MethodGet {
// 		tmpl.Execute(w, nil)
// 		return
// 	}
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "method not allowd", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	user := &User{
// 		Username: strings.TrimSpace(r.FormValue("username")), // (تحسين 3: كنزيدو TrimSpace باش نحيدو الفراغات)
// 		Password: r.FormValue("password"),
// 	}
// 	tmpl1, err := template.ParseFiles("./frontend/html/home.html")
// 	if err != nil {
// 		http.Error(w, "method not allowd", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	tmpl1.Execute(w, user)
// }
