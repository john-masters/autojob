package utils

import (
	"autojob/db"
	"autojob/models"
	"log"
	"sync"
)

func UpdateToApplyList() {
	log.Println("updating...")

	var userList []models.User

	err := db.SelectMemberUsersByID(&userList)
	if err != nil {
		log.Fatal(err)
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
