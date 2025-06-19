package helper

import (
	"fmt"
	"go-jwt-api/config"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func FormHandler(w http.ResponseWriter, r *http.Request) {
	html := `
 <!DOCTYPE html>
 <html>
 <head><title>Upload File</title></head>
 <body>
 <h2>Upload a File</h2>
 <form enctype="multipart/form-data" action="/upload" method="post">
 <input type="file" name="uploadFile" />
 <input type="submit" value="Upload" />
 </form>
 </body>
 </html>
 `
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("uploadFile")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save to temp file
	tempPath := filepath.Join(os.TempDir(), handler.Filename)
	tempFile, err := os.Create(tempPath)
	if err != nil {
		http.Error(w, "Error creating temp file", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()
	io.Copy(tempFile, file)

	// Call gRPC client
	message, err := config.UploadFileToGRPCServer(tempPath)
	if err != nil {
		http.Error(w, "gRPC upload failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "gRPC File upload successful: %s\n", message)
}
