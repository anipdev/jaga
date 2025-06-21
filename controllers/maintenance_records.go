package controllers

import (
	"math"
	"net/http"

	"jaga/dto"
	"jaga/models"
	"jaga/services"

	"github.com/gin-gonic/gin"
)

type MaintenanceRecordController interface {
	CreateMaintenanceRecord(c *gin.Context)
	GetMaintenanceRecordByID(c *gin.Context)
	GetMaintenanceRecords(c *gin.Context)
	UpdateMaintenanceRecord(c *gin.Context)
	UpdateMaintenanceRecordStatus(c *gin.Context)
	DeleteMaintenanceRecord(c *gin.Context)
}

type maintenanceRecordController struct {
	service services.MaintenanceRecordService
}

func NewMaintenanceRecordController(service services.MaintenanceRecordService) MaintenanceRecordController {
	return &maintenanceRecordController{service: service}
}

func (ctrl *maintenanceRecordController) CreateMaintenanceRecord(c *gin.Context) {
	var req dto.CreateMaintenanceRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	newRecord := &models.MaintenanceRecord{
		AssetID:         req.AssetID,
		ScheduleID:      req.ScheduleID,
		PerformedBy:     req.PerformedBy,
		Description:     req.Description,
		Status:          req.Status,
		MaintenanceDate: req.MaintenanceDate,
	}

	if err := ctrl.service.CreateMaintenanceRecord(newRecord); err != nil {
		if err.Error() == "asset not found" ||
			err.Error() == "maintenance schedule not found" ||
			err.Error() == "performed by not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create maintenance record: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.CreateMaintenanceRecordResponse{
		Message: "Maintenance record created successfully",
	})
}

func (ctrl *maintenanceRecordController) GetMaintenanceRecordByID(c *gin.Context) {
	recordID := c.Param("id")
	recordModel, err := ctrl.service.GetMaintenanceRecordByID(recordID)
	if err != nil {
		if err.Error() == "maintenance record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Maintenance record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve maintenance record: " + err.Error()})
		return
	}

	res := dto.GetMaintenanceRecordByIDResponse{
		Message: "Maintenance record retrieved successfully",
		MaintenanceRecord: dto.MaintenanceRecordDTO{
			ID:              recordModel.ID,
			AssetID:         recordModel.AssetID,
			AssetName:       recordModel.Asset.Name,
			ScheduleID:      recordModel.ScheduleID,
			PerformedBy:     recordModel.PerformedBy,
			Description:     recordModel.Description,
			Status:          recordModel.Status,
			MaintenanceDate: recordModel.MaintenanceDate,
			CreatedAt:       recordModel.CreatedAt,
			UpdatedAt:       recordModel.UpdatedAt,
		},
	}
	c.JSON(http.StatusOK, res)
}

func (ctrl *maintenanceRecordController) GetMaintenanceRecords(c *gin.Context) {
	var req dto.GetMaintenanceRecordsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters: " + err.Error()})
		return
	}

	recordsModel, totalItems, err := ctrl.service.GetMaintenanceRecords(
		req.Page,
		req.ItemsPerPage,
		req.SortBy,
		req.SortDir,
		req.AssetID,
		req.ScheduleID,
		req.Status,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve maintenance records: " + err.Error()})
		return
	}

	recordDTOs := make([]dto.MaintenanceRecordDTO, len(recordsModel))
	for i, record := range recordsModel {
		recordDTOs[i] = dto.MaintenanceRecordDTO{
			ID:              record.ID,
			AssetID:         record.AssetID,
			AssetName:       record.Asset.Name,
			ScheduleID:      record.ScheduleID,
			PerformedBy:     record.PerformedBy,
			Description:     record.Description,
			Status:          record.Status,
			MaintenanceDate: record.MaintenanceDate,
			CreatedAt:       record.CreatedAt,
			UpdatedAt:       record.UpdatedAt,
		}
	}

	totalPages := 0
	if req.ItemsPerPage > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(req.ItemsPerPage)))
	}

	c.JSON(http.StatusOK, dto.GetMaintenanceRecordsResponse{
		Message:            "Maintenance records retrieved successfully",
		MaintenanceRecords: recordDTOs,
		TotalItems:         int(totalItems),
		Page:               req.Page,
		ItemsPerPage:       req.ItemsPerPage,
		TotalPages:         totalPages,
	})
}

func (ctrl *maintenanceRecordController) UpdateMaintenanceRecord(c *gin.Context) {
	recordID := c.Param("id")
	var req dto.UpdateMaintenanceRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	updatedRecord := &models.MaintenanceRecord{
		ID:              recordID,
		AssetID:         req.AssetID,
		ScheduleID:      req.ScheduleID,
		PerformedBy:     req.PerformedBy,
		Description:     req.Description,
		Status:          req.Status,
		MaintenanceDate: req.MaintenanceDate,
	}

	if err := ctrl.service.UpdateMaintenanceRecord(updatedRecord); err != nil {
		if err.Error() == "maintenance record not found" ||
			err.Error() == "asset not found" ||
			err.Error() == "maintenance schedule not found" ||
			err.Error() == "performed by not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update maintenance record: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.UpdateMaintenanceRecordResponse{
		Message: "Maintenance record updated successfully!",
	})
}

func (ctrl *maintenanceRecordController) UpdateMaintenanceRecordStatus(c *gin.Context) {
	recordID := c.Param("id")
	var req dto.UpdateMaintenanceRecordStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if err := ctrl.service.UpdateMaintenanceRecordStatus(recordID, req.Status); err != nil {
		if err.Error() == "maintenance record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update maintenance record status: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.UpdateMaintenanceRecordResponse{
		Message: "Asset updated successfully!",
	})
}

func (ctrl *maintenanceRecordController) DeleteMaintenanceRecord(c *gin.Context) {
	recordID := c.Param("id")
	if err := ctrl.service.DeleteMaintenanceRecord(recordID); err != nil {
		if err.Error() == "maintenance record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete maintenance record: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.DeleteMaintenanceRecordResponse{
		Message: "Maintenance record deleted successfully!",
	})
}
