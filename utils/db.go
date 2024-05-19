package utils

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func DbConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./test.db")

	if err != nil {
		fmt.Println("Error opening database")
		return db, err
	}

	return db, err
}

func DbInit() {
	db, err := DbConnection()
	if err != nil {
		fmt.Println("Error initializing database")
		return
	}
	defer db.Close()
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			email TEXT NOT NULL,
			password TEXT NOT NULL
		);

		CREATE TABLE history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			role TEXT NOT NULL,
			start TEXT NOT NULL,
			finish TEXT,
			current BOOLEAN NOT NULL,
			duties TEXT,
			FOREIGN KEY (user_id) REFERENCES user (id)
		);
	`)
	if err != nil {
		fmt.Println("Error creating table")
		return
	}

	fmt.Println("Successfully initialized database")
}
