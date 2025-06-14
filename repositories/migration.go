package repositories

import (
	"jaga/config"
	"jaga/models"
)

func AutoMigrate() {
	config.DB.AutoMigrate(
		&models.User{},
		&models.AssetCategory{},
		&models.Asset{},
		&models.MaintenanceSchedule{},
		&models.MaintenanceRecord{},
	)
}
