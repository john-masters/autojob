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

	err = db.QueryRow("SELECT id, email, password FROM users WHERE email = $1;", email).Scan(
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

	err = db.QueryRow("SELECT id, first_name, last_name, email, password, is_member, is_admin FROM users WHERE id = $1;", ID).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.IsMember,
		&user.IsAdmin,
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

	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1;", email).Scan(&count)
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

	insertUserSQL := "INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4);"
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

	updateUserSQL := "UPDATE users SET first_name = $1, last_name = $2, email = $3, is_member = $4, password = $5 WHERE id = $6;"

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

	rows, err := db.Query("SELECT id, first_name, last_name, email FROM users WHERE is_member = TRUE;")
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

func SelectAllUsers(userList *[]models.User) error {
	db, err := conn()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users;")
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
			&user.Password,
			&user.IsMember,
			&user.IsAdmin,
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

func UpdateUserMemberStatusByID(ID int) error {
	var user models.User

	err := SelectUserByID(ID, &user)
	if err != nil {
		return err
	}

	db, err := conn()
	if err != nil {
		return err
	}
	defer db.Close()

	updateUserSQL := "UPDATE users SET is_member = $1 WHERE id = $2;"

	statement, err := db.Prepare(updateUserSQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	isMember := !user.IsMember

	result, err := statement.Exec(
		isMember,
		ID,
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

func UpdateUserAdminStatusByID(ID int) error {
	var user models.User

	err := SelectUserByID(ID, &user)
	if err != nil {
		return err
	}

	db, err := conn()
	if err != nil {
		return err
	}
	defer db.Close()

	updateUserSQL := "UPDATE users SET is_admin = $1 WHERE id = $2;"

	statement, err := db.Prepare(updateUserSQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	isAdmin := !user.IsAdmin

	result, err := statement.Exec(
		isAdmin,
		ID,
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

func DeleteUserByID(id int) error {
	db, err := conn()
	if err != nil {
		return err
	}
	defer db.Close()

	deleteHistorySQL := "DELETE FROM users WHERE id = $1;"
	statement, err := db.Prepare(deleteHistorySQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
