package lib

import (
	"fmt"
	"html/template"
	"net/http"

	"forum/database"
	"forum/models"
)

// Helpers is an empty struct that serves as a receiver for helper methods.
type Helpers struct{}

var Help Helpers

// this make the user create a new post; get (content, category), and insert everyting to the data base
func (a *Helpers) KeepUserCreatePost(w http.ResponseWriter, r *http.Request, userID string) models.Error {
	err := r.ParseForm()
	if err != nil {
		return models.Error{}.ParseErr(err, "DB", "Parsing form")
	}

	categories := r.Form["category"]
	content := r.FormValue("content")

	if len(content) == 0 || len(content) > 300 {
		return models.Error{}.ParseErr(err, "INV", "")
	}

	postID, Err := database.DB.CreatePost(userID, content)
	if Err.Exist {
		return Err
	}

	for _, category := range categories {
		categoryID, err := database.DB.GetCategoryIdByName(category)
		if err != nil {
			return models.Error{}.ParseErr(err, "INV", "")
		}
		err = database.DB.LinkPosttoCategory(postID, categoryID)
		if err != nil {
			return models.Error{}.ParseErr(err, "DB", "Linking <LinkPosttoCategory>")
		}
	}
	return models.NoErrors
}

// serves static files,
// if you want to serve folder or file add it as argument
func (a *Helpers) StaticsHandler(mux *http.ServeMux, statics ...string) {
	for _, v := range statics {
		mux.Handle("/"+v+"/", http.StripPrefix("/"+v+"/", http.FileServer(http.Dir(v))))
	}
}

// renders an HTML page with data.
// returns error if rendering fails.
func (a *Helpers) RenderPage(w http.ResponseWriter, path string, data any) models.Error {
	templ, err := template.ParseFiles(path)
	if err != nil {
		return models.Error{}.ParseErr(err, "TMPL", "Parsing file")
	}

	err = templ.Execute(w, data)
	if err != nil {
		return models.Error{}.ParseErr(err, "TMPL", "Execute file")
	}
	return models.NoErrors
}

// this check session sended by the browser if it is included in the data base
func (*Helpers) CheckSession(r *http.Request) (string, models.Error) {
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		return "", models.Error{}.ParseErr(err, "CK", "Missing cookie")
	}
	return database.DB.GetUserIDBySession(cookie.Value)
}

// ErrorPage renders error page with the specified HTTP status code.
func (a *Helpers) ErrorHandler(w http.ResponseWriter, err models.Error) {
	tmpl, parseErr := template.ParseFiles("templates/error.html")
	if parseErr != nil {
		fmt.Fprintln(w, "Internal server error")
		fmt.Println(parseErr)
		return
	}
	if err.Status < 100 {
		err.Status = http.StatusInternalServerError
	}
	w.WriteHeader(err.Status)
	parseErr = tmpl.Execute(w, err)
	if parseErr != nil {
		fmt.Fprintln(w, "Internal server error")
		fmt.Println(parseErr)
		return
	}
}

func (a *Helpers) IsPostCreator(userID, postID string) (bool, models.Error) {
	var creatorID string
	err := database.DB.Db.QueryRow(`
		SELECT creator_id 
		FROM posts 
		WHERE posts.id = ?
	`, postID).Scan(&creatorID)
	if err != nil {
		return false, models.Error{}.ParseErr(err, "DB", "Reading post creatorID <IsCreator>")
	}
	return creatorID == userID, models.NoErrors
}
