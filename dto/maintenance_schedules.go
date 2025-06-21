package dto

import "time"

type MaintenanceScheduleDTO struct {
	ID                  string    `json:"id"`
	AssetID             string    `json:"asset_id"`
	AssetName           string    `json:"asset_name"`
	ScheduleType        string    `json:"schedule_type"`
	IntervalDays        *int      `json:"interval_days,omitempty"`
	NextMaintenanceDate time.Time `json:"next_maintenance_date"`
	ScheduledBy         string    `json:"scheduled_by"`
	AssignedTo          string    `json:"assigned_to"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type CreateMaintenanceScheduleRequest struct {
	AssetID             string    `json:"asset_id" binding:"required"`
	ScheduleType        string    `json:"schedule_type" binding:"required,oneof=periodic conditional"`
	IntervalDays        *int      `json:"interval_days"`
	NextMaintenanceDate time.Time `json:"next_maintenance_date" binding:"required"`
	ScheduledBy         string    `json:"scheduled_by"`
	AssignedTo          string    `json:"assigned_to"`
}

type CreateMaintenanceScheduleResponse struct {
	Message string `json:"message"`
}

type GetMaintenanceScheduleByIDResponse struct {
	Message             string                 `json:"message"`
	MaintenanceSchedule MaintenanceScheduleDTO `json:"maintenance_schedule"`
}

type GetMaintenanceSchedulesRequest struct {
	Page         int    `form:"page,default=1"`
	ItemsPerPage int    `form:"items_per_page,default=10"`
	SortBy       string `form:"sort_by,default=next_maintenance_date"`
	SortDir      string `form:"sort_dir,default=asc"`
	AssetID      string `form:"asset_id"`
	ScheduleType string `form:"schedule_type"`
}

type GetMaintenanceSchedulesResponse struct {
	Message              string                   `json:"message"`
	MaintenanceSchedules []MaintenanceScheduleDTO `json:"maintenance_schedules"`
	TotalItems           int                      `json:"total_items"`
	Page                 int                      `json:"page"`
	ItemsPerPage         int                      `json:"items_per_page"`
	TotalPages           int                      `json:"total_pages"`
}

type UpdateMaintenanceScheduleRequest struct {
	ScheduleType        string    `json:"schedule_type" binding:"omitempty,oneof=periodic conditional"`
	IntervalDays        *int      `json:"interval_days"`
	NextMaintenanceDate time.Time `json:"next_maintenance_date"`
	ScheduledBy         string    `json:"scheduled_by"`
	AssignedTo          string    `json:"assigned_to"`
}

type UpdateMaintenanceScheduleResponse struct {
	Message string `json:"message"`
}

type DeleteMaintenanceScheduleResponse struct {
	Message string `json:"message"`
}
