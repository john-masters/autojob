package handlers

import (
	"autojob/components"
	"autojob/middleware"
	"autojob/models"
	"autojob/utils"
	"fmt"
	"net/http"
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

	_, err = statement.Exec(user.ID, name, role, start, finish, isCurrent, duties)
	if err != nil {
		fmt.Println("Error executing SQL statement:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	component := components.ExperienceForm("GET")
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

	var experience models.Experience
	err = db.QueryRow("SELECT * FROM experiences WHERE user_id = ?", user.ID).Scan(&experience.ID, &experience.UserID, &experience.Name, &experience.Role, &experience.Start, &experience.Finish, &experience.Current, &experience.Duties)
	if err != nil {
		fmt.Println("Database query error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println(experience)

	component := components.ExperienceForm("POST")
	component.Render(r.Context(), w)
}
