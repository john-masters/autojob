package routes

import (
	"autojob/handlers"
	"autojob/middleware"
	"net/http"
)

func HistoryRoutes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("POST /", middleware.RequireAuth(handlers.CreateHistory))
	router.HandleFunc("GET /", middleware.RequireAuth(handlers.GetHistory))

	router.HandleFunc("GET /{id}", middleware.RequireAuth(handlers.GetSingleHistory))
	router.HandleFunc("POST /{id}", middleware.RequireAuth(handlers.UpdateSingleHistory))
	router.HandleFunc("DELETE /{id}", middleware.RequireAuth(handlers.DeleteSingleHistory))

	return router
}
