package routes

import (
	"go-jwt-api/api"
	"go-jwt-api/helper"
	"go-jwt-api/middleware"

	"time"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	// Global middlewares
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.RecoveryMiddleware)
	router.Use(middleware.TimeoutMiddleware(5 * time.Second))

	// Public routes
	router.HandleFunc("/", helper.FormHandler)
	router.HandleFunc("/login", api.Login).Methods("POST")
	router.HandleFunc("/refresh", api.Refresh).Methods("POST")
	router.HandleFunc("/logout", api.Logout).Methods("POST")
	router.HandleFunc("/register", api.Register).Methods("POST")
	router.HandleFunc("/upload", helper.UploadHandler)

	// Protected routes
	router.HandleFunc("/items", middleware.JWTMiddleware(api.GetItems, "admin")).Methods("GET")
	router.HandleFunc("/items/{id}", middleware.JWTMiddleware(api.GetItem, "admin", "user")).Methods("GET")
	router.HandleFunc("/items", middleware.JWTMiddleware(api.CreateItem, "admin")).Methods("POST")
	router.HandleFunc("/items/{id}", middleware.JWTMiddleware(api.UpdateItem, "admin")).Methods("PUT")
	router.HandleFunc("/items/{id}", middleware.JWTMiddleware(api.DeleteItem, "admin")).Methods("DELETE")
	router.HandleFunc("/startworker", middleware.JWTMiddleware(helper.TriggerWorker, "admin", "user")).Methods("POST")
}
