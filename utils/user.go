package utils

import (
	"autojob/db"
	"autojob/models"
	"encoding/json"
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

	var jobHistory []models.History
	err = db.SelectHistoriesByUserID(user.ID, &jobHistory)
	if err != nil {
		log.Fatal("Error getting job history: ", err)
	}

	var coverLetter models.Letter

	err = db.SelectLetterByUserID(user.ID, &coverLetter)
	if err != nil {
		log.Fatal("Error getting cover letter: ", err)
	}

	jsonHistory, err := json.Marshal(jobHistory)
	if err != nil {
		log.Fatal(err)
	}

	var applyList []models.Apply

	for _, job := range jobs {
		data := models.UserPromptData{
			CoverLetterContent: coverLetter.Content,
			JobHistory:         string(jsonHistory),
			JobTitle:           job.Title,
			JobDescription:     job.Description,
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Fatal(err)
		}
		jsonResponse, err := askGPT(string(jsonData))
		if err != nil {
			log.Fatal(err)
		}

		var response models.Response

		err = json.Unmarshal([]byte(jsonResponse), &response)
		if err != nil {
			log.Fatal(err)
		}

		if response.IsMatch {
			applyList = append(applyList, models.Apply{
				Title:       job.Title,
				Company:     job.Company,
				Link:        job.Link,
				Description: job.Description,
				CoverLetter: response.CoverLetter,
			})
		}
	}
}
