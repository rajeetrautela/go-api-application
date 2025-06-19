package migrations

import (
	"go-jwt-api/config"
	"go-jwt-api/model"
)

func Migrate() {
	config.DB.AutoMigrate(&model.Item{})
	config.DB.AutoMigrate(&model.User{})
	config.DB.AutoMigrate(&model.RefreshToken{})
}
