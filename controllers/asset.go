package controllers

import (
	"encoding/json"
	"jaga/config"
	"jaga/dto"
	"jaga/middleware"
	"jaga/models"
	"jaga/services"
	"net/http"
)

func CreateAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var input dto.CreateAsset

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	userID := ctx.Value(middleware.ContextUserID).(string)
	var user models.User
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if err := services.CreateAsset(input, user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Asset created"))
}
