package helpers

import (
	"log"

	"real_time/backend/config"
)

func GetUsername(userId int) (string) {
	var username string
	err := config.Db.QueryRow("SELECT username FROM users WHERE id = ?", userId).Scan(&username)
	if err != nil {
		log.Printf("Error fetching sender username: %v", err)
		return ""
	}
	return username
}
