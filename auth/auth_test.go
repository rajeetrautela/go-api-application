package auth_test

import (
	"strings"
	"testing"
	"time"

	"go-jwt-api/auth"
)

func TestGenerateAndValidateJWT(t *testing.T) {
	username := "testuser"
	role := "admin"

	token, err := auth.GenerateJWT(username, role)
	if err != nil {
		t.Fatalf("GenerateJWT failed: %v", err)
	}

	if token == "" {
		t.Fatal("JWT token is empty")
	}

	claims, err := auth.ValidateJWT(token)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}

	if claims.Username != username || claims.Role != role {
		t.Errorf("Claims mismatch: got %+v", claims)
	}

	if time.Until(claims.ExpiresAt.Time) > 15*time.Minute {
		t.Error("JWT expiration time is not approximately 15 minutes")
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	username := "testuser"

	token, err := auth.GenerateRefreshToken(username)
	if err != nil {
		t.Fatalf("GenerateRefreshToken failed: %v", err)
	}

	if token == "" {
		t.Fatal("Refresh token is empty")
	}

	claims, err := auth.ValidateJWT(token)
	if err != nil {
		t.Fatalf("Refresh token validation failed: %v", err)
	}

	if claims.Username != username {
		t.Errorf("Refresh token username mismatch. Expected %s, got %s", username, claims.Username)
	}

	if time.Until(claims.ExpiresAt.Time) < 23*time.Hour {
		t.Error("Refresh token expiration time is less than expected")
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	_, err := auth.ValidateJWT("invalid.token.string")
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
}

func TestHashAndCheckPassword(t *testing.T) {
	password := "mySecurePassword123"

	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash == "" || hash == password {
		t.Error("Generated hash is invalid or same as plain password")
	}

	if !auth.CheckPasswordHash(password, hash) {
		t.Error("CheckPasswordHash failed to validate correct password")
	}
}

func TestCheckPasswordHash_Invalid(t *testing.T) {
	password := "myPassword"
	wrong := "wrongPassword"

	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if auth.CheckPasswordHash(wrong, hash) {
		t.Error("Expected CheckPasswordHash to return false for wrong password")
	}
}

func TestHashPassword_Strength(t *testing.T) {
	password := "test"
	hash1, _ := auth.HashPassword(password)
	hash2, _ := auth.HashPassword(password)

	if hash1 == hash2 {
		t.Error("Expected different hashes for same password due to salting")
	}

	if !strings.HasPrefix(hash1, "$2a$") && !strings.HasPrefix(hash1, "$2b$") {
		t.Error("Unexpected hash prefix, bcrypt expected")
	}
}
