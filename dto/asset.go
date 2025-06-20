package dto

import "time"

type CreateAssetRequest struct {
	Name         string     `json:"name" validate:"required"`
	CategoryID   string     `json:"category_id" validate:"required,uuid4"`
	Location     string     `json:"location"`
	PurchaseDate *time.Time `json:"purchase_date"`
	Condition    string     `json:"condition"`
	Status       string     `json:"status" validate:"required,oneof=ready under_maintenance need_maintenance"`
}

type CreateAssetResponse struct {
	Message string `json:"message"`
}

type GetAssetResponse struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	CategoryID   string     `json:"category_id"`
	CategoryName string     `json:"category_name"`
	Location     string     `json:"location"`
	PurchaseDate *time.Time `json:"purchase_date"`
	Condition    string     `json:"condition"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
