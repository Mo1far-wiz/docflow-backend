package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "api.db")

	if err != nil {
		panic("Could not connect to DB.")
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createTables()
}

func createTables() {
	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		firstName TEXT NOT NULL,
		lastName TEXT NOT NULL,
		dateOfBirth DATETIME NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`

	_, err := DB.Exec(createUserTable)

	if err != nil {
		panic("Could not create table.")
	}

	createDocsTable := `
	CREATE TABLE IF NOT EXISTS docs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		docName TEXT NOT NULL,
		dateTime DATETIME NOT NULL,
		faculty TEXT NOT NULL,
		specialty TEXT NOT NULL,
		yearOfStudy INTEGER,
		user_id INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	_, err = DB.Exec(createDocsTable)

	if err != nil {
		panic("Could not create table.")
	}
}
