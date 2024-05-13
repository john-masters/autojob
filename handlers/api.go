package handlers

import (
	"fmt"
	"net/http"
)

func ApiRoutes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from api route")
	})

	return router
}
