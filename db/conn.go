package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func conn() (*sql.DB, error) {
	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		fmt.Println("Error opening database:", err)
		return db, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return db, err
	}

	return db, nil
}
