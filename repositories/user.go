package repositories

import (
	"errors"
	"jaga/config"
	"jaga/models"
)

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := config.DB.First(&user, "email = ?", email)
	if result.Error != nil {
		return nil, errors.New("user not found")

	}
	return &user, nil
}

func CreateUser(user *models.User) error {
	if err := config.DB.Create(user).Error; err != nil {
		return errors.New("failed to create user")
	}
	return nil
}
