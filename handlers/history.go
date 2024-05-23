package handlers

import (
	"autojob/components"
	"autojob/middleware"
	"autojob/models"
	"fmt"
	"net/http"
	"strconv"
)

func CreateHistory(w http.ResponseWriter, r *http.Request) {
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

	name := r.FormValue("name")
	role := r.FormValue("role")
	start := r.FormValue("start")
	finish := r.FormValue("finish")
	current := r.FormValue("current")
	duties := r.FormValue("duties")

	switch {
	case name == "":
		fmt.Fprint(w, "Name is required")
		return
	case role == "":
		fmt.Fprint(w, "Role is required")
		return
	case start == "":
		fmt.Fprint(w, "Start date is required")
		return
	case duties == "":
		fmt.Fprint(w, "Duties is required")
		return
	}

	isCurrent := current == "on"

	err = InsertHistory(&models.History{
		UserID:  user.ID,
		Name:    name,
		Role:    role,
		Start:   start,
		Finish:  finish,
		Current: isCurrent,
		Duties:  duties,
	})
	if err != nil {
		fmt.Println("Error Creating History:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/job-history", http.StatusSeeOther)
}

func GetSingleHistory(w http.ResponseWriter, r *http.Request) {
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

	var history models.History

	err = SelectHistoryByIDAndUserID(intId, user.ID, &history)
	if err != nil {
		fmt.Println("Error getting history: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	component := components.HistoryForm("POST", history)
	component.Render(r.Context(), w)
}

func UpdateSingleHistory(w http.ResponseWriter, r *http.Request) {
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

	err = r.ParseForm()

	if err != nil {
		fmt.Println("Error parsing form")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	name := r.FormValue("name")
	role := r.FormValue("role")
	start := r.FormValue("start")
	finish := r.FormValue("finish")
	current := r.FormValue("current")
	duties := r.FormValue("duties")

	switch {
	case name == "":
		fmt.Fprint(w, "Name is required")
		return
	case role == "":
		fmt.Fprint(w, "Role is required")
		return
	case start == "":
		fmt.Fprint(w, "Start date is required")
		return
	case finish == "" && current != "on":
		fmt.Fprint(w, "Finish date is required")
		return
	case current == "on" && finish != "":
		fmt.Fprint(w, "Finish date should be empty")
		return
	case duties == "":
		fmt.Fprint(w, "Duties is required")
		return
	}

	isCurrent := current == "on"

	err = UpdateHistory(&models.History{
		Name:    name,
		Role:    role,
		Start:   start,
		Finish:  finish,
		Current: isCurrent,
		Duties:  duties,
		ID:      intId,
		UserID:  user.ID,
	})
	if err != nil {
		fmt.Println("Error updating history:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var history models.History

	err = SelectHistoryByIDAndUserID(intId, user.ID, &history)
	if err != nil {
		fmt.Println("Error querying the new history:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	component := components.HistoryForm("GET", history)
	component.Render(r.Context(), w)
}

func DeleteSingleHistory(w http.ResponseWriter, r *http.Request) {
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

	err = DeleteHistory(intId, user.ID)
	if err != nil {
		fmt.Println("Error deleting history: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// return ok so ui knows no issue in updating
	w.WriteHeader(http.StatusOK)
}
