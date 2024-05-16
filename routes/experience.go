package routes

import (
	"autojob/handlers"
	"autojob/middleware"
	"net/http"
)

func ExperienceRoutes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /", middleware.RequireAuth(handlers.CreateExperience))
	router.HandleFunc("GET /", middleware.RequireAuth(handlers.GetExperience))

	return router
}
