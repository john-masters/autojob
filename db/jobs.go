package db

import "autojob/models"

func InsertJob(job *models.Job) error {
	db, err := conn()
	if err != nil {
		return err
	}
	defer db.Close()

	insertJobSQL := `INSERT INTO jobs (user_id, title, company, link, description, cover_letter) VALUES (?, ?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertJobSQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(
		job.UserID,
		job.Title,
		job.Company,
		job.Link,
		job.Description,
		job.CoverLetter,
	)

	if err != nil {
		return err
	}
	return nil
}
