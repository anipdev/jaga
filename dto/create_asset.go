package dto

import "time"

type CreateAsset struct {
	Name         string     `json:"name" validate:"required"`
	CategoryID   string     `json:"category_id" validate:"required,uuid4"`
	Location     string     `json:"location"`
	PurchaseDate *time.Time `json:"purchase_date"`
	Condition    string     `json:"condition"`
	Status       string     `json:"status" validate:"required,oneof=ready under_maintenance need_maintenance"`
}
