package controllers

/*
import (
	"fmt"
	"net/http"
	"strings"

	"forum/server/backend/lib"
	DB "forum/server/database"
)

type Post struct {
	Id          string
	SenderName  string
	SenderId    int
	Content     string
	Category    string
	CreatedAt   string
	Likes       int
	Dislikes    int
	CommentsCnt int
}

func HandleFilter(repo *DB.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			///error method not allowed
			return
		}
		userId, havSession := lib.Help.CheckSession(r, repo)
		if !havSession {
			///rerror can't filter he is just a guest
			return
		}
		validUrl, CategorieFiltre := PostOrComment(r.URL.Path)
		if !validUrl {
			// error invalid path
			return
		}
		if CategorieFiltre == "user" {
			q := `SELECT * FROM posts WHERE user_id = ?`
			Rows, err := repo.Db.Query(q, userId)
			if err != nil {
				// err server
				fmt.Println(err)
				return
			}
			for Rows.Next() {
				var post Post
				err := Rows.Scan(post.Id,post.SenderId)
			}

		}
	}
}

func PostOrComment(url string) (bool, string) {
	Spleted := strings.Split(url, "/")
	if len(Spleted) == 2 {
		return Spleted[0] == "Filter" && (Spleted[1] == "user" || Spleted[1] == "liked"), Spleted[1]
		// Interaction/posts/Like
	}
	return false, ""
}
*/

