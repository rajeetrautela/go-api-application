package config_test

import (
	"go-jwt-api/config"
	"os"
	"testing"
)

func TestConnectDatabase_DefaultDSN(t *testing.T) {
	// Clear env vars to trigger default DSN path
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_PORT")

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("ConnectDatabase panicked with default DSN: %v", r)
		}
	}()

	config.ConnectDatabase("DEV")
	if config.DB == nil {
		t.Fatal("Expected DB to be initialized with default DSN, got nil")
	}
}

func TestConnectDatabase_SkipConnection(t *testing.T) {
	config.SkipDBConnect = true
	config.ConnectDatabase("TEST")
	if config.DB != nil {
		t.Fatal("Expected DB to be nil when skipping connection")
	}
}
