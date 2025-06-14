package models

import "time"

type Asset struct {
	ID                  string `gorm:"primaryKey;type:char(36)"`
	Name                string `gorm:"type:varchar(100);not null"`
	CategoryID          string `gorm:"type:char(36);not null"`
	Location            string `gorm:"type:varchar(100)"`
	PurchaseDate        *time.Time
	LastMaintenanceDate *time.Time
	Condition           string `gorm:"type:varchar(50)"`
	Status              string `gorm:"type:enum('ready','under_maintenance','need_maintenance');not null"`
	AddedBy             string `gorm:"type:char(36)"`
	CreatedAt           time.Time
	UpdatedAt           time.Time

	Category AssetCategory `gorm:"foreignKey:CategoryID"`
}
