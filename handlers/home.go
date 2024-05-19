package handlers

import (
	"autojob/components"
	"autojob/middleware"
	"autojob/models"
	"autojob/utils"
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

	db, err := utils.DbConnection()
	if err != nil {
		fmt.Println("Error initializing database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM history WHERE user_id = ? ORDER BY start DESC", user.ID)
	if err != nil {
		fmt.Println("Database query error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var historyList []models.History
	for rows.Next() {
		var history models.History
		err := rows.Scan(&history.ID, &history.UserID, &history.Name, &history.Role, &history.Start, &history.Finish, &history.Current, &history.Duties)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		historyList = append(historyList, history)

	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating rows:", err)
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

	component := components.LetterPage(user)
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
