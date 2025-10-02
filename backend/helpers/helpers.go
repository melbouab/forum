package helpers

import (
	"net/http"
	"text/template"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	tmpl, _ := template.ParseFiles("frontend/html/login.html")
	tmpl.Execute(w, nil)
}
