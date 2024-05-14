package main

import (
	"autojob/components"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		component := components.Home()

		component.Render(r.Context(), w)
	})

	http.HandleFunc("GET /sign-up", func(w http.ResponseWriter, r *http.Request) {
		component := components.Signup()

		component.Render(r.Context(), w)
	})

	http.HandleFunc("POST /auth/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("meesage rec'd")

		err := r.ParseForm()

		if err != nil {
			fmt.Println("Error parsing form")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")
		fmt.Printf("Email: %s, Password: %s\n", email, password)

		// TODO: Add your authentication logic here

		fmt.Fprintf(w, "Login successful for email: %s", email)

	})

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
