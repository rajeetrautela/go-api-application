package tokenstore_test

import (
	"go-jwt-api/config"
	"go-jwt-api/test/test_util"
	"go-jwt-api/tokenstore"
	"testing"
)

func TestMain(m *testing.M) {
	test_util.TestMain(m)
}
func TestStoreValidateDelete(t *testing.T) {
	// Clean up before and after test
	cleanup := func() {
		config.DB.Exec("DELETE FROM refresh_tokens WHERE username = ?", "testuser")
	}
	cleanup()
	defer cleanup()

	token := "testtoken123"
	username := "testuser"

	// Test Store
	if err := tokenstore.Store(token, username); err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	// Test Validate
	gotUsername, err := tokenstore.Validate(token)
	if err != nil {
		t.Fatalf("Validate failed: %v", err)
	}
	if gotUsername != username {
		t.Errorf("Validate returned username %q; want %q", gotUsername, username)
	}

	// Test Delete
	if err := tokenstore.Delete(token); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Validate again after deletion should error
	_, err = tokenstore.Validate(token)
	if err == nil {
		t.Errorf("Validate succeeded after deletion, expected error")
	}
}
