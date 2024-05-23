package middleware

import (
	"autojob/models"
	"autojob/utils"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type contextKey string

const UserContextKey contextKey = "user"

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := r.Cookie("Authorization")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		token, err := jwt.Parse(tokenString.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if int64(time.Now().Unix()) > int64(claims["exp"].(float64)) {
				fmt.Println("expired cookie")
				w.WriteHeader(http.StatusUnauthorized)
			}

			var user models.User

			db, err := utils.DbConnection()
			if err != nil {
				fmt.Println("Error initializing database")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer db.Close()

			err = db.QueryRow("SELECT id, first_name, last_name, email, password FROM users WHERE id = ?", claims["sub"]).Scan(
				&user.ID,
				&user.FirstName,
				&user.LastName,
				&user.Email,
				&user.Password,
			)
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

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			r = r.WithContext(ctx)

			next(w, r)

		} else {
			fmt.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
}
