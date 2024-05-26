package db

import "autojob/models"

func SelectHistoriesByUserID(userID int, histories *[]models.History) error {
	db, err := db()
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM history WHERE user_id = ? ORDER BY start DESC", userID)
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
	db, err := db()
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.QueryRow("SELECT * FROM history WHERE id = ? AND user_id = ?", id, userID).Scan(
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

func InsertHistory(history *models.History) error {
	db, err := db()
	if err != nil {
		return err
	}
	defer db.Close()

	insertHistorySQL := `INSERT INTO history (user_id, name, role, start, finish, current, duties) VALUES (?, ?, ?, ?, ?, ?, ?)`
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
	db, err := db()
	if err != nil {
		return err
	}
	defer db.Close()

	updateHistorySQL := `UPDATE history SET name = ?, role = ?, start = ?, finish = ?, current = ?, duties = ? WHERE id = ? AND user_id = ?`

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
	db, err := db()
	if err != nil {
		return err
	}
	defer db.Close()

	deleteHistorySQL := `DELETE FROM history WHERE id = ? AND user_id = ?`
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
