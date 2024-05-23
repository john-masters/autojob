package utils

import (
	"autojob/models"
	"log"
)

func UpdateToApplyList() {
	log.Println("updating...")
	db, err := DbConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, first_name, last_name, email FROM users WHERE is_member = TRUE")
	if err != nil {
		log.Fatalln("Database query error:", err)
	}

	var userList []models.User
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
		userList = append(userList, user)

	}

	if err := rows.Err(); err != nil {
		log.Fatalln("Error iterating rows:", err)
	}

	log.Println(userList)

	// get all users from the databasse where is_member = true

	// start new go func for each user

	// var jobs []models.Job
	// scrapeJobData(&jobs)

}
