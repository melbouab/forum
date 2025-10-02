package main

import (
	"crypto/tls"
	"fmt"
	"forum/backend/helpers"
	"forum/database"
	"log"
	"net/http"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

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
