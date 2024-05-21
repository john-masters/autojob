package utils

import (
	"autojob/models"
	"log"
)

func UpdateToApplyList() {
	log.Println("updating...")

	var jobs []models.Job

	scrapeJobData(&jobs)
}
