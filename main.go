package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"jaga/config"
	"jaga/models"
	"jaga/repositories"
	"jaga/routes"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, proceeding with system env")
	}
}

func SeedSuperUser(userRepo repositories.UserRepository) {
	email := os.Getenv("SUPERUSER_EMAIL")
	password := os.Getenv("SUPERUSER_PASSWORD")
	name := os.Getenv("SUPERUSER_NAME")

	if email == "" || password == "" || name == "" {
		log.Println("SUPERUSER credentials (SUPERUSER_EMAIL, SUPERUSER_PASSWORD, SUPERUSER_NAME) are not fully set in .env. Skipping super user seeding.")
		return
	}

	existingUser, err := userRepo.GetUserByEmail(email)
	if err == nil && existingUser != nil {
		log.Println("Super user already exists.")
		return
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Fatalf("Error checking for existing super user: %v", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password for super user: %v", err)
	}

	superUser := models.User{
		ID:           uuid.New().String(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         "super_user",
	}

	_, err = userRepo.CreateUser(&superUser)
	if err != nil {
		log.Fatalf("Failed to create super user: %v", err)
	}

	log.Println("Super user created successfully!")
}

func main() {
	db := config.InitDB()
	config.AutoMigrate(db)

	userRepository := repositories.NewUserRepository(db)
	SeedSuperUser(userRepository)

	router := routes.RegisterRoutes()

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", router)

}
