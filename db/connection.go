package db

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
