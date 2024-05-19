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
	router.HandleFunc("GET /job-history", middleware.RequireAuth(handlers.HistoryPage))
	router.HandleFunc("GET /settings", middleware.RequireAuth(handlers.SettingsPage))

	return router
}
