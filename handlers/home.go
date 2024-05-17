package handlers

import (
	"autojob/components"
	"autojob/middleware"
	"autojob/models"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	component := components.HomePage()
	component.Render(r.Context(), w)
}

func SignupPage(w http.ResponseWriter, r *http.Request) {
	component := components.SignupPage()
	component.Render(r.Context(), w)
}

func AccountPage(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	component := components.AccountPage(user)
	component.Render(r.Context(), w)
}

func ExperiencePage(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	component := components.ExperiencePage(user)
	component.Render(r.Context(), w)
}

func CreateExperiencePage(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	component := components.CreateExperiencePage(user)
	component.Render(r.Context(), w)
}

func SettingsPage(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	component := components.SettingsPage(user)
	component.Render(r.Context(), w)
}
