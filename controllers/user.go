package controllers

import (
	"encoding/json"
	"net/http"

	"jaga/dto"
	"jaga/middleware"
	"jaga/services"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var input dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	userRole := ctx.Value(middleware.ContextRole).(string)

	if err := services.CreateUser(input, userRole); err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created"))
}
