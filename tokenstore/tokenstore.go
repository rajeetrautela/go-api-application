package tokenstore

import (
	"go-jwt-api/config"
	"go-jwt-api/model"
)

func Store(token, username string) error {
	rt := model.RefreshToken{
		Token:    token,
		Username: username,
	}
	return config.DB.Create(&rt).Error
}

func Validate(token string) (string, error) {
	var rt model.RefreshToken
	err := config.DB.Where("token = ?", token).First(&rt).Error
	if err != nil {
		return "", err
	}
	return rt.Username, nil
}

func Delete(token string) error {
	return config.DB.Where("token = ?", token).Delete(&model.RefreshToken{}).Error
}
