package middleware

import (
	"fmt"
	"net/http"
)

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello from da middlewares")
		next(w, r)
	}
}
