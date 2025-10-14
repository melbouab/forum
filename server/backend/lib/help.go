package lib

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"forum/server/database"

	"golang.org/x/crypto/bcrypt"
)

// Helpers is an empty struct that serves as a receiver for helper methods.
type Helpers struct{}

var Help Helpers

type Error struct {
	Message string
	Error   string
}
type Data struct {
	Error   string
	StatusE string
}

func GenerateNameForImages(fn string) string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprint(hex.EncodeToString(b) + fn)
}

// this make the user create a new post; get (title, content, category, uploaded image), and insert everyting to the data base
func (*Helpers) KeepUserCreatePost(r *http.Request, w http.ResponseWriter, repo *database.Repo, userID int) error {
	r.ParseMultipartForm(10 << 20)
	title := r.FormValue("title")
	content := r.FormValue("content")
	catStr := r.FormValue("category")
	imagePath := r.FormValue("imagepath")

	file, h, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		imagePath = GenerateNameForImages(h.Filename)
		dst, err := os.Create("web/uploads/" + imagePath)
		if err == nil {
			defer dst.Close()
			io.Copy(dst, file)
		}
	}
	if title != "" && content != "" && catStr != "" {
		categoryid, err := strconv.Atoi(catStr)
		if err == nil {
			t := time.Now().Format("2006-01-02 15:04:05")
			err := repo.CreatePost(userID, title, content, categoryid, imagePath, t)
			if err != nil {
				fmt.Println("error from 'server/backend/lib/help.go': ", err)
			}
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

// this check session sended by the browser if it is included in the data base
func (*Helpers) CheckSession(r *http.Request, repo *database.Repo) (int, bool) {
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		return 0, false
	}
	userID, err := repo.GetUserIDBySession(cookie.Value)
	if err != nil {
		return 0, false
	}
	return userID, true
}

func (*Helpers) VerifyPassword(password, encodedHash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(encodedHash), []byte(password))
	if err != nil {
		return errors.New("incorrect password")
	}
	return nil
}

func (a *Helpers) RegisterGET(w http.ResponseWriter) {
	tmpl, err := template.ParseFiles("./web/html/signup.html")
	if err != nil {
		a.ErrorPage(w, http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// this function parse and execute the login page
func (*Helpers) LoginGET(w http.ResponseWriter) {
	tmpl, err := template.ParseFiles("./web/html/login.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// ErrorPage renders error page with the specified HTTP status code.
func (a *Helpers) ErrorPage(w http.ResponseWriter, status int) {
	msg := http.StatusText(status)

	if status == 404 {
		msg = "Page " + msg
	}

	tmp, err := template.New("error").Parse(errorPage)
	if err != nil {
		http.Error(w, strconv.Itoa(status)+" "+msg, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var buf bytes.Buffer
	errExec := tmp.Execute(&buf, Data{Error: msg, StatusE: strconv.Itoa(status)})
	if errExec != nil {
		a.InternalServerError(w)
		return
	}
	w.WriteHeader(status)
	w.Write(buf.Bytes())
}

// renders 500 Internal Server Error page with embedded HTML and CSS.
func (a *Helpers) InternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
		    <meta charset="UTF-8">
		    <meta name="viewport" content="width=device-width, initial-scale=1.0">
		    <title></title>
		</head>

		 <style>
		        html,
		        body {
		            color: white;
		            font-family: Arial, Helvetica, sans-serif;
		            background-color: rgb(48, 47, 47);
		            height: 100vh;
		            display: flex;
		            justify-content: center;
		            align-items: center;
		            flex-direction: column;
		        }



		        .err {
		            color: red;
		            font-size: 40px;
		            font-weight: bold;
		            margin-bottom: 20px;
		        }
		    </style>
		<body>
		    <div class="err">500 Status Internal Server Error</div>
		</body>
		</html>`))
}

var errorPage = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GT</title>
</head>

 <style>
        html,
        body {
            color: white;
            font-family: Arial, Helvetica, sans-serif;
            background-color: rgb(48, 47, 47);
            height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            flex-direction: column;
        }



        .err {
            color: red;
            font-size: 40px;
            font-weight: bold;
            margin-bottom: 20px;
        }
    </style>
<body>
    <div class="err">{{.StatusE}} {{.Error}}  </div>
</body>
</html>`

func (*Helpers) ErrLogin(w http.ResponseWriter, errMsg string, message string, templ string, status int) {
	data := Error{
		Message: message,
		Error:   errMsg,
	}
	w.WriteHeader(status)
	tmpl, err := template.ParseFiles("./web/html/" + templ)
	if err != nil {
		return
	}
	if execErr := tmpl.Execute(w, data); execErr != nil {
		http.Error(w, execErr.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *Helpers) IsUserCridentialCorrect(isexist bool, password, confirmPassword, email, username string) (bool, string) {
	emailRegaxp := regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$`)
	if isexist {
		return false, "Username or email already exists"
	}
	if !emailRegaxp.MatchString(email) {
		return false, "invalid email adress"
	}
	passwordregex := regexp.MustCompile(`^(?=.*[a-b])(?=.*[A-B])[a-zA-Z]{8,12}$`)

	if (password != confirmPassword) && !passwordregex.MatchString(password) {
		return false, "Passwords do not match"
	}
	return true, ""
}
