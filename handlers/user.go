package handlers

import (
	"autojob/components"
	"autojob/db"
	"autojob/middleware"
	"autojob/models"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	component := components.SettingsForm(user, "POST")
	component.Render(r.Context(), w)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
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
	member := r.FormValue("member")
	isMember := member == "on"

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

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error generating hash from password:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.UpdateUserByID(&models.User{
		ID:        user.ID,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  string(hash),
		IsMember:  isMember,
	})
	if err != nil {
		fmt.Println("Error updating user:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var updatedUser models.User

	err = db.SelectUserByID(user.ID, &updatedUser)
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

	component := components.SettingsForm(updatedUser, "GET")
	component.Render(r.Context(), w)
}

func MakeMember(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error converting string to int: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.UpdateUserMemberStatusByID(intId)
	if err != nil {
		fmt.Println("Error updating user:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func MakeAdmin(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error converting string to int: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.UpdateUserAdminStatusByID(intId)
	if err != nil {
		fmt.Println("Error updating user:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	id := r.PathValue("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error converting string to int: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.DeleteUserByID(intId)
	if err != nil {
		fmt.Println("Error updating user:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
