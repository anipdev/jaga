package dto

import "time"

type MaintenanceRecordDTO struct {
	ID              string    `json:"id"`
	AssetID         string    `json:"asset_id"`
	AssetName       string    `json:"asset_name"`
	ScheduleID      *string   `json:"schedule_id,omitempty"`
	PerformedBy     *string   `json:"performed_by_user_id,omitempty"`
	PerformedByName *string   `json:"performed_by_user_name,omitempty"`
	Description     string    `json:"description"`
	Status          string    `json:"status"`
	MaintenanceDate time.Time `json:"maintenance_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateMaintenanceRecordRequest struct {
	AssetID         string    `json:"asset_id" binding:"required"`
	ScheduleID      *string   `json:"schedule_id,omitempty"`
	PerformedBy     *string   `json:"performed_by,omitempty"`
	Description     string    `json:"description" binding:"required,min=5,max=500"`
	Status          string    `json:"status" binding:"required,oneof=pending in_progress on_hold finished failed cancelled"`
	MaintenanceDate time.Time `json:"maintenance_date" binding:"required"`
}

type CreateMaintenanceRecordResponse struct {
	Message string `json:"message"`
}

type GetMaintenanceRecordByIDResponse struct {
	Message           string               `json:"message"`
	MaintenanceRecord MaintenanceRecordDTO `json:"maintenance_record"`
}

type GetMaintenanceRecordsRequest struct {
	Page         int    `form:"page,default=1"`
	ItemsPerPage int    `form:"items_per_page,default=10"`
	SortBy       string `form:"sort_by,default=maintenance_date"`
	SortDir      string `form:"sort_dir,default=desc"`
	AssetID      string `form:"asset_id,omitempty"`
	ScheduleID   string `form:"schedule_id,omitempty"`
	Status       string `form:"status,omitempty"`
}

type GetMaintenanceRecordsResponse struct {
	Message            string                 `json:"message"`
	MaintenanceRecords []MaintenanceRecordDTO `json:"maintenance_records"`
	TotalItems         int                    `json:"total_items"`
	Page               int                    `json:"page"`
	ItemsPerPage       int                    `json:"items_per_page"`
	TotalPages         int                    `json:"total_pages"`
}

type UpdateMaintenanceRecordRequest struct {
	AssetID         string    `json:"asset_id,omitempty"`
	ScheduleID      *string   `json:"schedule_id,omitempty"`
	PerformedBy     *string   `json:"performed_by,omitempty"`
	Description     string    `json:"description,omitempty" binding:"omitempty,min=5,max=500"`
	Status          string    `json:"status,omitempty" binding:"omitempty,oneof=pending in_progress on_hold finished failed cancelled"`
	MaintenanceDate time.Time `json:"maintenance_date,omitempty"`
}

type UpdateMaintenanceRecordStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending in_progress on_hold finished failed cancelled"`
}

type UpdateMaintenanceRecordResponse struct {
	Message string `json:"message"`
}

type DeleteMaintenanceRecordResponse struct {
	Message string `json:"message"`
}
