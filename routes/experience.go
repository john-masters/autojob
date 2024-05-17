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
	// TODO: Implement get and post single with experience form response
	// remember to validate the user id matches the id of the experience
	router.HandleFunc("GET /{id}", middleware.RequireAuth(handlers.GetSingleExperience))

	return router
}
