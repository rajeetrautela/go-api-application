package helper_test

import (
	"bytes"
	"fmt"
	"go-jwt-api/helper"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestFormHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	helper.FormHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if ct := resp.Header.Get("Content-Type"); ct != "text/html" {
		t.Errorf("expected Content-Type text/html; got %s", ct)
	}

	body, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(body, []byte("<form")) {
		t.Errorf("expected body to contain <form>")
	}
}

func TestUploadHandler(t *testing.T) {
	// Mock the UploadFileToGRPCServer function
	orig := helper.UploadFileToGRPCServer
	helper.UploadFileToGRPCServer = func(path string) (string, error) {
		return "mocked upload successful", nil
	}
	defer func() { helper.UploadFileToGRPCServer = orig }()

	// Prepare a buffer to hold multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create a form file field
	part, err := writer.CreateFormFile("uploadFile", "testfile.txt")
	if err != nil {
		t.Fatalf("CreateFormFile error: %v", err)
	}

	// Write file content to the form file field
	fileContent := "This is a test file"
	_, err = io.Copy(part, strings.NewReader(fileContent))
	if err != nil {
		t.Fatalf("io.Copy error: %v", err)
	}

	writer.Close()

	// Create a POST request with the multipart form data
	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the handler
	helper.UploadHandler(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status 200 OK but got %d", status)
	}

	// Check the response body contains the mocked message
	expected := "gRPC File upload successful: mocked upload successful"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("Expected response body to contain %q, got %q", expected, rr.Body.String())
	}

	// clean ups
	os.Remove("./uploads/testfile.txt")

}

func TestUploadHandler_WrongMethod(t *testing.T) {
	req := httptest.NewRequest("GET", "/upload", nil)
	w := httptest.NewRecorder()

	helper.UploadHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405; got %d", resp.StatusCode)
	}
}

func TestUploadHandler_BadForm(t *testing.T) {
	req := httptest.NewRequest("POST", "/upload", strings.NewReader("bad data"))
	req.Header.Set("Content-Type", "multipart/form-data; boundary=xyz")
	w := httptest.NewRecorder()

	helper.UploadHandler(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400; got %d", resp.StatusCode)
	}
}

func TestUploadHandler_GrpcUploadFail(t *testing.T) {
	orig := helper.UploadFileToGRPCServer
	helper.UploadFileToGRPCServer = func(path string) (string, error) {
		return "", fmt.Errorf("mocked gRPC upload failure")
	}
	defer func() { helper.UploadFileToGRPCServer = orig }()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("uploadFile", "failfile.txt")
	part.Write([]byte("content"))
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	helper.UploadHandler(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500; got %d", resp.StatusCode)
	}

	os.Remove("./uploads/failfile.txt")
}
