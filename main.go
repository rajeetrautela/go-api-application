package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"go-jwt-api/auth"
	"go-jwt-api/middleware"
	"go-jwt-api/model"
	"go-jwt-api/repository"

	"go-jwt-api/config"
	"go-jwt-api/db"
	"go-jwt-api/migrations"
	"go-jwt-api/scheduler"
	"go-jwt-api/tokenstore"
	"go-jwt-api/worker"

	"github.com/gorilla/mux"
)

func getItems(w http.ResponseWriter, r *http.Request) {
	var items []model.Item
	if err := config.DB.Find(&items).Error; err != nil {
		http.Error(w, "Failed to fetch items", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(items)
}

func getItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	item, err := repository.GetItemByID(uint(id))
	if err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(item)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	var item model.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := repository.CreateItem(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(item)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updated model.Item
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updated.ID = uint(id)

	if err := repository.UpdateItem(&updated); err != nil {
		http.Error(w, "Failed to update item", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(updated)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := repository.DeleteItem(uint(id)); err != nil {
		http.Error(w, "Failed to delete item", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// it log's in user and provide a valid JWT in return
func login(w http.ResponseWriter, r *http.Request) {
	var creds model.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var user model.User
	if err := config.DB.Where("username = ?", creds.Username).First(&user).Error; err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if !auth.CheckPasswordHash(creds.Password, user.Password) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, _ := auth.GenerateJWT(user.Username, user.Role)
	refreshToken, _ := auth.GenerateRefreshToken(user.Username)
	tokenstore.Store(refreshToken, user.Username)

	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  token,
		"refresh_token": refreshToken,
	})
}

func register(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := repository.CreateUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	username, err := tokenstore.Validate(req.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	token, _ := auth.GenerateJWT(username, "user")
	json.NewEncoder(w).Encode(map[string]string{"access_token": token})
}

func logout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	tokenstore.Delete(req.RefreshToken)
	w.WriteHeader(http.StatusOK)
}

func triggerWorker(w http.ResponseWriter, r *http.Request) {
	// Default values
	const defNumJobs = 5
	const defNumWorkers = 3

	var req struct {
		JobCount    int `json:"job_count"`
		WorkerCount int `json:"worker_count"`
	}

	numJobs := defNumJobs
	numWorkers := defNumWorkers

	// Try to decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("Invalid JSON. Using default values.")
	} else {
		// Use values from request if decoding succeeds
		if req.JobCount > 0 {
			numJobs = req.JobCount
		}
		if req.WorkerCount > 0 {
			numWorkers = req.WorkerCount
		}
	}

	// Dispatch jobs and start workers
	jobs := worker.DispatchJobs(numJobs)
	results := make(chan string, numJobs)

	for w := 1; w <= numWorkers; w++ {
		go worker.StartWorker(r.Context(), w, jobs, results)
	}

	var output []string
	for a := 1; a <= numJobs; a++ {
		output = append(output, <-results)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)

}

func formHandler(w http.ResponseWriter, r *http.Request) {
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

func uploadHandler(w http.ResponseWriter, r *http.Request) {
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
	message, err := UploadFileToGRPCServer(tempPath)
	if err != nil {
		http.Error(w, "gRPC upload failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "gRPC File upload successful: %s\n", message)
}

// func register(w http.ResponseWriter, r *http.Request) {
// 	var creds user.User
// 	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
// 		http.Error(w, "Bad request", http.StatusBadRequest)
// 		return
// 	}
// 	// format it wrt to user struct and then validate and save into DB
// }

func startHTTPServer() {
	go scheduler.StartCronJobs()
	config.ConnectDatabase()
	// Test the connection
	if config.DB != nil {
		log.Println("✅ Successfully connected to the database!")
		migrations.Migrate()
		fmt.Println("📦 Database migrated successfully too Hurray!")
		fmt.Println("📦 Now Seeding!!!!")

		// Make sure if its not required you can just comment it here
		// as seed requires just one time to create some dummy data to proceed with
		db.Seed()
	} else {
		log.Fatal("❌ Failed to connect to the database.")
	}
	router := mux.NewRouter()

	// global middlewares
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.RecoveryMiddleware)
	router.Use(middleware.TimeoutMiddleware(5 * time.Second)) // 5s timeout
	router.HandleFunc("/", formHandler)
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/refresh", refresh).Methods("POST")
	router.HandleFunc("/logout", logout).Methods("POST")
	router.HandleFunc("/items", middleware.JWTMiddleware(getItems, "admin")).Methods("GET")
	router.HandleFunc("/items/{id}", middleware.JWTMiddleware(getItem, "admin", "user")).Methods("GET")
	router.HandleFunc("/items", middleware.JWTMiddleware(createItem, "admin")).Methods("POST")

	router.HandleFunc("/register", register).Methods("POST")

	router.HandleFunc("/items/{id}", middleware.JWTMiddleware(updateItem, "admin")).Methods("PUT")
	router.HandleFunc("/items/{id}", middleware.JWTMiddleware(deleteItem, "admin")).Methods("DELETE")
	router.HandleFunc("/startworker", middleware.JWTMiddleware(triggerWorker, "admin", "user")).Methods("POST")
	router.HandleFunc("/upload", uploadHandler)

	log.Println("Http Server started at :8001")
	log.Fatal(http.ListenAndServe(":8001", router))
}

func main() {
	go startHTTPServer()
	go startGRPCServer()

	select {}
}
