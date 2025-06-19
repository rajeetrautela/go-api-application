package db

import (
	"go-jwt-api/auth"
	"go-jwt-api/config"
	"go-jwt-api/model"
	"log"
)

func Seed() {
	users := []model.User{
		{Username: "admin", Password: "admin123", Role: "admin"},
		{Username: "user1", Password: "user123", Role: "user"},
	}

	for _, user := range users {
		hashedPassword, err := auth.HashPassword(user.Password)
		if err != nil {
			log.Printf("Failed to hash password for user %s: %v", user.Username, err)
			continue
		}
		user.Password = hashedPassword

		if err := config.DB.Create(&user).Error; err != nil {
			log.Printf("Failed to seed user %s: %v", user.Username, err)
		} else {
			log.Printf("Seeded user %s", user.Username)
		}
	}
}
