package main

import (
	"autojob/handlers"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	// utils.DbInit()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mux := http.NewServeMux()

	mux.Handle("/", handlers.HomeRoutes())
	mux.Handle("/auth/", http.StripPrefix("/auth", handlers.AuthRoutes()))

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
