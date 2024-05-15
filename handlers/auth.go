package handlers

import (
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

	db, err := utils.DbConnection()
	if err != nil {
		fmt.Println("Error initializing database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var user models.User
	err = db.QueryRow("SELECT id, email, password FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.Password)
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

	db, err := utils.DbConnection()
	if err != nil {
		fmt.Println("Error initializing database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var userCount int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&userCount)
	if err != nil {
		fmt.Println("Error querying database:", err)
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

	insertUserSQL := `INSERT INTO users (first_name, last_name, email, password) VALUES (?, ?, ?, ?)`
	statement, err := db.Prepare(insertUserSQL)
	if err != nil {
		fmt.Println("Error preparing SQL statement:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(firstName, lastName, email, string(hash))
	if err != nil {
		fmt.Println("Error executing SQL statement:", err)
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
