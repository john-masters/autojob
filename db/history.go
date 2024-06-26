package db

import (
	"autojob/models"
)

func SelectHistoriesByUserID(userID int, histories *[]models.History) error {
	db, err := conn()
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM history WHERE user_id = $1 ORDER BY start DESC;", userID)
	if err != nil {
		return err
	}

	for rows.Next() {
		var history models.History
		err := rows.Scan(&history.ID, &history.UserID, &history.Name, &history.Role, &history.Start, &history.Finish, &history.Current, &history.Duties)
		if err != nil {
			return err
		}
		*histories = append(*histories, history)

	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

func SelectHistoryByIDAndUserID(id int, userID int, history *models.History) error {
	db, err := conn()
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.QueryRow("SELECT * FROM history WHERE id = $1 AND user_id = $2;", id, userID).Scan(
		&history.ID,
		&history.UserID,
		&history.Name,
		&history.Role,
		&history.Start,
		&history.Finish,
		&history.Current,
		&history.Duties,
	)

	if err != nil {
		return err
	}

	return nil
}

func SelectHistoryCount(user *models.User, count *int) error {
	db, err := conn()
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.QueryRow("SELECT COUNT(*) FROM history WHERE user_id = $1;", &user.ID).Scan(count)
	if err != nil {
		return err
	}
	return nil
}

func InsertHistory(history *models.History) error {
	db, err := conn()
	if err != nil {
		return err
	}
	defer db.Close()

	insertHistorySQL := "INSERT INTO history (user_id, name, role, start, finish, current, duties) VALUES ($1, $2, $3, $4, $5, $6, $7);"
	statement, err := db.Prepare(insertHistorySQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(
		history.UserID,
		history.Name,
		history.Role,
		history.Start,
		history.Finish,
		history.Current,
		history.Duties,
	)

	if err != nil {
		return err
	}
	return nil
}

func UpdateHistory(history *models.History) error {
	db, err := conn()
	if err != nil {
		return err
	}
	defer db.Close()

	updateHistorySQL := "UPDATE history SET name = $1, role = $2, start = $3, finish = $4, current = $5, duties = $6 WHERE id = $7 AND user_id = $8;"

	statement, err := db.Prepare(updateHistorySQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	result, err := statement.Exec(
		history.Name,
		history.Role,
		history.Start,
		history.Finish,
		history.Current,
		history.Duties,
		history.ID,
		history.UserID,
	)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		return err
	}
	return nil
}

func DeleteHistory(id int, userID int) error {
	db, err := conn()
	if err != nil {
		return err
	}
	defer db.Close()

	deleteHistorySQL := "DELETE FROM history WHERE id = $1 AND user_id = $2;"
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
