package db

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func Init() {
	db, err := conn()
	if err != nil {
		fmt.Println("Error initializing database")
		return
	}
	defer db.Close()
	_, err = db.Exec(`
		DROP TABLE IF EXISTS queries CASCADE;
		DROP TABLE IF EXISTS jobs CASCADE;
		DROP TABLE IF EXISTS letters CASCADE;
		DROP TABLE IF EXISTS history CASCADE;
		DROP TABLE IF EXISTS users CASCADE;

		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			email TEXT NOT NULL,
			password TEXT NOT NULL,
			is_member BOOLEAN NOT NULL DEFAULT FALSE,
			is_admin BOOLEAN NOT NULL DEFAULT FALSE
		);

		CREATE TABLE IF NOT EXISTS history (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			role TEXT NOT NULL,
			start TEXT NOT NULL,
			finish TEXT,
			current BOOLEAN NOT NULL,
			duties TEXT,
			FOREIGN KEY (user_id) REFERENCES users (id)
		);

		CREATE TABLE IF NOT EXISTS letters (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (id)
		);

		CREATE TABLE IF NOT EXISTS queries (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			query TEXT NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users (id)
		);

		CREATE TABLE IF NOT EXISTS jobs (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			company TEXT NOT NULL,
			link TEXT NOT NULL,
			description TEXT NOT NULL,
			cover_letter TEXT NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users (id)
		);
	`)

	if err != nil {
		fmt.Println("Error creating table", err)
		return
	}

	insertAdminSQL := `
	INSERT INTO users (
		first_name,
		last_name,
		email,
		password,
		is_member,
		is_admin
	) VALUES (
		'Admin',
		'User',
		'admin@example.com',
		$1,
		TRUE,
		TRUE
	) ON CONFLICT DO NOTHING;`

	statement, err := db.Prepare(insertAdminSQL)
	if err != nil {
		fmt.Println("Error preparing statement", err)
		return
	}
	defer statement.Close()

	password := os.Getenv("ADMIN_PASS")

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		fmt.Println("Error generating hash from password:", err)
		return
	}

	_, err = statement.Exec(string(hash))
	if err != nil {
		fmt.Println("Error inserting admin user", err)
		return
	}

	fmt.Println("Successfully initialized database")
}
