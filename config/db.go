package config

import (
	"errors"
	"fmt"
	"jaga/models"
	"jaga/repositories"
	"log"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

var JWTSecret string
var JWTExpHours int

var DefaultSuperUserEmail = "adminjaga@gmail.com"
var DefaultSuperUserPassword = "@AdminJaga99"
var DefaultSuperUserName = "Super Admin"
var DefaultSuperUserRole = "super_user"

func InitDB() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	JWTSecret = os.Getenv("JWT_SECRET_KEY")
	if JWTSecret == "" {
		log.Fatal("JWT_SECRET_KEY environment variable not set. This is critical for JWT security.")
	}

	jwtExpStr := os.Getenv("JWT_EXPIRATION_HOURS")
	if jwtExpStr == "" {
		log.Println("JWT_EXPIRATION_HOURS environment variable not set. Defaulting to 24 hours.")
		JWTExpHours = 24 // Default value
	} else {
		hours, err := strconv.Atoi(jwtExpStr)
		if err != nil {
			log.Fatalf("Invalid value for JWT_EXPIRATION_HOURS: %v. Must be an integer.", err)
		}
		JWTExpHours = hours
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	fmt.Println("Successfully connected to database!")
	return DB
}

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.AssetCategory{},
		&models.Asset{},
		&models.MaintenanceSchedule{},
		&models.MaintenanceRecord{},
	)
	if err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}
	fmt.Println("Database migrated successfully!")
}

func SeedSuperUser(userRepo repositories.UserRepository) {

	log.Printf("Attempting to seed super user with default credentials: Email=%s, Name=%s", DefaultSuperUserEmail, DefaultSuperUserName)

	existingSuperUser, err := userRepo.GetUserByRole(DefaultSuperUserRole)
	if err == nil && existingSuperUser != nil {
		log.Printf("A user with role '%s' already exists (email: %s). Skipping seeding.", DefaultSuperUserRole, existingSuperUser.Email)
		return
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Fatalf("Error checking for existing user by role %s: %v", DefaultSuperUserRole, err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(DefaultSuperUserPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash default password for super user: %v", err)
	}

	superUser := models.User{
		ID:           uuid.New().String(),
		Name:         DefaultSuperUserName,
		Email:        DefaultSuperUserEmail,
		PasswordHash: string(hashedPassword),
		Role:         DefaultSuperUserRole,
	}

	_, err = userRepo.CreateUser(&superUser)
	if err != nil {
		log.Fatalf("Failed to create super user with default credentials: %v", err)
	}

	log.Println("Super user created successfully with default credentials!")
	log.Printf("Default Super User Details: Email=%s (Please change the password in production!)", DefaultSuperUserEmail)
}
