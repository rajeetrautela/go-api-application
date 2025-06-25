package main

import (
	"fmt"
	"log"
	"net/http"

	"go-jwt-api/config"
	"go-jwt-api/db"
	"go-jwt-api/migrations"
	"go-jwt-api/routes"
	"go-jwt-api/scheduler"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func startHTTPServer() {
	loadEnv()
	go scheduler.StartCronJobs()
	config.ConnectDatabase("DEV")

	if config.DB != nil {
		log.Println("✅ Successfully connected to the database!")
		migrations.Migrate()
		fmt.Println("📦 Database migrated successfully too Hurray!")
		fmt.Println("📦 Now Seeding!!!!")
		db.Seed()
	} else {
		log.Fatal("❌ Failed to connect to the database.")
	}

	router := mux.NewRouter()
	routes.RegisterRoutes(router)

	log.Println("Http Server started at :8085")

	log.Fatal(http.ListenAndServe("0.0.0.0:8085", router))
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
func main() {
	go startHTTPServer()
	go config.StartGRPCServer()
	select {}
}
