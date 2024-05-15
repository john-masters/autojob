package routes

import (
	"autojob/handlers"
	"net/http"
)

func AuthRoutes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /login", handlers.Login)
	router.HandleFunc("POST /signup", handlers.Signup)

	return router
}
