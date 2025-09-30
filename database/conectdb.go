package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // درايفر SQLite
)

// Connect كتحل الكونيكسيون و كترجعها
func Connect() (*sql.DB, error) {
	// (تصحيح 1: حيدنا defer db.Close() من هنا)
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return nil, err
	}

	// كنتأكدو أن الكونيكسيون خدامة مزيان
	if err = db.Ping(); err != nil {
		db.Close() // إلا كان شي مشكل فالـ ping، كنسدو الكونيكسيون
		return nil, err
	}

	log.Println("Connected to database successfully")
	return db, nil
}
