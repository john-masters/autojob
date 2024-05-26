package db

import "autojob/models"

func SelectQueriesByUserID(userID int, queriesList *[]models.Query) error {
	db, err := conn()
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM queries WHERE user_id = ?", userID)
	if err != nil {
		return err
	}

	for rows.Next() {
		var query models.Query
		err := rows.Scan(
			&query.ID,
			&query.UserID,
			&query.Query,
		)
		if err != nil {
			return err
		}
		*queriesList = append(*queriesList, query)

	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

func InsertQuery(query *models.Query) error {
	db, err := conn()
	if err != nil {
		return err
	}
	defer db.Close()

	insertQuerySQL := `INSERT INTO queries (user_id, query) VALUES (?, ?)`
	statement, err := db.Prepare(insertQuerySQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(
		query.UserID,
		query.Query,
	)

	if err != nil {
		return err
	}
	return nil
}

func DeleteQuery(id int, userID int) error {
	db, err := conn()
	if err != nil {
		return err
	}
	defer db.Close()

	deleteHistorySQL := `DELETE FROM queries WHERE id = ? AND user_id = ?`
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

func SelectQueryCount(user *models.User, count *int) error {
	db, err := conn()
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.QueryRow("SELECT COUNT(*) FROM queries WHERE user_id = ?", &user.ID).Scan(count)
	if err != nil {
		return err
	}
	return nil
}
