package handlers

import (
	"autojob/db"
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

func Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		fmt.Println("Error parsing form")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	switch {
	case email == "":
		fmt.Fprint(w, "Email is required")
		return
	case password == "":
		fmt.Fprint(w, "Password is required")
		return
	}

	var user models.User
	err = db.SelectUserByEmail(email, &user)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Fprint(w, "Invalid email address or password")
		} else {
			fmt.Println("Database query error:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Fprint(w, "Invalid email address or password")
		return
	}

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

	fmt.Fprint(w, "Log in successful. Go to <a href='/account'>my account</a>.")
}

func Signup(w http.ResponseWriter, r *http.Request) {
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

	isValidPassword := utils.ValidatePassword(password)
	if !isValidPassword {
		fmt.Fprint(w, "Password does not meet complexity requirements")
		return
	}

	isValidEmail := utils.ValidateEmail(email)
	if !isValidEmail {
		fmt.Fprint(w, "Invalid email address")
		return
	}

	var userCount int

	err = db.SelectUserCountByEmail(email, &userCount)
	if err != nil {
		fmt.Println("Error getting user count:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if userCount > 0 {
		fmt.Fprint(w, "An account with this email already exists")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		fmt.Println("Error generating hash from password:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.InsertUser(&models.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  string(hash),
	})
	if err != nil {
		fmt.Println("Error creating user:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Sign up successful, please <a href='/'>log in</a>.")
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "Authorization",
		Value:   "",
		Expires: time.Unix(0, 0),
		Path:    "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
