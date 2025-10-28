package models

import (
	"database/sql"
	"fmt"
	"net/http"
)

type Error struct {
	Message string
	Status  int
	Exist   bool
	Type    string
}

func (err Error) ParseErr(e error, errType, msg string) Error {
	return Error{
		Message: fmt.Sprintf("ERROR: %v\n%v", msg, e),
		Type:    errType,
		Exist:   true,
	}
}

func (err Error) MethodNotAllowed() Error {
	return Error{
		Message: "Sorry! This action is not allowed. Please use the correct method.",
		Status:  http.StatusMethodNotAllowed,
	}
}

func (err Error) InternalServerErr() Error {
	return Error{
		Message: "Something went wrong on our end. Please try again later.",
		Status:  http.StatusInternalServerError,
	}
}

func (err Error) PageNotFound() Error {
	return Error{
		Message: "Oops! The page you’re looking for doesn’t exist.",
		Status:  http.StatusNotFound,
	}
}

func (err Error) BadRequest() Error {
	return Error{
		Message: "Oops! It looks like one of the selected categories isn’t allowed.",
		Status:  http.StatusBadRequest,
	}
}

func (err Error) Forbidden() Error {
	return Error{
		Message: "Oops! You don’t have permission to perform this action.",
		Status:  http.StatusForbidden,
	}
}

var NoErrors = Error{Exist: false}

type HomeData struct {
	Logged     bool
	User       string
	Posts      []Post
	Categories []string
}

type User struct {
	Id        int
	Name      string
	Email     string
	Password  string
	SessionID sql.NullString
}

type Post struct {
	Id           int
	CreatorId    int
	CreatorName  string
	Content      string
	Categories   []Category
	CreatedAt    string
	IsLiked      int
	Likes        int
	Dislikes     int
	CommentCount int
	Comments     []Comment
}

type Category struct {
	ID   int
	Name string
}

type Comment struct {
	Id         int
	SenderId   int
	SenderName string
	Content    string
	IsLiked    int
	Likes      int
	Dislikes   int
	CreatedAt  string
}
