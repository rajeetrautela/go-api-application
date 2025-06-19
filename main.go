package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"go-jwt-api/middleware"

	"go-jwt-api/config"
	"go-jwt-api/db"
	"go-jwt-api/migrations"
	"go-jwt-api/scheduler"

	"go-jwt-api/api"
	"go-jwt-api/helper"

	"github.com/gorilla/mux"
)

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
	router.HandleFunc("/", helper.FormHandler)
	router.HandleFunc("/login", api.Login).Methods("POST")
	router.HandleFunc("/refresh", api.Refresh).Methods("POST")
	router.HandleFunc("/logout", api.Logout).Methods("POST")
	router.HandleFunc("/items", middleware.JWTMiddleware(api.GetItems, "admin")).Methods("GET")
	router.HandleFunc("/items/{id}", middleware.JWTMiddleware(api.GetItem, "admin", "user")).Methods("GET")
	router.HandleFunc("/items", middleware.JWTMiddleware(api.CreateItem, "admin")).Methods("POST")

	router.HandleFunc("/register", api.Register).Methods("POST")

	router.HandleFunc("/items/{id}", middleware.JWTMiddleware(api.UpdateItem, "admin")).Methods("PUT")
	router.HandleFunc("/items/{id}", middleware.JWTMiddleware(api.DeleteItem, "admin")).Methods("DELETE")
	router.HandleFunc("/startworker", middleware.JWTMiddleware(helper.TriggerWorker, "admin", "user")).Methods("POST")
	router.HandleFunc("/upload", helper.UploadHandler)

	log.Println("Http Server started at :8001")
	log.Fatal(http.ListenAndServe(":8001", router))
}

func main() {
	go startHTTPServer()
	go config.StartGRPCServer()

	select {}
}
