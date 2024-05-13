package main

import (
	"autojob/handlers"
	"fmt"
	"net/http"
)

func main() {
	router := http.NewServeMux()

	router.Handle("/", handlers.HomeRoutes())
	router.Handle("/api", handlers.ApiRoutes())

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
