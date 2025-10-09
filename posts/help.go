package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
)

type PostRepo struct {
	Db *sql.DB
}

type Post struct {
	ID        int       `json:"id"`
	UserId    int       `json:"userid"`
	UserName  string    `json:"username"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	ImagePath string    `json:"imagepath"`
	Category  Category  `json:"category"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
	Comments  []Comment `json:"comments"`
	CreatedAt string    `json:"createdat"`
}

type User struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type Comment struct{}

// create table of user if not exist
func (r *PostRepo) CreateTableUserIfNotExist() error {
	_, er := r.Db.Exec(CreateUserTabel)
	if er != nil {
		return er
	}
	return nil
}

// it create a table post if it not exist at the first time
func (r *PostRepo) CreateTablePostIfNotExist() error {
	_, er := r.Db.Exec(CreatePostTable)
	if er != nil {
		return er
	}
	return nil
}

// this create the sessions table if it doesn't exist
func (r *PostRepo) CreateTableSessionIfNotExist() error {
	_, er := r.Db.Exec(CreateTableSessions)
	if er != nil {
		return er
	}
	return nil
}

// create categories table
func (r *PostRepo) CreateTableCategoryIfNotExist() error {
	_, err := r.Db.Exec(CreateTableOfCategories)
	if err != nil {
		return err
	}
	return nil
}

// insert default categories
func (r *PostRepo) InsertCategories() error {
	categories := []string{"Sport", "Politic", "Economic", "Music", "Education"}
	for _, name := range categories {
		_, err := r.Db.Exec("INSERT OR IGNORE INTO categories(name) VALUES(?)", name)
		if err != nil {
			return err
		}
	}
	return nil
}

// insert user
func (r *PostRepo) InsertUserInUserTable(username, email, password string) error {
	_, err := r.Db.Exec(Insert_User_Into_User_Table_Query, username, email, password)
	if err != nil {
		return fmt.Errorf("error at inserting user: %w", err)
	}
	return nil
}

// check if the user exist in the data base
func (r *PostRepo) CheckIfUserExistInDB(username, password string) (int, error) {
	var id int
	err := r.Db.QueryRow(Check_If_User_Exist_In_Data_Base_Query, username, password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// get user by session id
func (r *PostRepo) GetUserIDBySession(sessionID string) (int, error) {
	var userid int
	err := r.Db.QueryRow(Select_User_By_Session_Id_Query, sessionID).Scan(&userid)
	if err != nil {
		return 0, err
	}
	return userid, nil
}

// create session for a user
func (r *PostRepo) CreateSessionforUser(userId int) (string, error) {
	generate := func() (string, error) {
		b := make([]byte, 32)
		_, err := rand.Read(b)
		if err != nil {
			return "", err
		}
		return base64.URLEncoding.EncodeToString(b), nil
	}

	sessionId, er := generate()
	if er != nil {
		return "", fmt.Errorf("error at generating session id: %w", er)
	}
	_, er = r.Db.Exec(Insert_Session_ID_Into_User_Query, userId, sessionId)
	if er != nil {
		return "", fmt.Errorf("failed to create session: %w", er)
	}
	return sessionId, nil
}

// create post
func (r *PostRepo) CreatePost(userID int, title, content string, categoryID int, imagePath, createdAt string) error {
	_, err := r.Db.Exec(CreateNewPostQuery, userID, title, content, categoryID, imagePath, createdAt)
	return err
}

// this return all the posts or an error
func (r *PostRepo) GetAllPosts() ([]Post, error) {
	rows, err := r.Db.Query(Select_all_posts_from_db)
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

// start the data base
var CreateUserTabel = `
		CREATE TABLE IF NOT EXISTS users(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE,
			email TEXT UNIQUE,
			password TEXT
		);
	`

var CreatePostTable = `
		CREATE TABLE IF NOT EXISTS posts(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			userid INTEGER,
			title TEXT,
			content TEXT,
			categoryid INTEGER,
			imagepath TEXT,
			likes INTEGER DEFAULT 0,
			dislikes INTEGER DEFAULT 0,
			createdat DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (userid) REFERENCES users(id),
			FOREIGN KEY (categoryid) REFERENCES categories(id)
		);
	`

var CreateTableSessions = `
		CREATE TABLE IF NOT EXISTS sessions(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			userid INTEGER,
			sessionid TEXT UNIQUE,
			createdat DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (userid) REFERENCES users(id)
		);
	`

var CreateTableOfCategories = `
		CREATE TABLE IF NOT EXISTS categories(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE
		);
	`

func StartDataBase(repo *PostRepo) {
	if err := repo.CreateTableUserIfNotExist(); err != nil {
		fmt.Println("error creating user table: ", err)
		return
	}
	if err := repo.CreateTablePostIfNotExist(); err != nil {
		fmt.Println("error creating post table: ", err)
		return
	}
	if err := repo.CreateTableSessionIfNotExist(); err != nil {
		fmt.Println("error creating session table: ", err)
		return
	}
	if err := repo.CreateTableCategoryIfNotExist(); err != nil {
		fmt.Println("error creating categories table: ", err)
		return
	}
	if err := repo.InsertCategories(); err != nil {
		fmt.Println("error creating categories table: ", err)
		return
	}
}

// queries
var (
	Insert_User_Into_User_Table_Query      = "INSERT INTO users(username, email, password) VALUES(?, ?, ?)"
	Check_If_User_Exist_In_Data_Base_Query = "SELECT id FROM users WHERE username = ? AND password = ?"
	Select_User_By_Session_Id_Query        = "SELECT userid FROM sessions WHERE sessionid = ?"
	Insert_Session_ID_Into_User_Query      = "INSERT INTO sessions(userid, sessionid) VALUES(?, ?)"
)

var CreateNewPostQuery = `
INSERT INTO posts (userid, title, content, categoryid, imagepath, likes, dislikes, createdat) VALUES (?, ?, ?, ?, ?, 0, 0, ?)
`

var Select_all_posts_from_db = `
	SELECT posts.id, posts.userid, COALESCE(users.username,'Unknown') AS username, 
		posts.title, 
		posts.content, 
		posts.imagepath,
  		categories.id AS categoryid, 
		categories.name AS categoryname, 
		posts.likes, 
		posts.dislikes, 
		posts.createdat
	FROM posts
	LEFT JOIN users ON posts.userid = users.id
	LEFT JOIN categories ON posts.categoryid = categories.id
	ORDER BY posts.createdat DESC
`
