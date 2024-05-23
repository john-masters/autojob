package utils

import "autojob/models"

func processUser(user models.User) {
	var jobs []models.Job
	scrapeJobData(&jobs)
}
