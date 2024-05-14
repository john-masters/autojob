package main

import (
	"autojob/handlers"
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", handlers.HomeRoutes())
	mux.Handle("/auth/", http.StripPrefix("/auth", handlers.AuthRoutes()))

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
