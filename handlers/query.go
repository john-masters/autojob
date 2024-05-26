package handlers

import (
	"autojob/db"
	"autojob/middleware"
	"autojob/models"
	"fmt"
	"net/http"
	"strconv"
)

func CreateQuery(w http.ResponseWriter, r *http.Request) {
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

	query := r.FormValue("query")

	if query == "" {

		fmt.Fprint(w, "Name is required")
		return
	}

	err = db.InsertQuery(&models.Query{
		UserID: user.ID,
		Query:  query,
	})
	if err != nil {
		fmt.Println("Error Creating Query:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/search-terms", http.StatusSeeOther)
}

func DeleteSingleQuery(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error converting string to int: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	err = db.DeleteQuery(intId, user.ID)
	if err != nil {
		fmt.Println("Error deleting query: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// return ok so ui knows no issue in updating
	w.WriteHeader(http.StatusOK)
}

func GetQueryCount(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	var queryCount int

	err := db.SelectQueryCount(&user, &queryCount)
	if err != nil {
		fmt.Println("Error getting user count:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, queryCount)
}
