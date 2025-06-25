package middleware_test

import (
	"fmt"
	"go-jwt-api/auth"
	"go-jwt-api/middleware"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestJWTMiddleware_ValidTokenAndRole(t *testing.T) {
	token, err := auth.GenerateJWT("testuser", "admin")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", token)

	rr := httptest.NewRecorder()
	called := false

	handler := middleware.JWTMiddleware(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}, "admin")

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}
	if !called {
		t.Error("Expected handler to be called")
	}
}

func TestJWTMiddleware_MissingToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler := middleware.JWTMiddleware(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called without token")
	}, "admin")

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Missing token") {
		t.Errorf("Unexpected body: %s", rr.Body.String())
	}
}

func TestJWTMiddleware_InvalidRole(t *testing.T) {
	token, _ := auth.GenerateJWT("testuser", "user")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", token)
	rr := httptest.NewRecorder()

	handler := middleware.JWTMiddleware(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called with wrong role")
	}, "admin")

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("Expected 403, got %d", rr.Code)
	}
}

func TestLoggingMiddleware(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/logtest", nil)
	rr := httptest.NewRecorder()

	handler := middleware.LoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}
}

func TestRecoveryMiddleware(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler := middleware.RecoveryMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("something went wrong")
	}))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "Internal Server Error") {
		t.Errorf("Unexpected body: %s", rr.Body.String())
	}
}

func TestTimeoutMiddleware(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler := middleware.TimeoutMiddleware(50 * time.Millisecond)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		fmt.Fprintln(w, "Done")
	}))

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected 503 due to timeout, got %d", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "⏱ Request timed out") {
		t.Errorf("Unexpected timeout message: %s", rr.Body.String())
	}
}
