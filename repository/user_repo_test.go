package repository_test

import (
	"go-jwt-api/config"
	"go-jwt-api/model"
	"go-jwt-api/repository"
	"testing"
)

func TestCreateUser(t *testing.T) {
	user := model.User{
		Username: "newuser",
		Password: "plainpassword",
		Role:     "admin",
	}
	err := repository.CreateUser(&user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	var stored model.User
	err = config.DB.First(&stored, user.ID).Error
	if err != nil {
		t.Fatalf("User not found in DB: %v", err)
	}
	if stored.Password == "plainpassword" {
		t.Error("Password was not hashed")
	}
}

func TestCreateUser_Duplicate(t *testing.T) {
	_ = repository.CreateUser(&model.User{
		Username: "dupeuser",
		Password: "pass",
		Role:     "admin",
	})

	err := repository.CreateUser(&model.User{
		Username: "dupeuser",
		Password: "anotherpass",
		Role:     "user",
	})
	if err == nil {
		t.Error("Expected error for duplicate username, got nil")
	}
}

func TestGetUserByID(t *testing.T) {
	user := model.User{
		Username: "getuser",
		Password: "pass",
		Role:     "user",
	}
	_ = repository.CreateUser(&user)

	fetched, err := repository.GetUserByID(user.ID)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}
	if fetched.Username != user.Username || fetched.Role != user.Role {
		t.Errorf("Fetched user mismatch: expected %+v, got %+v", user, fetched)
	}
}

func TestUpdateUser(t *testing.T) {
	user := model.User{
		Username: "updateuser",
		Password: "pass",
		Role:     "user",
	}
	_ = repository.CreateUser(&user)

	user.Role = "admin"
	err := repository.UpdateUser(&user)
	if err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}

	var updated model.User
	_ = config.DB.First(&updated, user.ID)
	if updated.Role != "admin" {
		t.Errorf("Expected role 'admin', got '%s'", updated.Role)
	}
}

func TestDeleteUser(t *testing.T) {
	user := model.User{
		Username: "tobedeleted",
		Password: "pass",
		Role:     "user",
	}
	_ = repository.CreateUser(&user)

	err := repository.DeleteUser(user.ID)
	if err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}

	var found model.User
	err = config.DB.First(&found, user.ID).Error
	if err == nil {
		t.Errorf("Expected user to be deleted, but found in DB")
	}
}
