package models

import "time"

type User struct {
	ID           string `gorm:"primaryKey;type:char(36)"`
	Name         string `gorm:"type:varchar(100);not null"`
	Email        string `gorm:"type:varchar(100);unique;not null"`
	PasswordHash string `gorm:"type:text;not null"`
	Role         string `gorm:"type:enum('super_user','admin','technician','manager');not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
