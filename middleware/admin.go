package middleware

import (
	"autojob/models"
	"net/http"
)

func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(UserContextKey).(models.User)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !user.IsAdmin {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next(w, r)
	}
}
