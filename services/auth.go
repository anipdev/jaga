package services

import (
	"errors"
	"jaga/config"
	"jaga/models"

	"golang.org/x/crypto/bcrypt"
)

func Authenticate(email, password string) (*models.User, error) {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("incorrect username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("incorrect username or password")
	}

	return &user, nil
}
