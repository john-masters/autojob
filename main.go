package main

import (
	"autojob/routes"
	"autojob/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	utils.DbInit()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mux := http.NewServeMux()

	mux.Handle("/", routes.HomeRoutes())
	mux.Handle("/auth/", http.StripPrefix("/auth", routes.AuthRoutes()))
	mux.Handle("/user/", http.StripPrefix("/user", routes.UserRoutes()))
	mux.Handle("/history/", http.StripPrefix("/history", routes.HistoryRoutes()))

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
