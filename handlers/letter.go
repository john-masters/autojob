package handlers

import (
	"autojob/middleware"
	"autojob/models"
	"autojob/utils"
	"fmt"
	"net/http"
)

func CreateLetter(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	err := r.ParseForm()

	if err != nil {
		fmt.Println("Error parsing form")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	content := r.FormValue("content")

	if content == "" {
		fmt.Fprint(w, "Content is required")
	}

	db, err := utils.DbConnection()
	if err != nil {
		fmt.Println("Error initializing database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	insertHistorySQL := `INSERT INTO letters (user_id, content) VALUES (?, ?)`
	statement, err := db.Prepare(insertHistorySQL)
	if err != nil {
		fmt.Println("Error preparing SQL statement:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(user.ID, content)

	if err != nil {
		fmt.Println("Error executing SQL statement:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/cover-letter", http.StatusSeeOther)
}

func DeleteLetter(w http.ResponseWriter, r *http.Request) {
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

	deleteHistorySQL := `DELETE FROM letters WHERE user_id = ?`
	statement, err := db.Prepare(deleteHistorySQL)
	if err != nil {
		fmt.Println("Error preparing SQL statement:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(user.ID)

	if err != nil {
		fmt.Println("Error executing SQL statement:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/cover-letter", http.StatusSeeOther)

}
