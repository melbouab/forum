package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"forum/backend/lib"
	"forum/database"
	"forum/models"
)

func HandleLikes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		lib.Help.ErrorHandler(w, models.Error{}.MethodNotAllowed())
		return
	}

	userID, Err := lib.Help.CheckSession(r)
	if Err.Exist {
		if Err.Type == "CK" {
			http.Redirect(w, r, "/logout", http.StatusSeeOther)
			return
		}
		fmt.Println(Err.Message)
		lib.Help.ErrorHandler(w, models.Error{}.InternalServerErr())
		return
	}

	is_Post, is_comment, interactionType, posteOrComment_id := PostOrComment(r)
	if !is_Post && !is_comment {
		lib.Help.ErrorHandler(w, models.Error{}.BadRequest())
		return
	}

	idCheckQuery := ""
	switch is_Post {
	case true:
		idCheckQuery = "SELECT id FROM posts WHERE id = ?"
	case false:
		idCheckQuery = "SELECT id FROM comments WHERE id = ?"
	}

	var Id int
	err := database.DB.Db.QueryRow(idCheckQuery, posteOrComment_id).Scan(&Id)
	if err != nil {
		if err == sql.ErrNoRows {
			lib.Help.ErrorHandler(w, models.Error{}.BadRequest())
			fmt.Println(err)
			return
		}
		lib.Help.ErrorHandler(w, models.Error{}.InternalServerErr())
		fmt.Println(err)
		return
	}

	likeOrNot := -1

	if interactionType == "like" {
		likeOrNot = 1
	}

	if is_Post {
		er := PostInteractions("post_id", posteOrComment_id, userID, likeOrNot, database.DB.Db)
		if er.Exist {
			fmt.Println(Err.Message)
			lib.Help.ErrorHandler(w, models.Error{}.InternalServerErr())
			return
		}

		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	er := PostInteractions("Coment_id", posteOrComment_id, userID, likeOrNot, database.DB.Db)
	if er.Exist {
		fmt.Println(Err.Message)
		lib.Help.ErrorHandler(w, models.Error{}.InternalServerErr())
		return
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func PostInteractions(PostOrComment, posteOrComment_id, userId string, likeOrNot int, repo *sql.DB) models.Error {
	var is_like int
	var id_interaction int
	AlreadyHaveInteraction := fmt.Sprintf("SELECT is_like, id FROM interactions WHERE %s = ? AND user_id = ?", PostOrComment)

	err := repo.QueryRow(AlreadyHaveInteraction, posteOrComment_id, userId).Scan(&is_like, &id_interaction)
	if err == sql.ErrNoRows {
		AddInteraction := fmt.Sprintf(`INSERT INTO interactions (user_id, %s, is_like) VALUES (?, ?, ?)`, PostOrComment)

		_, err = repo.Exec(AddInteraction, userId, posteOrComment_id, likeOrNot)
		if err != nil {
			return models.Error{}.ParseErr(err, "DB", "Inserting <PostInteractions>")
		}
	} else if err != nil {
		return models.Error{}.ParseErr(err, "DB", "Reading is_like <PostInteractions>")
	} else {
		if likeOrNot == is_like {
			DeletInteraction := `DELETE FROM interactions WHERE id = ?`
			_, err = repo.Exec(DeletInteraction, id_interaction)
			if err != nil {
				return models.Error{}.ParseErr(err, "DB", "Deleting <PostInteractions>")
			}
		} else {
			UPDATEInteraction := fmt.Sprintf(`UPDATE interactions SET is_like = %d  WHERE id = ?`, likeOrNot)

			_, err = repo.Exec(UPDATEInteraction, id_interaction)
			if err != nil {
				return models.Error{}.ParseErr(err, "DB", "Updating is_like <PostInteractions>")
			}
		}
	}

	return models.NoErrors
}

func PostOrComment(r *http.Request) (bool, bool, string, string) {
	id := r.URL.Query().Get("id")
	_, err := strconv.Atoi(id)
	if err != nil {
		return false, false, "", ""
	}

	react := r.URL.Query().Get("react")
	postOrComment := r.URL.Query().Get("type")

	return postOrComment == "post" && (react == "like" || react == "dislike"), postOrComment == "comment" && (react == "like" || react == "dislike"), react, id
}

func FetchPostsLikes(posts *[]models.Post, userID string, repo *sql.DB) models.Error {
	err := models.Error{}
	for i := range *posts {
		post := &(*posts)[i]
		Querylikes := fmt.Sprintf(`SELECT COUNT(*) FROM interactions WHERE is_like = 1 AND post_id = %d`, post.Id)
		post.Likes, err = Get_Count(Querylikes, repo)
		if err.Exist {
			return err
		}

		QueryDislikes := fmt.Sprintf(`SELECT COUNT(*) FROM interactions WHERE is_like = -1 AND post_id = %d`, post.Id)
		post.Dislikes, err = Get_Count(QueryDislikes, repo)
		if err.Exist {
			return err
		}

		Queryislikes := fmt.Sprintf(`SELECT is_like FROM interactions WHERE  post_id = %v AND user_id = %v`, post.Id, userID)
		post.IsLiked, err = Get_Count(Queryislikes, repo)
		if err.Exist {
			fmt.Println(post.Id, userID)
			return err
		}
	}
	return models.NoErrors
}

func FetchlikesComment(comments *[]models.Comment, userId string, repo *sql.DB) models.Error {
	err := models.Error{}
	for i := range *comments {
		comment := &(*comments)[i]
		Querylikes := fmt.Sprintf(`SELECT COUNT(*) FROM interactions WHERE is_like = 1 AND Coment_id = %v`, comment.Id)
		comment.Likes, err = Get_Count(Querylikes, repo)
		if err.Exist {
			return err
		}

		QueryDislikes := fmt.Sprintf(`SELECT COUNT(*) FROM interactions WHERE is_like = -1 AND Coment_id = %v`, comment.Id)
		comment.Dislikes, err = Get_Count(QueryDislikes, repo)
		if err.Exist {
			return err
		}

		Queryislikes := fmt.Sprintf(`SELECT is_like FROM interactions WHERE  Coment_id = %v AND user_id = %v`, comment.Id, userId)
		comment.IsLiked, err = Get_Count(Queryislikes, repo)
		if err.Exist {
			return err
		}
	}
	return models.NoErrors
}

func Get_Count(Query string, repo *sql.DB) (int, models.Error) {
	var Count int
	row := repo.QueryRow(Query)
	err := row.Scan(&Count)
	if err != nil && err != sql.ErrNoRows {
		return 0, models.Error{}.ParseErr(err, "DB", "Scanning row <Get_Count>")
	}

	return Count, models.NoErrors
}
