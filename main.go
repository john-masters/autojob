package main

import (
	"autojob/components"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		component := components.Home()

		component.Render(r.Context(), w)
	})

	http.HandleFunc("/sign-up", func(w http.ResponseWriter, r *http.Request) {
		component := components.Signup()

		component.Render(r.Context(), w)
	})

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
