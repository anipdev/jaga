package services

import (
	"errors"

	"jaga/models"
	"jaga/repositories"
	"jaga/utils"

	"gorm.io/gorm"
)

type UserService interface {
	GetUsers(page, itemsPerPage int, sortBy, sortDir string) ([]models.User, int64, error)
	GetUserByID(userID string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User, requesterRole string) error
	UpdateUser(user *models.User, requesterRole string) error
	DeleteUser(userID, requesterRole string) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetUsers(page, itemsPerPage int, sortBy, sortDir string) ([]models.User, int64, error) {
	users, totalItems, err := s.userRepo.GetUsers(page, itemsPerPage, sortBy, sortDir)
	if err != nil {
		return nil, 0, err
	}
	return users, totalItems, nil
}

func (s *userService) GetUserByID(userID string) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) CreateUser(user *models.User, requesterRole string) error {
	if requesterRole == "admin" && user.Role == "admin" {
		return errors.New("admin cannot create a user with 'admin' role")
	}

	_, err := s.userRepo.GetUserByEmail(user.Email)
	if err == nil {
		return errors.New("email already registered")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err // Other database error
	}

	if user.ID == "" {
		user.ID = utils.GenerateUUID()
	}

	hashedPassword, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user.PasswordHash = hashedPassword

	return s.userRepo.CreateUser(user)
}

func (s *userService) UpdateUser(user *models.User, requesterRole string) error {
	existingUser, err := s.userRepo.GetUserByID(user.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	if requesterRole == "admin" && user.Role == "admin" {
		return errors.New("admin cannot update a user to 'admin' role")
	}

	if user.PasswordHash != "" {
		hashedPassword, err := utils.HashPassword(user.PasswordHash)
		if err != nil {
			return errors.New("failed to hash password")
		}
		user.PasswordHash = hashedPassword
	} else {
		user.PasswordHash = existingUser.PasswordHash
	}

	return s.userRepo.UpdateUser(user)
}

func (s *userService) DeleteUser(userID, requesterRole string) error {
	existingUser, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	if requesterRole == "admin" && existingUser.Role == "admin" {
		return errors.New("admin cannot delete another admin")
	}

	return s.userRepo.DeleteUser(userID)
}
