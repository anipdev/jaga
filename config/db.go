package config

import (
	"fmt"
	"jaga/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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
