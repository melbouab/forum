package database

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"time"
	"unicode"

	models "forum/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Repo struct {
	Db *sql.DB
}

var DB = &Repo{}

func (r *Repo) InitDB() models.Error {
	var err error
	r.Db, err = sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		return models.Error{}.ParseErr(err, "DB", "Openning db")
	}

	sqlBytes, err := os.ReadFile("./database/start.sql")
	if err != nil {
		return models.Error{}.ParseErr(err, "DB", "Reading start.sql")
	}
	_, err = r.Db.Exec(string(sqlBytes))
	if err != nil {
		return models.Error{}.ParseErr(err, "DB", "Creating tables")
	}

	// add seeds here to test
	// var count int
	// err = r.Db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	// if err != nil {
	// 	return err
	// }
	// if count == 0 {
	// 	seeds, err := os.ReadFile("./database/seeds.sql")
	// 	if err != nil {
	// 		return err
	// 	}
	// 	_, err = r.Db.Exec(string(seeds))
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return models.NoErrors
}

// this return post seleced by id or an error (posts)
func (r *Repo) GetPostWithCommentsByID(id string) (models.Post, models.Error) {
	var post models.Post

	rows, err := r.Db.Query(`SELECT 
	comments.id,
	comments.sender_id,
	users.name,
	comments.content,
	comments.created_at
	FROM comments
	LEFT JOIN users ON comments.sender_id = users.id
	WHERE comments.post_id = ?
	ORDER BY comments.created_at DESC;`, id)
	if err != nil {
		return post, models.Error{}.ParseErr(err, "DB", "Reading post's comments <GetPostWithCommentsByID>")
	}
	defer rows.Close()

	var createdAt time.Time
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.Id, &comment.SenderId, &comment.SenderName, &comment.Content, &createdAt)
		if err != nil {
			return post, models.Error{}.ParseErr(err, "DB", "Scaning row <GetPostWithCommentsByID>")
		}

		comment.CreatedAt = r.GetTime(createdAt)

		post.Comments = append(post.Comments, comment)
	}

	err = r.Db.QueryRow(`
    SELECT posts.id, posts.creator_id, users.name AS creator_name, posts.content, COUNT(comments.id) AS comment_count, posts.created_at
    FROM posts
    LEFT JOIN users ON users.id = posts.creator_id
    LEFT JOIN comments ON comments.post_id = posts.id
    WHERE posts.id = ?
    GROUP BY posts.id;`, id).Scan(&post.Id, &post.CreatorId, &post.CreatorName, &post.Content, &post.CommentCount, &createdAt)
	if err != nil {
		return post, models.Error{}.ParseErr(err, "DB", "Reading post data <GetPostWithCommentsByID>")
	}

	post.CreatedAt = r.GetTime(createdAt)

	catRows, err := r.Db.Query(`SELECT 
	categories.id, 
	categories.category_name
	FROM categories
	JOIN post_categories ON categories.id = post_categories.category_id
	WHERE post_categories.post_id = ?;`, id)
	if err != nil {
		return post, models.Error{}.ParseErr(err, "DB", "Reading categories <GetPostWithCommentsByID>")
	}
	defer catRows.Close()

	for catRows.Next() {
		var cat models.Category
		err := catRows.Scan(&cat.ID, &cat.Name)
		if err != nil {
			return post, models.Error{}.ParseErr(err, "DB", "Scaning categories <GetPostWithCommentsByID>")
		}
		post.Categories = append(post.Categories, cat)
	}
	return post, models.NoErrors
}

// get comment by id
func (r *Repo) GetPostIdByCommentId(commentId string) (string, models.Error) {
	var postId string
	err := r.Db.QueryRow(`SELECT post_id FROM comments WHERE id = ?`, commentId).Scan(&postId)
	if err == sql.ErrNoRows {
		return "", models.NoErrors
	}
	if err != nil {
		return "", models.Error{}.ParseErr(err, "DB", "Reading post id <GetPostIdByCommentId>")
	}
	return postId, models.NoErrors
}

// create comment
func (r *Repo) CreateComment(senderId, postId, content string) models.Error {
	_, err := r.Db.Exec(`INSERT INTO comments (content, sender_id, post_id, created_at) VALUES (?,?,?,?)`, content, senderId, postId, time.Now())
	if err != nil {
		return models.Error{}.ParseErr(err, "DB", "Inserting comment <CreateComment>")
	}
	return models.NoErrors
}

// delet comment
func (r *Repo) DeleteCommentfromDB(commentID, userID string) models.Error {
	_, err := r.Db.Exec("DELETE FROM comments WHERE id = ? AND sender_id = ?", commentID, userID)
	if err != nil {
		return models.Error{}.ParseErr(err, "DB", "Deleting comment <DeleteCommentfromDB>")
	}
	return models.NoErrors
}

// this delet post from database by id
func (r *Repo) DeletePostfromDB(PostID, UserID string) error {
	_, er := r.Db.Exec("DELETE FROM posts WHERE id = ? AND creator_id = ?", PostID, UserID)
	return er
}

// this return user seleced by id or an error (users)
func (r *Repo) GetUserName(userID string) (string, models.Error) {
	var name string
	err := r.Db.QueryRow("SELECT name FROM users WHERE users.id = ?", userID).Scan(&name)
	if err != nil {
		return "", models.Error{}.ParseErr(err, "DB", "Reading <GetUserName>")
	}
	return name, models.NoErrors
}

// Get categorie by id (categories)
func (r *Repo) GetAllCategories() ([]string, models.Error) {
	var categories []string
	rows, err := r.Db.Query(`SELECT category_name FROM categories`)
	if err != nil {
		return nil, models.Error{}.ParseErr(err, "DB", "Reading categories <GetAllCategories>")
	}
	for rows.Next() {
		var cat string
		rows.Scan(&cat)
		categories = append(categories, cat)
	}
	return categories, models.NoErrors
}

// Get categorie by id (categories)
func (r *Repo) GetCategoryIdByName(cat string) (int, error) {
	var id int
	err := r.Db.QueryRow(`SELECT id FROM categories WHERE category_name = ?`, cat).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// this query is an sql code that select all posts from the posts table in the data base
var Query_Select_all_posts = `
	SELECT posts.id, 
	posts.creator_id, 
	users.name AS creator_name, 
	posts.content, 
	COUNT(comments.id) AS comment_count, 
	posts.created_at
	FROM posts
	LEFT JOIN users ON users.id = posts.creator_id
	LEFT JOIN comments ON comments.post_id = posts.id
	GROUP BY posts.id
	ORDER BY posts.created_at DESC;
`

var Query_Select_category = `
	SELECT categories.id, 
	categories.category_name
	FROM categories
	JOIN post_categories ON categories.id = post_categories.category_id
	WHERE post_categories.post_id = ?;`

// this returns all the posts or an error
func (r *Repo) GetAllPosts() ([]models.Post, models.Error) {
	rows, err := r.Db.Query(Query_Select_all_posts)
	if err != nil {
		return nil, models.Error{}.ParseErr(err, "DB", "Reading posts <GetAllPosts>")
	}

	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post

		var createdAt time.Time

		err := rows.Scan(&post.Id, &post.CreatorId, &post.CreatorName, &post.Content, &post.CommentCount, &createdAt)
		if err != nil {
			return nil, models.Error{}.ParseErr(err, "DB", "Scaning post <GetAllPosts>")
		}

		catR, err := r.Db.Query(Query_Select_category, post.Id)
		if err != nil {
			return nil, models.Error{}.ParseErr(err, "DB", "Reading categories <GetAllPosts>")
		}

		defer catR.Close()

		var cats []models.Category
		for catR.Next() {
			var cat models.Category
			if err := catR.Scan(&cat.ID, &cat.Name); err != nil {
				return nil, models.Error{}.ParseErr(err, "DB", "Scanning category <GetAllPosts>")
			}
			cats = append(cats, cat)
		}

		post.CreatedAt = r.GetTime(createdAt)

		post.Categories = cats
		posts = append(posts, post)
	}

	return posts, models.NoErrors
}

// create post (posts)
func (r *Repo) CreatePost(creatorID, content string) (int, models.Error) {
	res, err := r.Db.Exec(`INSERT INTO posts (creator_id, content) VALUES (?, ?)`, creatorID, content)
	if err != nil {
		return 0, models.Error{}.ParseErr(err, "DB", "Inserting <CreatePost>")
	}
	lst, err := res.LastInsertId()
	if err != nil {
		return 0, models.Error{}.ParseErr(err, "DB", "Reading post_id <CreatePost>")
	}
	return int(lst), models.NoErrors
}

func (r *Repo) LinkPosttoCategory(postID, categoryID int) error {
	_, err := r.Db.Exec(`INSERT OR IGNORE INTO post_categories (post_id, category_id) VALUES (?, ?)`, postID, categoryID)
	return err
}

func (r *Repo) CreateCategory(name string) (int, error) {
	res, err := r.Db.Exec(`INSERT OR IGNORE INTO categories (category_name) VALUES (?)`, name)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

// check if user exist (users)
func (r *Repo) ValidSignUp(name, email, password, confirmeP string) models.Error {
	if len(name) < 3 || len(name) > 12 || strings.ContainsAny(name, " \t\n\r") {
		return models.Error{}.ParseErr(nil, "INV", "Oops! Username size — try between 3 and 12 characters.")
	}
	for _, char := range name {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '_' {
			return models.Error{}.ParseErr(nil, "INV", fmt.Sprintf("Oops! Username format — unallowed character ( %v ).", string(char)))
		}
	}

	if len(email) > 100 {
		return models.Error{}.ParseErr(nil, "INV", "Oops! Email size — try below 100 characters.")
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return models.Error{}.ParseErr(nil, "INV", "Oops! Email format — try a valid email.")
	}

	if len(password) < 8 || len(password) > 64 {
		return models.Error{}.ParseErr(nil, "INV", "Oops! Password size — try between 8 and 64 characters.")
	}

	if password != confirmeP {
		return models.Error{}.ParseErr(nil, "INV", "Oops! The passwords don’t match — double-check and try again.")
	}

	var exists int
	err = r.Db.QueryRow(`SELECT 1 FROM users WHERE name = ? LIMIT 1`, name).Scan(&exists)
	if err == nil {
		return models.Error{}.ParseErr(nil, "INV", "Oops! Someone’s already using that username.")
	} else if err != sql.ErrNoRows {
		return models.Error{}.ParseErr(err, "DB", "Reading name <ValidSignUp>")
	}

	err = r.Db.QueryRow(`SELECT 1 FROM users WHERE email = ? LIMIT 1`, email).Scan(&exists)
	if err == nil {
		return models.Error{}.ParseErr(nil, "INV", "Oops! Someone’s already using that email.")
	} else if err != sql.ErrNoRows {
		return models.Error{}.ParseErr(err, "DB", "Reading email <ValidSignUp>")
	}

	return models.NoErrors
}

// this hashed the passord and create new user into the data base (users)
func (r *Repo) CreateNewUser(username, email, password string) models.Error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.Error{}.ParseErr(err, "HASH", "Hashing password <CreateNewUser>")
	}
	_, err = r.Db.Exec(`
		INSERT INTO users 
			(name, email, password) 
		VALUES (?, ?, ?)
	`, username, email, hashed)
	if err != nil {
		return models.Error{}.ParseErr(err, "DB", "Inserting user <CreateNewUser>")
	}
	return models.NoErrors
}

// this is get the user forom the database (users) based on the Email
func (r *Repo) ValidSignIn(email, password string) models.Error {
	var mail, pw string
	err := r.Db.QueryRow(`
		SELECT email, password 
		FROM users 	
		WHERE email = ?
	`, email).Scan(&mail, &pw)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Error{}.ParseErr(nil, "INV", "Oops! Something’s off — make sure your email is correct.")
		}
		return models.Error{}.ParseErr(nil, "DB", "ERROR: Reading <ValidSignIn>")
	}

	err = bcrypt.CompareHashAndPassword([]byte(pw), []byte(password))
	if err != nil {
		return models.Error{}.ParseErr(nil, "INV", "Oops! Something’s off — make sure your password is correct.")
	}

	return models.NoErrors
}

// create session for a user
func (r *Repo) SetSessionForUser(w http.ResponseWriter, email string) models.Error {
	sessionId := uuid.NewString()

	_, err := r.Db.Exec(`UPDATE users SET session_id = ? WHERE email = ?`, sessionId, email)
	if err != nil {
		return models.Error{}.ParseErr(err, "DB", "Insering <SetSessionforUser>")
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionId,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	return models.NoErrors
}

// get user by session id
func (r *Repo) GetUserIDBySession(sessionID string) (string, models.Error) {
	var userID string
	err := r.Db.QueryRow("SELECT id FROM users WHERE session_id = ?", sessionID).Scan(&userID)
	if err == sql.ErrNoRows {
		return "", models.Error{}.ParseErr(err, "CK", "Missing sessionID <GetUserIDBySession>")
	}
	if err != nil {
		return "", models.Error{}.ParseErr(err, "DB", "Reading user_id <GetUserIDBySession>")
	}
	return userID, models.NoErrors
}

func (r *Repo) GetTime(t time.Time) string {
	tSince := time.Since(t)

	seconds := int(tSince.Seconds())
	minutes := int(tSince.Minutes())
	hours := int(tSince.Hours())
	days := hours / 24
	weeks := days / 7

	switch {
	case seconds < 60:
		return fmt.Sprintf("%d seconds ago", seconds)
	case minutes < 60:
		return fmt.Sprintf("%d minutes ago", minutes)
	case hours < 24:
		return fmt.Sprintf("%d hours ago", hours)
	case days < 7:
		return fmt.Sprintf("%d days ago", days)
	default:
		return fmt.Sprintf("%d weeks ago", weeks)
	}
}

func (r *Repo) GetFromCategoriesWhere(uI string, iL, iP bool, arr []string) ([]models.Post, models.Error) {
	posts := []models.Post{}

	var categoriesSQL string
	if len(arr) > 0 {
		for i := range arr {
			arr[i] = fmt.Sprintf("'%s'", arr[i])
		}
		categoriesSQL = strings.Join(arr, ", ")
	}

	query := `SELECT 
	p.id,
	p.creator_id,
	u.name AS creator_name,
	p.content,
	COUNT(cmt.id) AS comment_count,
	p.created_at,
	GROUP_CONCAT(c.category_name) AS categories,
	IFNULL(i.is_like, 0) AS is_like
	FROM posts AS p
	LEFT JOIN post_categories AS pc ON p.id = pc.post_id
	LEFT JOIN categories AS c ON pc.category_id = c.id
	LEFT JOIN comments AS cmt ON cmt.post_id = p.id
	JOIN users AS u ON p.creator_id = u.id
	LEFT JOIN interactions AS i 
		ON i.post_id = p.id AND i.user_id = ?
	WHERE 1=1`

	if categoriesSQL != "" {
		query += fmt.Sprintf(" AND c.category_name IN (%s)", categoriesSQL)
	}

	if iL {
		query += ` AND i.is_like = 1`
	}

	if iP {
		query += fmt.Sprintf(" AND p.creator_id = %v", uI)
	}

	query += " GROUP BY p.id ORDER BY p.created_at DESC;"

	rows, err := r.Db.Query(query, uI)
	if err != nil {
		return posts, models.Error{}.ParseErr(err, "DB", "Reading <GetFromCategoriesWhere>")
	}
	defer rows.Close()

	for rows.Next() {
		post := models.Post{}
		var cats sql.NullString
		err := rows.Scan(
			&post.Id,
			&post.CreatorId,
			&post.CreatorName,
			&post.Content,
			&post.CommentCount,
			&post.CreatedAt,
			&cats,
			&post.IsLiked,
		)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return posts, models.Error{}.ParseErr(err, "DB", "Scaning <GetFromCategoriesWhere>")
		}

		if cats.Valid && cats.String != "" {
			for _, c := range strings.Split(cats.String, ",") {
				post.Categories = append(post.Categories, models.Category{Name: c})
			}
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return posts, models.Error{}.ParseErr(err, "DB", "Iterating row <GetFromCategoriesWhere>")
	}

	return posts, models.NoErrors
}
