package routes

import (
	"autojob/handlers"
	"autojob/middleware"
	"net/http"
)

func UserRoutes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("GET /", middleware.RequireAuth(handlers.GetUser))
	router.HandleFunc("POST /", middleware.RequireAuth(handlers.UpdateUser))

	router.HandleFunc("PUT /member/{id}", middleware.RequireAuth(middleware.RequireAdmin(handlers.MakeMember)))
	router.HandleFunc("PUT /admin/{id}", middleware.RequireAuth(middleware.RequireAdmin(handlers.MakeAdmin)))
	router.HandleFunc("DELETE /{id}", middleware.RequireAuth(middleware.RequireAdmin(handlers.DeleteUser)))

	return router
}
