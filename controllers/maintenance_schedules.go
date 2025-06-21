package controllers

import (
	"math"
	"net/http"

	"jaga/dto"
	"jaga/models"
	"jaga/services"

	"github.com/gin-gonic/gin"
)

type MaintenanceScheduleController interface {
	CreateMaintenanceSchedule(c *gin.Context)
	GetMaintenanceScheduleByID(c *gin.Context)
	GetMaintenanceSchedules(c *gin.Context)
	UpdateMaintenanceSchedule(c *gin.Context)
	DeleteMaintenanceSchedule(c *gin.Context)
}

type maintenanceScheduleController struct {
	service services.MaintenanceScheduleService
}

func NewMaintenanceScheduleController(service services.MaintenanceScheduleService) MaintenanceScheduleController {
	return &maintenanceScheduleController{service: service}
}

// CreateMaintenanceSchedule godoc
// @Summary Create a new maintenance schedule
// @Description Create a new maintenance schedule for a specific asset
// @Tags MaintenanceSchedules
// @Accept json
// @Produce json
// @Param request body dto.CreateMaintenanceScheduleRequest true "Create Maintenance Schedule Request"
// @Success 201 {object} dto.CreateMaintenanceScheduleResponse
// @Router /v1/maintenance-schedules [post]
func (ctrl *maintenanceScheduleController) CreateMaintenanceSchedule(c *gin.Context) {
	var req dto.CreateMaintenanceScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	schedule := &models.MaintenanceSchedule{
		AssetID:             req.AssetID,
		ScheduleType:        req.ScheduleType,
		IntervalDays:        req.IntervalDays,
		NextMaintenanceDate: req.NextMaintenanceDate,
		ScheduledBy:         req.ScheduledBy,
		AssignedTo:          req.AssignedTo,
	}

	if err := ctrl.service.CreateMaintenanceSchedule(schedule); err != nil {
		if err.Error() == "asset not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create maintenance schedule: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.CreateMaintenanceScheduleResponse{
		Message: "Maintenance schedule created successfully",
	})
}

// GetMaintenanceScheduleByID godoc
// @Summary Get a maintenance schedule by ID
// @Description Retrieve a maintenance schedule using its ID
// @Tags MaintenanceSchedules
// @Produce json
// @Param id path string true "Maintenance Schedule ID"
// @Success 200 {object} dto.GetMaintenanceScheduleByIDResponse
// @Router /v1/maintenance-schedules/{id} [get]
func (ctrl *maintenanceScheduleController) GetMaintenanceScheduleByID(c *gin.Context) {
	scheduleID := c.Param("id")

	scheduleModel, err := ctrl.service.GetMaintenanceScheduleByID(scheduleID)
	if err != nil {
		if err.Error() == "maintenance schedule not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Maintenance schedule not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve maintenance schedule: " + err.Error()})
		return
	}

	res := dto.GetMaintenanceScheduleByIDResponse{
		Message: "Maintenance schedule retrieved successfully",
		MaintenanceSchedule: dto.MaintenanceScheduleDTO{
			ID:                  scheduleModel.ID,
			AssetID:             scheduleModel.AssetID,
			AssetName:           scheduleModel.Asset.Name,
			ScheduleType:        scheduleModel.ScheduleType,
			IntervalDays:        scheduleModel.IntervalDays,
			NextMaintenanceDate: scheduleModel.NextMaintenanceDate,
			ScheduledBy:         scheduleModel.ScheduledBy,
			AssignedTo:          scheduleModel.AssignedTo,
			CreatedAt:           scheduleModel.CreatedAt,
			UpdatedAt:           scheduleModel.UpdatedAt,
		},
	}
	c.JSON(http.StatusOK, res)
}

// GetMaintenanceSchedules godoc
// @Summary Get maintenance schedules
// @Description Retrieve a paginated list of maintenance schedules
// @Tags MaintenanceSchedules
// @Produce json
// @Success 200 {object} dto.GetMaintenanceSchedulesResponse
// @Router /v1/maintenance-schedules [get]
func (ctrl *maintenanceScheduleController) GetMaintenanceSchedules(c *gin.Context) {
	var req dto.GetMaintenanceSchedulesRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters: " + err.Error()})
		return
	}

	schedulesModel, totalItems, err := ctrl.service.GetMaintenanceSchedules(
		req.Page,
		req.ItemsPerPage,
		req.SortBy,
		req.SortDir,
		req.AssetID,
		req.ScheduleType,
		req.StartDate,
		req.EndDate,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve maintenance schedules: " + err.Error()})
		return
	}

	scheduleDTOs := make([]dto.MaintenanceScheduleDTO, len(schedulesModel))
	for i, schedule := range schedulesModel {
		scheduleDTOs[i] = dto.MaintenanceScheduleDTO{
			ID:                  schedule.ID,
			AssetID:             schedule.AssetID,
			AssetName:           schedule.Asset.Name,
			ScheduleType:        schedule.ScheduleType,
			IntervalDays:        schedule.IntervalDays,
			NextMaintenanceDate: schedule.NextMaintenanceDate,
			ScheduledBy:         schedule.ScheduledBy,
			AssignedTo:          schedule.AssignedTo,
			CreatedAt:           schedule.CreatedAt,
			UpdatedAt:           schedule.UpdatedAt,
		}
	}

	totalPages := 0
	if req.ItemsPerPage > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(req.ItemsPerPage)))
	} else if totalItems > 0 {
		totalPages = 1
	}

	c.JSON(http.StatusOK, dto.GetMaintenanceSchedulesResponse{
		Message:              "Maintenance schedules retrieved successfully",
		MaintenanceSchedules: scheduleDTOs,
		TotalItems:           int(totalItems),
		Page:                 req.Page,
		ItemsPerPage:         req.ItemsPerPage,
		TotalPages:           totalPages,
	})
}

// UpdateMaintenanceSchedule godoc
// @Summary Update a maintenance schedule
// @Description Update an existing maintenance schedule by its ID
// @Tags MaintenanceSchedules
// @Produce json
// @Param id path string true "Maintenance Schedule ID"
// @Param request body dto.UpdateMaintenanceScheduleRequest true "Update Maintenance Schedule Request"
// @Success 200 {object} dto.UpdateMaintenanceScheduleResponse
// @Router /v1/maintenance-schedules/{id} [put]
func (ctrl *maintenanceScheduleController) UpdateMaintenanceSchedule(c *gin.Context) {
	scheduleID := c.Param("id")
	var req dto.UpdateMaintenanceScheduleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	existingSchedule, err := ctrl.service.GetMaintenanceScheduleByID(scheduleID)
	if err != nil {
		if err.Error() == "maintenance schedule not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Maintenance schedule not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve existing schedule: " + err.Error()})
		return
	}

	updatedSchedule := &models.MaintenanceSchedule{
		ID:                  scheduleID,
		AssetID:             existingSchedule.AssetID,
		ScheduleType:        req.ScheduleType,
		IntervalDays:        req.IntervalDays,
		NextMaintenanceDate: req.NextMaintenanceDate,
		ScheduledBy:         req.ScheduledBy,
		AssignedTo:          req.AssignedTo,
	}

	if err := ctrl.service.UpdateMaintenanceSchedule(updatedSchedule); err != nil {
		if err.Error() == "maintenance schedule not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update maintenance schedule: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.UpdateMaintenanceScheduleResponse{
		Message: "Maintenance schedule updated successfully!",
	})
}

// DeleteMaintenanceSchedule godoc
// @Summary Delete a maintenance schedule
// @Description Delete an existing maintenance schedule by its ID
// @Tags MaintenanceSchedules
// @Produce json
// @Param id path string true "Maintenance Schedule ID"
// @Success 200 {object} dto.DeleteMaintenanceScheduleResponse
// @Router /v1/maintenance-schedules/{id} [delete]
func (ctrl *maintenanceScheduleController) DeleteMaintenanceSchedule(c *gin.Context) {
	scheduleID := c.Param("id")

	if err := ctrl.service.DeleteMaintenanceSchedule(scheduleID); err != nil {
		if err.Error() == "maintenance schedule not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete maintenance schedule: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.DeleteMaintenanceScheduleResponse{
		Message: "Maintenance schedule deleted successfully!",
	})
}
