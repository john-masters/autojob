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

func CreateExperience(w http.ResponseWriter, r *http.Request) {
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

	insertExperienceSQL := `INSERT INTO experiences (user_id, name, role, start, finish, current, duties) VALUES (?, ?, ?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertExperienceSQL)
	if err != nil {
		fmt.Println("Error preparing SQL statement:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	res, err := statement.Exec(user.ID, name, role, startTime, finishTime, isCurrent, duties)
	if err != nil {
		fmt.Println("Error executing SQL statement:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		fmt.Println("Error getting last insert ID:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var experience models.Experience

	err = db.QueryRow("SELECT id, user_id, name, role, start, finish, current, duties FROM experiences WHERE id = ?", lastInsertID).Scan(
		&experience.ID, &experience.UserID, &experience.Name, &experience.Role, &experience.Start, &experience.Finish, &experience.Current, &experience.Duties)
	if err != nil {
		fmt.Println("Error querying the new experience:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	component := components.ExperienceForm("GET", experience)
	component.Render(r.Context(), w)
}

func GetExperience(w http.ResponseWriter, r *http.Request) {
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

	rows, err := db.Query("SELECT * FROM experiences WHERE user_id = ?", user.ID)
	if err != nil {
		fmt.Println("Database query error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var experiences []models.Experience
	for rows.Next() {
		var experience models.Experience
		err := rows.Scan(&experience.ID, &experience.UserID, &experience.Name, &experience.Role, &experience.Start, &experience.Finish, &experience.Current, &experience.Duties)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		experiences = append(experiences, experience)

	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating rows:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	component := components.ExperienceList(experiences)
	component.Render(r.Context(), w)
}

func GetSingleExperience(w http.ResponseWriter, r *http.Request) {
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

	var experience models.Experience
	err = db.QueryRow("SELECT * FROM experiences WHERE id = ? AND user_id = ?", id, user.ID).Scan(&experience.ID, &experience.UserID, &experience.Name, &experience.Role, &experience.Start, &experience.Finish, &experience.Current, &experience.Duties)
	if err != nil {
		fmt.Println("Database query error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	component := components.ExperienceForm("POST", experience)
	component.Render(r.Context(), w)
}

func UpdateSingleExperience(w http.ResponseWriter, r *http.Request) {
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

	db, err := utils.DbConnection()
	if err != nil {
		fmt.Println("Error initializing database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	updateExperienceSQL := `UPDATE experiences SET name = ?, role = ?, start = ?, finish = ?, current = ?, duties = ? WHERE id = ? AND user_id = ?`

	statement, err := db.Prepare(updateExperienceSQL)
	if err != nil {
		fmt.Println("Error preparing SQL statement:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	result, err := statement.Exec(name, role, start, finish, isCurrent, duties, id, user.ID)
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

	var experience models.Experience

	err = db.QueryRow("SELECT * FROM experiences WHERE id = ? AND user_id = ?", id, user.ID).Scan(
		&experience.ID,
		&experience.UserID,
		&experience.Name,
		&experience.Role,
		&experience.Start,
		&experience.Finish,
		&experience.Current,
		&experience.Duties,
	)
	if err != nil {
		fmt.Println("Error querying the new experience:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	component := components.ExperienceForm("GET", experience)
	component.Render(r.Context(), w)
}

func DeleteSingleExperience(w http.ResponseWriter, r *http.Request) {
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

	deleteExperienceSQL := `DELETE FROM experiences WHERE id = ? AND user_id = ?`
	statement, err := db.Prepare(deleteExperienceSQL)
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
