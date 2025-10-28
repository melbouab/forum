package controllers

import (
	"fmt"
	"net/http"

	"forum/backend/lib"
	"forum/database"
	"forum/models"
)

func Home(w http.ResponseWriter, r *http.Request) {
	logged := true
	userID, err := lib.Help.CheckSession(r)
	if err.Exist {
		if err.Type == "CK" {
			logged = false
		} else {
			lib.Help.ErrorHandler(w, models.Error{}.InternalServerErr())
			fmt.Println(err.Message)
			return
		}
	}

	var posts []models.Post
	if r.Method == http.MethodPost {
		Err := r.ParseForm()
		if Err != nil {
			lib.Help.ErrorHandler(w, models.Error{}.InternalServerErr())
			fmt.Println("ParseForm", err)
			return
		}

		categories := r.Form["category"]
		selectedCategories := []string{}
		liked := false
		posted := false

		for _, category := range categories {
			if category == "created" {
				posted = true
				continue
			} else if category == "liked" {
				liked = true
				continue
			}
			selectedCategories = append(selectedCategories, category)
		}

		posts, err = database.DB.GetFromCategoriesWhere(userID, liked, posted, selectedCategories)
		if err.Exist {
			lib.Help.ErrorHandler(w, models.Error{}.InternalServerErr())
			fmt.Println("ParseForm", err)
			return
		}

	} else {
		posts, err = database.DB.GetAllPosts()
		if err.Exist {
			lib.Help.ErrorHandler(w, models.NoErrors.InternalServerErr())
			fmt.Println(err.Message)
			return
		}

		err = FetchPostsLikes(&posts, userID, database.DB.Db)
		if err.Exist {
			lib.Help.ErrorHandler(w, models.Error{}.InternalServerErr())
			fmt.Println(err.Message)
			return
		}
	}

	username := "Guest"
	if logged {
		username, err = database.DB.GetUserName(userID)
		if err.Exist {
			lib.Help.ErrorHandler(w, models.NoErrors.InternalServerErr())
			fmt.Println(err.Message)
			return
		}
	}

	categories, err := database.DB.GetAllCategories()
	if err.Exist {
		lib.Help.ErrorHandler(w, models.NoErrors.InternalServerErr())
		fmt.Println(err.Message)
		return
	}

	homeData := models.HomeData{
		Logged:     logged,
		User:       username,
		Posts:      posts,
		Categories: categories,
	}
	Err := lib.Help.RenderPage(w, "templates/home.html", homeData)
	if Err.Exist {
		lib.Help.ErrorHandler(w, models.NoErrors.InternalServerErr())
		fmt.Println(err.Message)
	}
}
