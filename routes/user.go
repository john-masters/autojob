package routes

import (
	"autojob/handlers"
	"autojob/middleware"
	"net/http"
)

func UserRoutes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("GET /", middleware.RequireAuth(handlers.Get))
	router.HandleFunc("POST /", middleware.RequireAuth(handlers.Post))

	return router
}
