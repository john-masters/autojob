package db

import "autojob/models"

func InsertJob(job *models.Job) error {
	db, err := conn()
	if err != nil {
		return err
	}
	defer db.Close()

	insertJobSQL := "INSERT INTO jobs (user_id, title, company, link, description, cover_letter) VALUES ($1, $2, $3, $4, $5, $6);"
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

func SelectJobCountByEmail(job *models.Job, count *int) error {
	db, err := conn()
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.QueryRow("SELECT COUNT(*) FROM jobs WHERE user_id = $1 AND link = $2;", &job.UserID, &job.Link).Scan(&count)
	if err != nil {
		return err
	}
	return nil
}

func SelectJobsByUserID(userID int, jobs *[]models.Job) error {
	db, err := conn()
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM jobs WHERE user_id = $1;", userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var job models.Job
		err = rows.Scan(
			&job.ID,
			&job.UserID,
			&job.Title,
			&job.Company,
			&job.Link,
			&job.Description,
			&job.CoverLetter,
		)
		if err != nil {
			return err
		}
		*jobs = append(*jobs, job)
	}
	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

func DeleteJob(id int, userID int) error {
	db, err := conn()
	if err != nil {
		return err
	}
	defer db.Close()

	deleteHistorySQL := `DELETE FROM jobs WHERE id = $1 AND user_id = $2`
	statement, err := db.Prepare(deleteHistorySQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(id, userID)
	if err != nil {
		return err
	}
	return nil
}

func SelectJobCount(user *models.User, count *int) error {
	db, err := conn()
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.QueryRow("SELECT COUNT(*) FROM jobs WHERE user_id = $1;", &user.ID).Scan(count)
	if err != nil {
		return err
	}
	return nil
}
