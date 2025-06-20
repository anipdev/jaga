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

	_, err := ctrl.AssetService.CreateAsset(newAsset)
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

func (ctrl *AssetController) GetAssetByID(c *gin.Context) {
	assetID := c.Param("id")

	assetModel, err := ctrl.AssetService.GetAssetByID(assetID)
	if err != nil {
		if err.Error() == "asset not found" {
			c.JSON(http.StatusNotFound, dto.GetAssetByIDResponse{
				Message: "Asset not found",
				Asset:   dto.AssetDTO{},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.GetAssetByIDResponse{
			Message: "Failed to retrieve asset: " + err.Error(),
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

func (ctrl *AssetController) UpdateAsset(c *gin.Context) {
	assetID := c.Param("id")
	var req dto.UpdateAssetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	req.ID = assetID

	existingAsset, err := ctrl.AssetService.GetAssetByID(assetID)
	if err != nil {
		if err.Error() == "asset not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get existing asset: " + err.Error()})
		return
	}

	if req.Name != "" {
		existingAsset.Name = req.Name
	}
	if req.CategoryID != "" {
		existingAsset.CategoryID = req.CategoryID
	}
	if req.Location != "" {
		existingAsset.Location = req.Location
	}
	if req.PurchaseDate != nil {
		existingAsset.PurchaseDate = req.PurchaseDate
	}
	if req.LastMaintenanceDate != nil {
		existingAsset.LastMaintenanceDate = req.LastMaintenanceDate
	}
	if req.Condition != "" {
		existingAsset.Condition = req.Condition
	}
	if req.Status != "" {
		existingAsset.Status = req.Status
	}
	if req.AddedBy != "" {
		existingAsset.AddedBy = req.AddedBy
	}

	_, err = ctrl.AssetService.UpdateAsset(existingAsset)
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
