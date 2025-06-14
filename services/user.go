package services

import (
	"errors"
	"fmt"
	"log"
	"os"

	"jaga/dto"
	"jaga/models"
	"jaga/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func SeedSuperUser() {
	email := os.Getenv("SUPERUSER_EMAIL")
	password := os.Getenv("SUPERUSER_PASSWORD")
	name := os.Getenv("SUPERUSER_NAME")

	existing, err := repositories.GetUserByEmail(email)
	if err == nil && existing != nil {
		return
	}

	if (email == "" || password == "" || name == "") && existing == nil {
		log.Fatal("SUPERUSER credentials are missing. Please set SUPERUSER_EMAIL, SUPERUSER_PASSWORD, and SUPERUSER_NAME in .env.")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password")
	}

	superUser := models.User{
		ID:           uuid.New().String(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         "super_user",
	}

	if err := repositories.CreateUser(&superUser); err != nil {
		log.Fatal("Failed to create super user")
	}

	log.Println("Super user created")
}

func CreateUser(input dto.CreateUserRequest, userRole string) error {
	isAllowed := false
	switch userRole {
	case "super_user":
		isAllowed = input.Role == "admin" || input.Role == "manager" || input.Role == "technician"
	case "admin":
		isAllowed = input.Role == "manager" || input.Role == "technician"
	}
	if !isAllowed {
		return fmt.Errorf("%s is not allowed to create %s", userRole, input.Role)
	}

	existingUser, err := repositories.GetUserByEmail(input.Email)
	if err == nil && existingUser.ID != "" {
		return errors.New("email is already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user := models.User{
		ID:           uuid.New().String(),
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		Role:         input.Role,
	}

	return repositories.CreateUser(&user)
}
