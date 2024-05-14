package handlers

import (
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

		fmt.Printf("First Name: %s,\nLast Name: %s,\nEmail: %s,\nPassword: %s\n", firstName, lastName, email, password)

		// TODO: Add your authentication logic here

		fmt.Fprintf(w, "Signup successful for email: %s", email)
	})

	return router

}
