package repository

import (
	"errors"
	"go-jwt-api/auth"
	"go-jwt-api/config"
	"go-jwt-api/model"

	"gorm.io/gorm"
)

func CreateUser(user *model.User) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var existing model.User
		if err := tx.Where("username = ?", user.Username).First(&existing).Error; err == nil {
			return errors.New("username already exists")
		}

		hashedPassword, err := auth.HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword

		return tx.Create(user).Error
	})
}

func GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := config.DB.First(&user, id).Error
	return &user, err
}

func UpdateUser(user *model.User) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Save(user).Error
	})
}

func DeleteUser(id uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&model.User{}, id).Error
	})
}
