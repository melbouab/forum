package helpers

import (
	"net/http"
	"text/template"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	succesffuly(w)
}
func succesffuly(w http.ResponseWriter) {
	tmpl, err := template.ParseFiles("./frontend/html/home.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
