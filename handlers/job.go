package handlers

import (
	"autojob/db"
	"autojob/middleware"
	"autojob/models"
	"autojob/utils"
	"fmt"
	"net/http"
	"strconv"
)

func DeleteSingleJob(w http.ResponseWriter, r *http.Request) {
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

	err = db.DeleteJob(intId, user.ID)
	if err != nil {
		fmt.Println("Error deleting query: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// return ok so ui knows no issue in updating
	w.WriteHeader(http.StatusOK)
}

func GetJobCount(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	var jobCount int

	err := db.SelectJobCount(&user, &jobCount)
	if err != nil {
		fmt.Println("Error getting user count:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, jobCount)
}

func TriggerScrape(w http.ResponseWriter, r *http.Request) {
	go func() {
		utils.UpdateToApplyList()
	}()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Scraping job triggered"))
}
