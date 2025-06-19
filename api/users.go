package api

import (
	"encoding/json"
	"go-jwt-api/auth"
	"go-jwt-api/config"
	"go-jwt-api/model"
	"go-jwt-api/repository"
	"go-jwt-api/tokenstore"
	"net/http"
)

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
