package models

import "time"

type AssetCategory struct {
	ID        string `gorm:"primaryKey;type:char(36)"`
	Name      string `gorm:"type:varchar(100);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
