package handlers

import (
	"autojob/components"
	"autojob/middleware"
	"autojob/models"
	"autojob/utils"
	"fmt"
	"net/http"
	"time"
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

	startTime, err := time.Parse("2006-01", start)
	if err != nil {
		fmt.Println("Error parsing start date:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var finishTime time.Time
	if finish != "" {
		finishTime, err = time.Parse("2006-01", finish)
		if err != nil {
			fmt.Println("Error parsing finish date:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	db, err := utils.DbConnection()
	if err != nil {
		fmt.Println("Error initializing database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	insertHistorySQL := `INSERT INTO history (user_id, name, role, start, finish, current, duties) VALUES (?, ?, ?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertHistorySQL)
	if err != nil {
		fmt.Println("Error preparing SQL statement:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(user.ID, name, role, startTime, finishTime, isCurrent, duties)
	if err != nil {
		fmt.Println("Error executing SQL statement:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/job-history", http.StatusSeeOther)
}

func GetSingleHistory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

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

	var history models.History
	err = db.QueryRow("SELECT * FROM history WHERE id = ? AND user_id = ?", id, user.ID).Scan(&history.ID, &history.UserID, &history.Name, &history.Role, &history.Start, &history.Finish, &history.Current, &history.Duties)
	if err != nil {
		fmt.Println("Database query error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	component := components.HistoryForm("POST", history)
	component.Render(r.Context(), w)
}

func UpdateSingleHistory(w http.ResponseWriter, r *http.Request) {
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

	id := r.PathValue("id")
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

	startTime, err := time.Parse("2006-01", start)
	if err != nil {
		fmt.Println("Error parsing start date:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var finishTime time.Time
	if finish != "" {
		finishTime, err = time.Parse("2006-01", finish)
		if err != nil {
			fmt.Println("Error parsing finish date:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	db, err := utils.DbConnection()
	if err != nil {
		fmt.Println("Error initializing database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	updateHistorySQL := `UPDATE history SET name = ?, role = ?, start = ?, finish = ?, current = ?, duties = ? WHERE id = ? AND user_id = ?`

	statement, err := db.Prepare(updateHistorySQL)
	if err != nil {
		fmt.Println("Error preparing SQL statement:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	result, err := statement.Exec(name, role, startTime, finishTime, isCurrent, duties, id, user.ID)
	if err != nil {
		fmt.Println("Error executing SQL statement:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		fmt.Println("No rows updated")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var history models.History

	err = db.QueryRow("SELECT * FROM history WHERE id = ? AND user_id = ?", id, user.ID).Scan(
		&history.ID,
		&history.UserID,
		&history.Name,
		&history.Role,
		&history.Start,
		&history.Finish,
		&history.Current,
		&history.Duties,
	)
	if err != nil {
		fmt.Println("Error querying the new history:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	component := components.HistoryForm("GET", history)
	component.Render(r.Context(), w)
}

func DeleteSingleHistory(w http.ResponseWriter, r *http.Request) {
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

	id := r.PathValue("id")

	deleteHistorySQL := `DELETE FROM history WHERE id = ? AND user_id = ?`
	statement, err := db.Prepare(deleteHistorySQL)
	if err != nil {
		fmt.Println("Error preparing SQL statement:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(id, user.ID)
	if err != nil {
		fmt.Println("Error executing SQL statement:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// return ok so ui knows no issue in updating
	w.WriteHeader(http.StatusOK)
}
