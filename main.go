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
)

func startHTTPServer() {
	go scheduler.StartCronJobs()
	config.ConnectDatabase()

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

	log.Println("Http Server started at :8001")
	log.Fatal(http.ListenAndServe(":8001", router))
}

func main() {
	go startHTTPServer()
	go config.StartGRPCServer()
	select {}
}
