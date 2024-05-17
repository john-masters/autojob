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

	router.HandleFunc("GET /{id}", middleware.RequireAuth(handlers.GetSingleExperience))
	router.HandleFunc("POST /{id}", middleware.RequireAuth(handlers.UpdateSingleExperience))

	return router
}
