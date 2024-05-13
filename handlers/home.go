package handlers

import (
	"autojob/components"
	"net/http"
)

func HomeRoutes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		component := components.Home()

		component.Render(r.Context(), w)
	})

	router.HandleFunc("/sign-up", func(w http.ResponseWriter, r *http.Request) {
		component := components.Signup()

		component.Render(r.Context(), w)
	})

	return router
}
