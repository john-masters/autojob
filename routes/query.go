package routes

import (
	"autojob/handlers"
	"autojob/middleware"
	"net/http"
)

func QueryRoutes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /", middleware.RequireAuth(handlers.CreateQuery))
	router.HandleFunc("DELETE /{id}", middleware.RequireAuth(handlers.DeleteSingleQuery))

	return router
}
