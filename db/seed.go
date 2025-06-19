package db

import (
	"go-jwt-api/auth"
	"go-jwt-api/config"
	"go-jwt-api/model"
	"log"
)

func Seed() {
	// seeding Users
	users := []model.User{
		{Username: "testadminuser", Password: "admin123", Role: "admin"},
		{Username: "testuser", Password: "user123", Role: "user"},
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

	// Seeding Items

	items := []model.Item{
		{Name: "Laptop", Price: 75000},
		{Name: "Smartphone", Price: 30000},
		{Name: "Headphones", Price: 5000},
		{Name: "Monitor", Price: 15000},
	}

	for _, item := range items {
		if err := config.DB.Create(&item).Error; err != nil {
			log.Printf("Failed to seed item %s: %v", item.Name, err)
		} else {
			log.Printf("Seeded item %s", item.Name)
		}
	}

}
