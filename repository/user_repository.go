package repository

import (
	"go-jwt-api/config"
	"go-jwt-api/model"
)

func CreateUser(user *model.User) error {
	return config.DB.Create(user).Error
}

func GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := config.DB.First(&user, id).Error
	return &user, err
}

func UpdateUser(user *model.User) error {
	return config.DB.Save(user).Error
}

func DeleteUser(id uint) error {
	return config.DB.Delete(&model.User{}, id).Error
}
