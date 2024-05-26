package db

import (
	"autojob/models"
	"database/sql"
)

func SelectLetterByUserID(userID int, letter *models.Letter) error {
	db, err := conn()
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.QueryRow("SELECT * FROM letters WHERE user_id = ?", userID).Scan(
		&letter.ID,
		&letter.UserID,
		&letter.Content,
		&letter.CreatedAt,
	)

	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

func InsertLetter(letter *models.Letter) error {
	db, err := conn()
	if err != nil {
		return err
	}
	defer db.Close()

	insertHistorySQL := `INSERT INTO letters (user_id, content) VALUES (?, ?)`
	statement, err := db.Prepare(insertHistorySQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(letter.UserID, letter.Content)

	if err != nil {
		return err
	}

	return nil
}

func DeleteLetterByUserID(userID int) error {
	db, err := conn()
	if err != nil {
		return err
	}
	defer db.Close()

	deleteLetterSQL := `DELETE FROM letters WHERE user_id = ?`
	statement, err := db.Prepare(deleteLetterSQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(userID)
	if err != nil {
		return err
	}
	return nil
}

func SelectLetterCount(user *models.User, count *int) error {
	db, err := conn()
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.QueryRow("SELECT COUNT(*) FROM letters WHERE user_id = ?", &user.ID).Scan(count)
	if err != nil {
		return err
	}
	return nil
}
