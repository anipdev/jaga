package controllers

import (
	"math"
	"net/http"

	"jaga/dto"
	"jaga/models"
	"jaga/services"

	"github.com/gin-gonic/gin"
)

type AssetController struct {
	AssetService services.AssetService
}

func NewAssetController(assetService services.AssetService) *AssetController {
	return &AssetController{AssetService: assetService}
}

// CreateAsset godoc
// @Summary Create a new asset
// @Description Create a new asset with the provided details
// @Tags Assets
// @Accept json
// @Produce json
// @Param assetRequest body dto.CreateAssetRequest true "Create asset request"
// @Success 201 {object} dto.CreateAssetResponse
// @Router /v1/assets [post]
func (ctrl *AssetController) CreateAsset(c *gin.Context) {
	var req dto.CreateAssetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	addedByID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}
	req.AddedBy = addedByID.(string)

	newAsset := &models.Asset{
		Name:                req.Name,
		CategoryID:          req.CategoryID,
		Location:            req.Location,
		PurchaseDate:        req.PurchaseDate,
		LastMaintenanceDate: req.LastMaintenanceDate,
		Condition:           req.Condition,
		Status:              req.Status,
		AddedBy:             req.AddedBy,
	}

	err := ctrl.AssetService.CreateAsset(newAsset)
	if err != nil {
		if err.Error() == "category not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create asset: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.CreateAssetResponse{
		Message: "Asset created successfully!",
	})
}

// GetAssetByID godoc
// @Summary Get asset by ID
// @Description Retrieve an asset by its unique ID
// @Tags Assets
// @Accept json
// @Produce json
// @Param id path string true "Asset ID"
// @Success 200 {object} dto.GetAssetByIDResponse
// @Router /v1/assets/{id} [get]
func (ctrl *AssetController) GetAssetByID(c *gin.Context) {
	assetID := c.Param("id")

	assetModel, err := ctrl.AssetService.GetAssetByID(assetID)
	if err != nil {
		if err.Error() == "asset not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Asset not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve asset: " + err.Error(),
		})
		return
	}

	assetDTO := dto.AssetDTO{
		ID:                  assetModel.ID,
		Name:                assetModel.Name,
		CategoryID:          assetModel.CategoryID,
		CategoryName:        assetModel.Category.Name,
		Location:            assetModel.Location,
		PurchaseDate:        assetModel.PurchaseDate,
		LastMaintenanceDate: assetModel.LastMaintenanceDate,
		Condition:           assetModel.Condition,
		Status:              assetModel.Status,
		AddedBy:             assetModel.AddedBy,
		CreatedAt:           assetModel.CreatedAt,
		UpdatedAt:           assetModel.UpdatedAt,
	}

	c.JSON(http.StatusOK, dto.GetAssetByIDResponse{
		Message: "Asset retrieved successfully",
		Asset:   assetDTO,
	})
}

// GetAssets godoc
// @Summary Get list of assets
// @Description Retrieve a list of assets with pagination and optional filters
// @Tags Assets
// @Accept json
// @Produce json
// @Success 200 {object} dto.GetAssetsResponse
// @Router /v1/assets [get]
func (ctrl *AssetController) GetAssets(c *gin.Context) {
	var req dto.GetAssetsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters: " + err.Error()})
		return
	}

	assetsModel, totalItems, err := ctrl.AssetService.GetAssets(
		req.Page,
		req.ItemsPerPage,
		req.SortBy,
		req.SortDir,
		req.Search,
		req.CategoryID,
		req.Status,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve assets: " + err.Error()})
		return
	}

	assetDTOs := make([]dto.AssetDTO, len(assetsModel))
	for i, asset := range assetsModel {
		assetDTOs[i] = dto.AssetDTO{
			ID:                  asset.ID,
			Name:                asset.Name,
			CategoryID:          asset.CategoryID,
			CategoryName:        asset.Category.Name,
			Location:            asset.Location,
			PurchaseDate:        asset.PurchaseDate,
			LastMaintenanceDate: asset.LastMaintenanceDate,
			Condition:           asset.Condition,
			Status:              asset.Status,
			AddedBy:             asset.AddedBy,
			CreatedAt:           asset.CreatedAt,
			UpdatedAt:           asset.UpdatedAt,
		}
	}

	totalPages := 0
	if req.ItemsPerPage > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(req.ItemsPerPage)))
	}

	c.JSON(http.StatusOK, dto.GetAssetsResponse{
		Message:      "Assets retrieved successfully",
		Assets:       assetDTOs,
		TotalItems:   int(totalItems),
		Page:         req.Page,
		ItemsPerPage: req.ItemsPerPage,
		TotalPages:   totalPages,
	})
}

// UpdateAsset godoc
// @Summary Update an existing asset
// @Description Update the details of an existing asset by its ID
// @Tags Assets
// @Accept json
// @Produce json
// @Param id path string true "Asset ID"
// @Param assetRequest body dto.UpdateAssetRequest true "Update asset request"
// @Success 200 {object} dto.UpdateAssetResponse
// @Router /v1/assets/{id} [put]
func (ctrl *AssetController) UpdateAsset(c *gin.Context) {
	assetID := c.Param("id")
	var req dto.UpdateAssetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	updatedAsset := &models.Asset{
		ID:                  assetID,
		Name:                req.Name,
		CategoryID:          req.CategoryID,
		Location:            req.Location,
		PurchaseDate:        req.PurchaseDate,
		LastMaintenanceDate: req.LastMaintenanceDate,
		Condition:           req.Condition,
		Status:              req.Status,
		AddedBy:             req.AddedBy, // This field usually doesn't change on update. Consider its purpose.
	}

	err := ctrl.AssetService.UpdateAsset(updatedAsset)
	if err != nil {
		if err.Error() == "asset not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "category not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update asset: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.UpdateAssetResponse{
		Message: "Asset updated successfully!",
	})
}

// UpdateAssetStatus godoc
// @Summary Update the status of an asset
// @Description Update the status of an asset by its ID
// @Tags Assets
// @Accept json
// @Produce json
// @Param id path string true "Asset ID"
// @Param statusRequest body dto.UpdateAssetStatusRequest true "Update asset status request"
// @Success 200 {object} dto.UpdateAssetResponse
// @Router /v1/assets/{id}/status [put]
func (ctrl *AssetController) UpdateAssetStatus(c *gin.Context) {
	assetID := c.Param("id")
	var req dto.UpdateAssetStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	err := ctrl.AssetService.UpdateAssetStatus(assetID, req.Status)
	if err != nil {
		if err.Error() == "asset not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update asset status: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.UpdateAssetResponse{
		Message: "Asset updated successfully!",
	})
}

// DeleteAsset godoc
// @Summary Delete an asset
// @Description Delete an asset by its unique ID
// @Tags Assets
// @Accept json
// @Produce json
// @Param id path string true "Asset ID"
// @Success 200 {object} dto.DeleteAssetResponse
// @Router /v1/assets/{id} [delete]
func (ctrl *AssetController) DeleteAsset(c *gin.Context) {
	assetID := c.Param("id")

	err := ctrl.AssetService.DeleteAsset(assetID)
	if err != nil {
		if err.Error() == "asset not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete asset: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.DeleteAssetResponse{
		Message: "Asset deleted successfully!",
	})
}
