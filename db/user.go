package db

import (
	"autojob/models"
	"database/sql"
	"fmt"
	"log"
)

func SelectUserByEmail(email string, user *models.User) error {
	db, err := conn()
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.QueryRow("SELECT id, email, password FROM users WHERE email = ?", email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		} else {
			return err
		}
	}

	return nil
}

func SelectUserByID(ID int, user *models.User) error {
	db, err := conn()
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.QueryRow("SELECT id, first_name, last_name, email, is_member, password FROM users WHERE id = ?", ID).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.IsMember,
		&user.Password,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		} else {
			return err
		}
	}

	return nil
}

func SelectUserCountByEmail(email string, count *int) error {
	db, err := conn()
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	if err != nil {
		return err
	}
	return nil
}

func InsertUser(user *models.User) error {
	db, err := conn()
	if err != nil {
		return err
	}
	defer db.Close()

	insertUserSQL := `INSERT INTO users (first_name, last_name, email, password) VALUES (?, ?, ?, ?)`
	statement, err := db.Prepare(insertUserSQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUserByID(user *models.User) error {
	db, err := conn()
	if err != nil {
		return err
	}
	defer db.Close()

	updateUserSQL := `UPDATE users SET first_name = ?, last_name = ?, email = ?, is_member = ?, password = ? WHERE id = ?`

	statement, err := db.Prepare(updateUserSQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	result, err := statement.Exec(
		user.FirstName,
		user.LastName,
		user.Email,
		user.IsMember,
		user.Password,
		user.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

func SelectMemberUsersByID(userList *[]models.User) error {
	db, err := conn()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, first_name, last_name, email FROM users WHERE is_member = TRUE")
	if err != nil {
		log.Fatalln("Database query error:", err)
	}

	for rows.Next() {
		var user models.User

		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
		)
		if err != nil {
			log.Fatalln("Error scanning row:", err)
		}
		*userList = append(*userList, user)

	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}
