package handlers

import (
	"autojob/db"
	"autojob/middleware"
	"autojob/models"
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

	letter := models.Letter{
		UserID:  user.ID,
		Content: content,
	}

	err = db.InsertLetter(&letter)
	if err != nil {
		fmt.Println("Error inserting letter: ", err)
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

	err := db.DeleteLetterByUserID(user.ID)
	if err != nil {
		fmt.Println("Error deleting letter: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/cover-letter", http.StatusSeeOther)

}
