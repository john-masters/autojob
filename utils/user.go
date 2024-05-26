package utils

import (
	"autojob/db"
	"autojob/models"
	"log"
)

func processUser(user models.User) {
	var queriesList []models.Query
	err := db.SelectQueriesByUserID(user.ID, &queriesList)
	if err != nil {
		log.Fatal("Error getting queries: ", err)
	}

	var jobs []models.Job
	for _, query := range queriesList {
		scrapeJobData(&jobs, query.Query)
	}
	log.Println(jobs)

}
