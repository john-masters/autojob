package routes

import (
	"autojob/handlers"
	"autojob/middleware"
	"net/http"
)

func HomeRoutes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("GET /", handlers.HomePage)
	router.HandleFunc("GET /sign-up", handlers.SignupPage)
	router.HandleFunc("GET /account", middleware.RequireAuth(handlers.AccountPage))

	return router
}
