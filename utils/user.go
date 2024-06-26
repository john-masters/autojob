package utils

import (
	"autojob/db"
	"autojob/models"
	"encoding/json"
	"log"
	"sync"
)

func processUser(user models.User) {
	var queriesList []models.Query
	err := db.SelectQueriesByUserID(user.ID, &queriesList)
	if err != nil {
		log.Fatal("Error getting queries: ", err)
	}

	var scrapeData []models.ScrapeData
	var wg sync.WaitGroup

	for _, query := range queriesList {
		wg.Add(1)

		go func(query models.Query) {
			defer wg.Done()
			scrapeJobData(&scrapeData, query.Query)
		}(query)

		wg.Wait()
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
		log.Fatal("Error marshalling JSON", err)
	}

	var jobsList []models.Job

	for _, job := range scrapeData {
		data := models.UserPromptData{
			CoverLetterContent: coverLetter.Content,
			JobHistory:         string(jsonHistory),
			JobTitle:           job.Title,
			JobDescription:     job.Description,
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Fatal("Error marshalling JSON", err)
		}
		jsonResponse, err := askGPT(string(jsonData))
		log.Println("jsonResponse: ", jsonResponse)
		if err != nil {
			log.Fatal("Error asking GPT", err)
		}

		var response models.Response

		err = json.Unmarshal([]byte(jsonResponse), &response)
		if err != nil {
			log.Fatal("Error unmarshalling JSON", err)
		}

		if response.IsMatch {
			jobsList = append(jobsList, models.Job{
				UserID:      user.ID,
				Title:       job.Title,
				Company:     job.Company,
				Link:        job.Link,
				Description: job.Description,
				CoverLetter: response.CoverLetter,
			})
		}
	}

	log.Printf("Found %v jobs for user %v", len(jobsList), user.Email)

	for _, job := range jobsList {
		var count int
		err := db.SelectJobCountByEmail(&job, &count)

		if err != nil {
			log.Printf("Error getting job count for job %v: %v", job, err)
			continue
		}

		if count > 0 {
			log.Printf("Job %v already exists, skipping", job)
			continue
		}

		err = db.InsertJob(&job)
		if err != nil {
			log.Printf("Error inserting job %v: %v", job, err)
		}
	}
}
