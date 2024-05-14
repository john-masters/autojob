package handlers

import (
	"autojob/models"
	"autojob/utils"
	"database/sql"
	"fmt"
	"net/http"
)

func AuthRoutes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()

		if err != nil {
			fmt.Println("Error parsing form")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		switch {
		case email == "":
			fmt.Fprint(w, "Email is required")
			return
		case password == "":
			fmt.Fprint(w, "Password is required")
			return
		}

		db, err := utils.DbConnection()
		if err != nil {
			fmt.Println("Error initializing database")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer db.Close()

		var user models.User
		err = db.QueryRow("SELECT id, email, password FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Fprint(w, "Invalid email address or password")
			} else {
				fmt.Println("Database query error:", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		if user.Password != password {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Invalid email address or password")
			return
		}

		fmt.Fprintf(w, "Login successful for email: %s", email)

	})

	router.HandleFunc("POST /signup", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()

		if err != nil {
			fmt.Println("Error parsing form")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		switch {
		case firstName == "":
			fmt.Fprint(w, "First Name is required")
			return
		case lastName == "":
			fmt.Fprint(w, "Last Name is required")
			return
		case email == "":
			fmt.Fprint(w, "Email is required")
			return
		case password == "":
			fmt.Fprint(w, "Password is required")
			return
		}

		db, err := utils.DbConnection()
		if err != nil {
			fmt.Println("Error initializing database")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer db.Close()

		var userCount int
		err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&userCount)
		if err != nil {
			fmt.Println("Error querying database:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if userCount > 0 {
			fmt.Fprint(w, "An account with this email already exists")
			return
		}

		insertUserSQL := `INSERT INTO users (first_name, last_name, email, password) VALUES (?, ?, ?, ?)`
		statement, err := db.Prepare(insertUserSQL)
		if err != nil {
			fmt.Println("Error preparing SQL statement:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer statement.Close()

		_, err = statement.Exec(firstName, lastName, email, password)
		if err != nil {
			fmt.Println("Error executing SQL statement:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Signup successful for email: %s", email)
	})

	return router

}
