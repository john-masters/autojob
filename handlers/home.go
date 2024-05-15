package handlers

import (
	"autojob/components"
	"autojob/middleware"
	"autojob/models"
	"fmt"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	component := components.Home()
	component.Render(r.Context(), w)
}

func SignupPage(w http.ResponseWriter, r *http.Request) {
	component := components.Signup()
	component.Render(r.Context(), w)
}

func AccountPage(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}
	fmt.Println("I'm logged in as", user.FirstName, user.LastName)

	component := components.Account()
	component.Render(r.Context(), w)
}
