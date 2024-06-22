package routes

import (
	"autojob/handlers"
	"autojob/middleware"
	"net/http"
)

func JobRoutes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("DELETE /{id}", middleware.RequireAuth(handlers.DeleteSingleJob))
	router.HandleFunc("GET /count", middleware.RequireAuth(handlers.GetJobCount))
	router.HandleFunc("POST /scrape", middleware.RequireAuth(middleware.RequireAdmin(handlers.TriggerScrape)))

	return router
}
