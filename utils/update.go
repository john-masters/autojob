package utils

import (
	"autojob/models"
	"log"
	"sync"
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

	var wg sync.WaitGroup
	for _, user := range userList {
		wg.Add(1)
		go func(user models.User) {
			defer wg.Done()
			processUser(user)
		}(user)
	}
	wg.Wait()
}
