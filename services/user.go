package services

import (
	"errors"

	"jaga/models"
	"jaga/repositories"
	"jaga/utils"

	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(user *models.User, creatorRole string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.GetUserByEmail(email)
}

func (s *userService) CreateUser(user *models.User, creatorRole string) (*models.User, error) {
	if creatorRole == "admin" && user.Role == "admin" {
		return nil, errors.New("admin cannot create a user with 'admin' role")
	}

	_, err := s.userRepo.GetUserByEmail(user.Email)
	if err == nil {
		return nil, errors.New("email already registered")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}
	user.PasswordHash = hashedPassword

	return s.userRepo.CreateUser(user)
}
