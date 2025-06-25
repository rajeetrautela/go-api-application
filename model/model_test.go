package model_test

import (
	"testing"

	"go-jwt-api/config"
	"go-jwt-api/model"
	"go-jwt-api/test/test_util"
)

func TestMain(m *testing.M) {
	// Load environment variables from .env file
	test_util.TestMain(m)
}

func TestItemModel(t *testing.T) {
	db := config.DB

	item := model.Item{Name: "TestItem", Price: 100}
	if err := db.Create(&item).Error; err != nil {
		t.Fatalf("Failed to create item: %v", err)
	}

	var found model.Item
	if err := db.First(&found, item.ID).Error; err != nil {
		t.Fatalf("Failed to retrieve item: %v", err)
	}

	if found.Name != item.Name || found.Price != item.Price {
		t.Errorf("Item mismatch: expected %+v, got %+v", item, found)
	}
}

func TestUserModel(t *testing.T) {
	db := config.DB

	user := model.User{
		Username: "testuser",
		Password: "hashed-password",
		Role:     "admin",
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	var found model.User
	if err := db.First(&found, user.ID).Error; err != nil {
		t.Fatalf("Failed to retrieve user: %v", err)
	}

	if found.Username != user.Username || found.Role != user.Role {
		t.Errorf("User mismatch: expected %+v, got %+v", user, found)
	}
}

func TestRefreshTokenModel(t *testing.T) {
	db := config.DB

	token := model.RefreshToken{
		Token:    "some-refresh-token",
		Username: "testuser",
	}
	if err := db.Create(&token).Error; err != nil {
		t.Fatalf("Failed to create refresh token: %v", err)
	}

	var found model.RefreshToken
	if err := db.First(&found, token.ID).Error; err != nil {
		t.Fatalf("Failed to retrieve refresh token: %v", err)
	}

	if found.Token != token.Token || found.Username != token.Username {
		t.Errorf("Token mismatch: expected %+v, got %+v", token, found)
	}
}
