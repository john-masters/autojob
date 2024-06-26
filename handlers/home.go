package handlers

import (
	"autojob/components"
	"autojob/db"
	"autojob/middleware"
	"autojob/models"
	"fmt"
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

func HistoryPage(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	var historyList []models.History

	err := db.SelectHistoriesByUserID(user.ID, &historyList)
	if err != nil {
		fmt.Println("Error selecting histories by user ID:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	component := components.HistoryPage(historyList)
	component.Render(r.Context(), w)
}

func LetterPage(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}
	var letter models.Letter

	err := db.SelectLetterByUserID(user.ID, &letter)
	if err != nil {
		fmt.Println("Error selecting user by ID:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	component := components.LetterPage(letter)
	component.Render(r.Context(), w)
}

func QueriesPage(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	var queriesList []models.Query
	err := db.SelectQueriesByUserID(user.ID, &queriesList)
	if err != nil {
		fmt.Println("Error selecting histories by user ID:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	component := components.QueriesPage(&queriesList)
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

func JobsPage(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	var jobsList []models.Job
	err := db.SelectJobsByUserID(user.ID, &jobsList)
	if err != nil {
		fmt.Println("Error selecting jobs by user ID:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	component := components.JobsPage(&jobsList)
	component.Render(r.Context(), w)
}

func AdminPage(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	var userList []models.User

	err := db.SelectAllUsers(&userList)
	if err != nil {
		fmt.Println("Error getting all users:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	component := components.AdminPage(&userList)
	component.Render(r.Context(), w)
}
