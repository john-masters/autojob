package main

import (
	"autojob/routes"
	"autojob/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	// utils.DbInit()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mux := http.NewServeMux()

	mux.Handle("/", routes.HomeRoutes())
	mux.Handle("/auth/", http.StripPrefix("/auth", routes.AuthRoutes()))
	mux.Handle("/user/", http.StripPrefix("/user", routes.UserRoutes()))
	mux.Handle("/history/", http.StripPrefix("/history", routes.HistoryRoutes()))
	mux.Handle("/letter/", http.StripPrefix("/letter", routes.LetterRoutes()))

	serverErrChan := make(chan error)

	go func() {
		fmt.Println("Server running on http://localhost:8080")
		err := http.ListenAndServe(":8080", mux)
		if err != nil {
			serverErrChan <- err
		}
	}()

	c := cron.New()
	c.AddFunc("@every 10s", utils.UpdateToApplyList)
	c.Start()

	err = <-serverErrChan
	if err != nil {
		log.Fatal(err)
	}
}
