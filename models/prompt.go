package models

type UserPromptData struct {
	CoverLetterContent string `json:"coverLetterContent"`
	JobHistory         string `json:"jobHistory"`
	JobTitle           string `json:"jobTitle"`
	JobDescription     string `json:"jobDescription"`
}
