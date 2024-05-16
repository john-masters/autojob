package handlers

import (
	"autojob/components"
	"autojob/middleware"
	"autojob/models"
	"autojob/utils"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func Get(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	component := components.SettingsForm(user, false, "POST")
	component.Render(r.Context(), w)
}

func Post(w http.ResponseWriter, r *http.Request) {
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

	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	switch {
	case firstName == "":
		fmt.Fprint(w, "First Name is required")
		return
	case lastName == "":
		fmt.Fprint(w, "Last Name is required")
		return
	case email == "":
		fmt.Fprint(w, "Email is required")
		return
	case password == "":
		fmt.Fprint(w, "Password is required")
		return
	}

	db, err := utils.DbConnection()
	if err != nil {
		fmt.Println("Error initializing database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		fmt.Println("Error generating hash from password:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	updateUserSQL := `UPDATE users SET first_name = ?, last_name = ?, email = ?, password = ? WHERE id = ?`

	statement, err := db.Prepare(updateUserSQL)
	if err != nil {
		fmt.Println("Error preparing SQL statement:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	result, err := statement.Exec(firstName, lastName, email, string(hash), user.ID)
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

	var updatedUser models.User

	// get updated user

	err = db.QueryRow("SELECT * FROM users WHERE id = ?", user.ID).Scan(&updatedUser.ID, &updatedUser.FirstName, &updatedUser.LastName, &updatedUser.Email, &updatedUser.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Invalid ID in cookie:", err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			fmt.Println("Database query error:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// set cookie

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		fmt.Println("Error generating JWT: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "Authorization",
		Value:   tokenString,
		Expires: time.Now().Add(time.Hour * 24),
		Path:    "/",
	})

	component := components.SettingsForm(updatedUser, true, "GET")
	component.Render(r.Context(), w)
}
