package controllers

import (
	"fmt"
	"net/http"

	"forum/backend/lib"
	"forum/database"
	"forum/models"
)

func CommentsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := lib.Help.CheckSession(r)
	if err.Exist && err.Type != "CK" {
		lib.Help.ErrorHandler(w, models.Error{}.InternalServerErr())
		fmt.Println(err.Message)
		return
	}

	postID := r.URL.Query().Get("id")

	post, err := database.DB.GetPostWithCommentsByID(postID)
	if err.Exist {
		lib.Help.ErrorHandler(w, models.Error{}.InternalServerErr())
		fmt.Println(err.Message)
		return
	}

	err = FetchlikesComment(&post.Comments, userID, database.DB.Db)
	if err.Exist {
		lib.Help.ErrorHandler(w, models.Error{}.InternalServerErr())
		fmt.Println(err.Message)
		return
	}

	err = lib.Help.RenderPage(w, "templates/comments.html", post)
	if err.Exist {
		lib.Help.ErrorHandler(w, models.Error{}.InternalServerErr())
		fmt.Println(err.Message)
		return
	}
}
