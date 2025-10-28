package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/backend/controllers"
	"forum/config"
	"forum/database"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	if err := database.DB.InitDB(); err.Exist {
		log.Println(err.Message)
	}
}

func main() {
	mux := http.NewServeMux()
	controllers.RoutesHandle(mux)

	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: mux,
	}

	fmt.Println("http://localhost" + server.Addr)
	log.Fatal(server.ListenAndServe())
}
