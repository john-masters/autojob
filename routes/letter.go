package routes

import (
	"autojob/handlers"
	"autojob/middleware"
	"net/http"
)

func LetterRoutes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /", middleware.RequireAuth(handlers.CreateLetter))
	router.HandleFunc("DELETE /", middleware.RequireAuth(handlers.DeleteLetter))

	router.HandleFunc("GET /count", middleware.RequireAuth(handlers.GetLetterCount))

	return router
}
