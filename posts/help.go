package database

import (
	"database/sql"
	"fmt"
)

type PostRepo struct {
	Db *sql.DB
}

type Post struct {
	ID         int    `json:"id"`
	UserId     int    `json:"user_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CategoryID int    `json:"category_id"`
	CreatedAt  string `json:"created_at"`
}

type User struct {
	ID       int    `json:"id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// it create a table post if it not exist at the first time
func (r *PostRepo) CreateTablePostIfNotExist() error {
	_, er := r.Db.Exec(`
		CREATE TABLE IF NOT EXISTS posts(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			title TEXT,
			content TEXT,
			category_id INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	return fmt.Errorf("error creating table: %w", er)
}

// method used to insert posts content into the Post table
func (r *PostRepo) InsertPost(post Post) error {
	_, er := r.Db.Exec(
		"INSERT INTO posts (user_id, title, content, category_id,  created_at) VALUES (?, ?, ?, ?, ?)",
		post.UserId, post.Title, post.Content, post.CategoryID, post.CreatedAt,
	)
	if er != nil {
		return fmt.Errorf("error inserting: %w", er)
	}
	return nil
}

// this return all the posts or an error
func (r *PostRepo) GetAllPosts() ([]Post, error) {
	rows, er := r.Db.Query("SELECT * FROM posts")
	if er != nil {
		return nil, fmt.Errorf("error select table: %w", er)
	}
	defer rows.Close()

	var res []Post
	for rows.Next() {
		var p Post
		er := rows.Scan(&p.ID, &p.UserId, &p.Title, &p.Content, &p.CategoryID, &p.CreatedAt)
		if er != nil {
			return nil, fmt.Errorf("error scanning rows: %w", er)
		}
		res = append(res, p)
	}
	return res, nil
}

// this return a post asked
func (r *PostRepo) GetPost(id int) (Post, error) {
	var p Post
	er := r.Db.QueryRow("SELECT * FROM posts WHERE id = ?", id).Scan(&p.ID, &p.UserId, &p.Title, &p.Content, &p.CategoryID, &p.CreatedAt)
	if er != nil {
		return Post{}, fmt.Errorf("error selecting post: %w", er)
	}
	return p, nil
}

// this update post in the data base
func (r *PostRepo) UpdatePost(p Post) error {
	_, err := r.Db.Exec("UPDATE posts SET user_id=?, title=?, content=?, category_id=? WHERE id=? ",
		p.UserId, p.Title, p.Content, p.CategoryID, p.ID,
	)
	return err
}

// delete a post in the data base
func (r *PostRepo) Delete(id int) error {
	_, err := r.Db.Exec("DELETE FROM posts WHERE id=?", id)
	return err
}

// get post by category
func (r *PostRepo) GetPostByCategory(id int) ([]Post, error) {
	rows, er := r.Db.Query("SELECT * FROM posts WHERE  category_id=?", id)
	if er != nil {
		return []Post{}, fmt.Errorf("error select posts: %w", er)
	}
	defer rows.Close()

	var res []Post
	for rows.Next() {
		var p Post
		er := rows.Scan(&p.ID, &p.UserId, &p.Title, &p.Content, &p.CategoryID, &p.CreatedAt)
		if er != nil {
			return []Post{}, fmt.Errorf("error scanning row: %w", er)
		}
		res = append(res, p)
	}
	return res, nil
}