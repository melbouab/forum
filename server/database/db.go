package database

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type Repo struct {
	Db *sql.DB
}

type User struct {
	Id       int    `json:"id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Post struct {
	ID        int       `json:"id"`
	UserId    int       `json:"user_id"`
	UserName  string    `json:"user_name"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	ImagePath string    `json:"image_path"`
	Category  Category  `json:"category"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
	Comments  []Comment `json:"comments"`
	CreatedAt string    `json:"created_at"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Comment struct{}

func (r *Repo) InitDB() error {
	db, err := sql.Open("sqlite3", "./server/database/forum.db")
	if err != nil {
		return err
	}
	r.Db = db

	sqlBytes, err := os.ReadFile("./server/database/start.sql")
	if err != nil {
		return err
	}

	_, err = r.Db.Exec(string(sqlBytes))
	if err != nil {
		return err
	}

	// add seeds here to test
	var count int
	err = r.Db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		seeds, err := os.ReadFile("./server/database/seeds.sql")
		if err != nil {
			return err
		}
		_, err = r.Db.Exec(string(seeds))
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repo) DeletePostfromDB(postid, userID string) error {
	Uid, er := strconv.Atoi(userID)
	if er != nil {
		return er
	}
	Pid, er := strconv.Atoi(postid)
	if er != nil {
		return er
	}

	_, er = r.Db.Exec("DELETE FROM posts where id= ? AND user_id = ?", Pid, Uid)
	return er
}

// this query is an sql code that select all posts from the posts table in the data base
var Query_Select_all_posts = `
	SELECT posts.id, posts.user_id, COALESCE(users.user_name,'Unknown') AS user_name, 
		posts.title, 
		posts.content, 
		posts.image_path,
  		categories.id AS category_id,
		categories.name AS category_name,
		posts.likes, 
		posts.dislikes, 
		posts.created_at
	FROM posts
	LEFT JOIN users ON posts.user_id = users.id
	LEFT JOIN categories ON posts.category_id = categories.id
	ORDER BY posts.created_at DESC
`

// this return all the posts or an error (posts)
func (r *Repo) GetAllPosts() ([]Post, error) {
	rows, err := r.Db.Query(Query_Select_all_posts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		var catID int
		var catName string
		var imgPath sql.NullString
		err := rows.Scan(&p.ID, &p.UserId, &p.UserName, &p.Title, &p.Content, &imgPath, &catID, &catName, &p.Likes, &p.Dislikes, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		if imgPath.Valid {
			p.ImagePath = imgPath.String
		}
		p.Category = Category{ID: catID, Name: catName}

		posts = append(posts, p)
	}
	return posts, nil
}

// create post (posts)
func (r *Repo) CreatePost(userID int, title, content string, categoryID int, imagePath, createdAt string) error {
	_, err := r.Db.Exec("INSERT INTO posts (user_id, title, content, category_id, image_path, likes, dislikes, created_at) VALUES (?, ?, ?, ?, ?, 0, 0, ?)", userID, title, content, categoryID, imagePath, createdAt)
	return err
}

// check if user exist (users)
func (r *Repo) IsUserExistInDB(username, email string) (bool, error) {
	var idUser int
	err := r.Db.QueryRow("SELECT id FROM users WHERE user_name = ? OR email = ?", username, email).Scan(&idUser)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// this generate a random string will help for hashing password or session
func (r *Repo) generate() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// get user by session id
func (r *Repo) GetUserIDBySession(sessionID string) (int, error) {
	var userid int
	err := r.Db.QueryRow("SELECT user_id FROM sessions WHERE session_id = ?", sessionID).Scan(&userid)
	if err != nil {
		return 0, err
	}
	return userid, nil
}

// this hashed the passord and create new user into the data base (users)
func (r *Repo) CreateNewUser(username, email, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = r.Db.Exec("INSERT INTO users (user_name, email, password) VALUES (?, ?, ?)", username, email, hashed)
	return err
}

// this is get the user forom the database (users) based on the username
func (r *Repo) GetUserByUsername(username string) (*User, error) {
	var user User
	err := r.Db.QueryRow("SELECT id, email, user_name, password FROM users WHERE user_name = ?", username).Scan(&user.Id, &user.Email, &user.UserName, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// create session for a user
func (r *Repo) CreateSessionforUser(userId int) (string, error) {
	sessionId, er := r.generate()

	if er != nil {
		return "", er
	}
	_, er = r.Db.Exec("INSERT into sessions(user_id, session_id) VALUES(?, ?) ON CONFLICT(user_id) DO UPDATE SET session_id=excluded.session_id", userId, sessionId)
	if er != nil {
		return "", er
	}
	return sessionId, nil
}
