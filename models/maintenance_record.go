package models

import "time"

type MaintenanceRecord struct {
	ID              string    `gorm:"primaryKey;type:char(36)"`
	AssetID         string    `gorm:"type:char(36);not null"`
	ScheduleID      *string   `gorm:"type:char(36)"`
	PerformedBy     *string   `gorm:"type:char(36)"`
	Description     string    `gorm:"type:text"`
	Status          string    `gorm:"type:enum('pending','in_progress','on_hold','finished','failed','cancelled');not null"`
	MaintenanceDate time.Time `gorm:"not null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time

	Asset Asset `gorm:"foreignKey:AssetID"`
}
