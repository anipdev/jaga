package dto

import "time"

type AssetDTO struct {
	ID                  string     `json:"id"`
	Name                string     `json:"name"`
	CategoryID          string     `json:"category_id"`
	CategoryName        string     `json:"category_name"`
	Location            string     `json:"location"`
	PurchaseDate        *time.Time `json:"purchase_date"`
	LastMaintenanceDate *time.Time `json:"last_maintenance_date"`
	Condition           string     `json:"condition"`
	Status              string     `json:"status"`
	AddedBy             string     `json:"added_by"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type CreateAssetRequest struct {
	Name                string     `json:"name" binding:"required,min=2,max=100"`
	CategoryID          string     `json:"category_id" binding:"required"`
	Location            string     `json:"location" binding:"omitempty,max=100"`
	PurchaseDate        *time.Time `json:"purchase_date" binding:"omitempty"`
	LastMaintenanceDate *time.Time `json:"last_maintenance_date" binding:"omitempty"`
	Condition           string     `json:"condition" binding:"omitempty,max=50"`
	Status              string     `json:"status" binding:"required,oneof=ready under_maintenance need_maintenance"`
	AddedBy             string     `json:"added_by" binding:"required"`
}

type CreateAssetResponse struct {
	Message string `json:"message"`
}

type GetAssetByIDResponse struct {
	Message string   `json:"message"`
	Asset   AssetDTO `json:"asset"`
}

type GetAssetsRequest struct {
	Page         int    `json:"page" binding:"omitempty,min=1"`
	ItemsPerPage int    `json:"items_per_page" binding:"omitempty,min=1"`
	SortBy       string `json:"sort_by" binding:"omitempty,oneof=id name category_id location purchase_date last_maintenance_date condition status added_by created_at updated_at"`
	SortDir      string `json:"sort_dir" binding:"omitempty,oneof=asc desc"`
	Search       string `json:"search" binding:"omitempty"`
	CategoryID   string `json:"category_id" binding:"omitempty"`
	Status       string `json:"status" binding:"omitempty,oneof=ready under_maintenance need_maintenance"`
}

type GetAssetsResponse struct {
	Message      string     `json:"message"`
	Assets       []AssetDTO `json:"assets"`
	TotalItems   int        `json:"total_items"`
	Page         int        `json:"page"`
	ItemsPerPage int        `json:"items_per_page"`
	TotalPages   int        `json:"total_pages"`
}

type UpdateAssetRequest struct {
	ID                  string     `json:"id" binding:"required"`
	Name                string     `json:"name" binding:"omitempty,min=2,max=100"`
	CategoryID          string     `json:"category_id" binding:"omitempty"`
	Location            string     `json:"location" binding:"omitempty,max=100"`
	PurchaseDate        *time.Time `json:"purchase_date" binding:"omitempty"`
	LastMaintenanceDate *time.Time `json:"last_maintenance_date" binding:"omitempty"`
	Condition           string     `json:"condition" binding:"omitempty,max=50"`
	Status              string     `json:"status" binding:"omitempty,oneof=ready under_maintenance need_maintenance"`
	AddedBy             string     `json:"added_by" binding:"omitempty"`
}

type UpdateAssetStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=ready under_maintenance need_maintenance"`
}

type UpdateAssetResponse struct {
	Message string `json:"message"`
}

type DeleteAssetResponse struct {
	Message string `json:"message"`
}
