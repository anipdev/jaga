package models

import "time"

type MaintenanceSchedule struct {
	ID                  string `gorm:"primaryKey;type:char(36)"`
	AssetID             string `gorm:"type:char(36);not null"`
	ScheduleType        string `gorm:"type:enum('periodic','conditional');not null"`
	IntervalDays        *int
	NextMaintenanceDate time.Time `gorm:"not null"`
	ScheduledBy         string    `gorm:"type:char(36)"`
	AssignedTo          string    `gorm:"type:char(36)"`
	CreatedAt           time.Time
	UpdatedAt           time.Time

	Asset Asset `gorm:"foreignKey:AssetID"`
}
