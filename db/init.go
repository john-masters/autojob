package db

import "fmt"

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
			password TEXT NOT NULL,
			is_member BOOLEAN NOT NULL DEFAULT FALSE
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

		CREATE TABLE IF NOT EXISTS letters (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT (datetime('now')),
			FOREIGN KEY (user_id) REFERENCES users (id)
		);
	`)
	if err != nil {
		fmt.Println("Error creating table")
		return
	}

	fmt.Println("Successfully initialized database")
}
