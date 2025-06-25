package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-jwt-api/api"
	"go-jwt-api/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	"go-jwt-api/helper"
	"go-jwt-api/test/test_util"
)

var router *mux.Router

func TestMain(m *testing.M) {
	// Load environment variables from .env file
	test_util.TestMain(m)
}

func TestSetup(t *testing.T) {
	if router == nil {
		router = mux.NewRouter()
		router.HandleFunc("/items", MockJWTMiddleware(api.CreateItem, "admin")).Methods("POST")
		router.HandleFunc("/items", MockJWTMiddleware(api.GetItems, "admin")).Methods("GET")
		router.HandleFunc("/items/{id}", MockJWTMiddleware(api.GetItem, "admin", "user")).Methods("GET")
		router.HandleFunc("/items/{id}", MockJWTMiddleware(api.UpdateItem, "admin")).Methods("PUT")
		router.HandleFunc("/items/{id}", MockJWTMiddleware(api.DeleteItem, "admin")).Methods("DELETE")

		router.HandleFunc("/login", api.Login).Methods("POST")
		router.HandleFunc("/refresh", api.Refresh).Methods("POST")
		router.HandleFunc("/logout", api.Logout).Methods("POST")
		router.HandleFunc("/register", api.Register).Methods("POST")
		router.HandleFunc("/upload", helper.UploadHandler)
		router.HandleFunc("/startworker", MockJWTMiddleware(helper.TriggerWorker, "admin", "user")).Methods("POST")
	}
}

func MockJWTMiddleware(next http.HandlerFunc, _ ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Bypass auth checks for testing
		next.ServeHTTP(w, r)
	}
}

// helper for JSON requests
func makeJSONRequest(t *testing.T, method, url string, body interface{}) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("Failed to encode JSON: %v", err)
		}
	}
	req := httptest.NewRequest(method, url, &buf)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func TestItemCRUD(t *testing.T) {
	TestSetup(t)

	var created model.Item

	t.Run("Create", func(t *testing.T) {
		item := model.Item{Name: "CRUD Item", Price: 100}
		rr := makeJSONRequest(t, "POST", "/items", item)

		if rr.Code != http.StatusOK {
			t.Fatalf("Create failed: expected 200, got %d", rr.Code)
		}

		if err := json.NewDecoder(rr.Body).Decode(&created); err != nil {
			t.Fatalf("Failed to decode create response: %v", err)
		}
	})

	t.Run("Find", func(t *testing.T) {
		rr := makeJSONRequest(t, "GET", fmt.Sprintf("/items/%d", created.ID), nil)
		if rr.Code != http.StatusOK {
			t.Fatalf("Get failed: expected 200, got %d", rr.Code)
		}
		var fetched model.Item
		if err := json.NewDecoder(rr.Body).Decode(&fetched); err != nil {
			t.Fatalf("Failed to decode get response: %v", err)
		}
		if fetched.ID != created.ID {
			t.Errorf("Expected ID %d, got %d", created.ID, fetched.ID)
		}
	})

	t.Run("Update", func(t *testing.T) {
		updated := model.Item{Name: "Updated Item", Price: 999}
		rr := makeJSONRequest(t, "PUT", fmt.Sprintf("/items/%d", created.ID), updated)

		if rr.Code != http.StatusOK {
			t.Fatalf("Update failed: expected 200, got %d", rr.Code)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		rr := makeJSONRequest(t, "DELETE", fmt.Sprintf("/items/%d", created.ID), nil)
		if rr.Code != http.StatusNoContent {
			t.Fatalf("Delete failed: expected 204, got %d", rr.Code)
		}
	})

	t.Run("Verify deletion", func(t *testing.T) {
		rr := makeJSONRequest(t, "GET", fmt.Sprintf("/items/%d", created.ID), nil)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("Expected 404 after deletion, got %d", rr.Code)
		}
	})
}

func TestLoginFailures(t *testing.T) {
	// Invalid JSON
	body := strings.NewReader("invalid json")
	req := httptest.NewRequest("POST", "/login", body)
	rr := httptest.NewRecorder()
	api.Login(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid JSON, got %d", rr.Code)
	}

	// Missing user in DB (simulate by using a non-existent username)
	user := model.User{Username: "nonexistent", Password: "pass"}
	b, _ := json.Marshal(user)
	req = httptest.NewRequest("POST", "/login", bytes.NewReader(b))
	rr = httptest.NewRecorder()
	api.Login(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for unknown user, got %d", rr.Code)
	}
}

func TestRegisterFailures(t *testing.T) {
	// Invalid JSON
	req := httptest.NewRequest("POST", "/register", strings.NewReader("invalid"))
	rr := httptest.NewRecorder()
	api.Register(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid JSON, got %d", rr.Code)
	}

	// Attempt to register user with existing username or invalid data
	// You can mock repository.CreateUser to return error here to simulate failure.
}

func TestRefreshFailures(t *testing.T) {
	// Invalid JSON
	req := httptest.NewRequest("POST", "/refresh", strings.NewReader("invalid"))
	rr := httptest.NewRecorder()
	api.Refresh(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid JSON, got %d", rr.Code)
	}

	// Invalid token
	b := []byte(`{"refresh_token":"invalidtoken"}`)
	req = httptest.NewRequest("POST", "/refresh", bytes.NewReader(b))
	rr = httptest.NewRecorder()
	api.Refresh(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for invalid token, got %d", rr.Code)
	}
}

func TestLogoutFailures(t *testing.T) {
	// Invalid JSON
	req := httptest.NewRequest("POST", "/logout", strings.NewReader("invalid"))
	rr := httptest.NewRecorder()
	api.Logout(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid JSON, got %d", rr.Code)
	}
}

func TestGetItemFailures(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/items/{id}", api.GetItem).Methods("GET")

	// invalid id (non-int)
	req := httptest.NewRequest("GET", "/items/abc", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid ID, got %d", rr.Code)
	}

	// non-existent id (valid int but missing item)
	req = httptest.NewRequest("GET", "/items/9999999", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected 404 for missing item, got %d", rr.Code)
	}
}

func TestCreateItemFailures(t *testing.T) {
	// Invalid JSON
	req := httptest.NewRequest("POST", "/items", strings.NewReader("invalid"))
	rr := httptest.NewRecorder()
	api.CreateItem(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid JSON, got %d", rr.Code)
	}

	// Simulate repository.CreateItem failure (mock repository.CreateItem to return error)
	// For now, assume normal DB errors handled via repo layer.
}

func TestUpdateItemFailures(t *testing.T) {
	// Invalid ID param
	req := httptest.NewRequest("PUT", "/items/abc", nil)
	rr := httptest.NewRecorder()
	api.UpdateItem(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid ID, got %d", rr.Code)
	}

	// Invalid JSON body
	req = httptest.NewRequest("PUT", "/items/1", strings.NewReader("invalid"))
	rr = httptest.NewRecorder()
	api.UpdateItem(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid JSON, got %d", rr.Code)
	}

	// Simulate repository.UpdateItem failure by mocking to return error
}

func TestDeleteItemFailures(t *testing.T) {
	// Invalid ID param
	req := httptest.NewRequest("DELETE", "/items/abc", nil)
	rr := httptest.NewRecorder()
	api.DeleteItem(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid ID, got %d", rr.Code)
	}

	// Simulate repository.DeleteItem failure by mocking to return error
}
