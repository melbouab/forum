package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	Cntrlr "forum/server/backend/controllers"
	DB "forum/server/database"

	_ "github.com/mattn/go-sqlite3"
)

var repo *DB.Repo

func init() {
	repo = &DB.Repo{}
	if err := repo.InitDB(); err != nil {
		fmt.Println("error starting database: ", err)
		os.Exit(1)
	}
}

func main() {
	mux := http.NewServeMux()
	Cntrlr.RoutesHandle(mux, repo)

	server := &http.Server{
		Addr:    ":3000",
		Handler: mux,
	}

	fmt.Println("http://localhost:3000")
	log.Fatal(server.ListenAndServe())
}
