package repository

import (
	"go-jwt-api/config"
	"go-jwt-api/model"
)

func StoreRefreshToken(token, username string) error {
	rt := model.RefreshToken{
		Token:    token,
		Username: username,
	}
	return config.DB.Create(&rt).Error
}

func ValidateRefreshToken(token, username string) bool {
	var rt model.RefreshToken
	err := config.DB.Where("token = ? AND username = ?", token, username).First(&rt).Error
	return err == nil
}

func DeleteRefreshToken(token string) error {
	return config.DB.Where("token = ?", token).Delete(&model.RefreshToken{}).Error
}
