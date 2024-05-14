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

	fmt.Println("Successfully connected to database")

	return db, err
}

func DbInit() {
	db, err := DbConnection()
	if err != nil {
		fmt.Println("Error initializing database")
		return
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		first_name TEXT,
		last_name TEXT,
		email TEXT,
		password TEXT
	)`)
	if err != nil {
		fmt.Println("Error creating table")
		return
	}

	fmt.Println("Successfully initialized database")
}
