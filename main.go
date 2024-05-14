package main

import (
	"autojob/components"
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	homeMux := http.NewServeMux()
	authMux := http.NewServeMux()

	homeMux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		component := components.Home()

		component.Render(r.Context(), w)
	})

	homeMux.HandleFunc("GET /sign-up", func(w http.ResponseWriter, r *http.Request) {
		component := components.Signup()

		component.Render(r.Context(), w)
	})

	authMux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()

		if err != nil {
			fmt.Println("Error parsing form")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")
		fmt.Printf("Email: %s,\nPassword: %s\n", email, password)

		// TODO: Add your authentication logic here

		fmt.Fprintf(w, "Login successful for email: %s", email)

	})

	authMux.HandleFunc("POST /signup", func(w http.ResponseWriter, r *http.Request) {
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
		fmt.Printf("First Name: %s,\nLast Name: %s,\nEmail: %s,\nPassword: %s\n", firstName, lastName, email, password)

		// TODO: Add your authentication logic here

		fmt.Fprintf(w, "Signup successful for email: %s", email)
	})

	mux.Handle("/", homeMux)
	mux.Handle("/auth/", http.StripPrefix("/auth", authMux))

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}