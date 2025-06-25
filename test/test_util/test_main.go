package test_util

import (
	"fmt"
	"go-jwt-api/config"
	"go-jwt-api/migrations"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	loadEnv()
	config.ConnectDatabase("TEST")
	migrations.Migrate()
	cleanupTestDB()

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func cleanupTestDB() {
	config.DB.Exec("TRUNCATE TABLE items, users, refresh_tokens RESTART IDENTITY CASCADE;")
}

func loadEnv() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("failed to get working dir: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Working directory:", wd)

	// Fix the root path to point to go-api-application folder
	root := filepath.Join(wd, "..", ".env") // assuming wd is /.../go-api-application/test or similar
	fmt.Println("Loading .env from:", root)

	if err := godotenv.Load(root); err != nil {
		fmt.Printf("Failed to load .env from %s: %v\n", root, err)
		os.Exit(1)
	}
}
