package handlers

import (
	"autojob/utils"
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
			fmt.Println("Email is empty")
			w.WriteHeader(http.StatusBadRequest)
			return
		case password == "":
			fmt.Println("Password is empty")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Printf("Email: %s,\nPassword: %s\n", email, password)

		// TODO: Add your authentication logic here

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
			fmt.Println("First Name is empty")
			w.WriteHeader(http.StatusBadRequest)
			return
		case lastName == "":
			fmt.Println("Last Name is empty")
			w.WriteHeader(http.StatusBadRequest)
			return
		case email == "":
			fmt.Println("Email is empty")
			w.WriteHeader(http.StatusBadRequest)
			return
		case password == "":
			fmt.Println("Password is empty")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		db, err := utils.DbConnection()
		if err != nil {
			fmt.Println("Error initializing database")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer db.Close()

		insertUserSQL := `INSERT INTO users (first_name, last_name, email, password) VALUES (?, ?, ?, ?)`
		statement, err := db.Prepare(insertUserSQL)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer statement.Close()

		_, err = statement.Exec(firstName, lastName, email, password)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Inserted sample data successfully")

		rows, err := db.Query("SELECT id, first_name, last_name, email FROM users")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer rows.Close()

		fmt.Println("Querying data...")

		for rows.Next() {
			var id int
			var firstName string
			var lastName string
			var email string
			err = rows.Scan(&id, &firstName, &lastName, &email)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("ID: %d, First Name: %s, Last Name: %s, Email: %s\n", id, firstName, lastName, email)
		}

		fmt.Fprintf(w, "Signup successful for email: %s", email)
	})

	return router

}
