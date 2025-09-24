package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", Home)
	log.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("frontend/html/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err=r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// username:=r.FormValue("username")
	// password:=r.FormValue("password")

	tmpl.Execute(w, nil)
}
