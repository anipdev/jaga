package controllers

import (
	"jaga/dto"
	"jaga/models"
	"jaga/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AssetCategoryController interface {
	CreateAssetCategory(c *gin.Context)
	GetAssetCategories(c *gin.Context)
	GetAssetCategoryByID(c *gin.Context)
	UpdateAssetCategory(c *gin.Context)
	DeleteAssetCategory(c *gin.Context)
}

type assetCategoryController struct {
	service services.AssetCategoryService
}

func NewAssetCategoryController(service services.AssetCategoryService) AssetCategoryController {
	return &assetCategoryController{service: service}
}

// CreateAssetCategory godoc
// @Summary Create a new asset category
// @Description Create a new asset category with a specified name
// @Tags Asset Categories
// @Accept json
// @Produce json
// @Param assetCategoryRequest body dto.CreateAssetCategoryRequest true "Asset category request"
// @Success 201 {object} dto.CreateAssetCategoryResponse
// @Router /v1/asset-categories [post]
func (ctrl *assetCategoryController) CreateAssetCategory(c *gin.Context) {
	var req dto.CreateAssetCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	assetCategory := models.AssetCategory{
		Name: req.Name,
	}

	if err := ctrl.service.CreateAssetCategory(&assetCategory); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create asset category: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.CreateAssetCategoryResponse{
		Message: "Asset category created successfully",
	})
}

// GetAssetCategories godoc
// @Summary Get list of asset categories
// @Description Retrieve a list of asset categories
// @Tags Asset Categories
// @Accept json
// @Produce json
// @Success 200 {object} dto.GetAssetCategoriesResponse
// @Router /v1/asset-categories [get]
func (ctrl *assetCategoryController) GetAssetCategories(c *gin.Context) {

	assetCategories, err := ctrl.service.GetAssetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve asset categories: " + err.Error()})
		return
	}

	var assetCategoryDTOs []dto.AssetCategoryDTO
	for _, ac := range assetCategories {
		assetCategoryDTOs = append(assetCategoryDTOs, dto.AssetCategoryDTO{
			ID:   ac.ID,
			Name: ac.Name,
		})
	}

	c.JSON(http.StatusOK, dto.GetAssetCategoriesResponse{
		Message:         "Asset categories retrieved successfully",
		AssetCategories: assetCategoryDTOs,
	})
}

// GetAssetCategoryByID godoc
// @Summary Get asset category by ID
// @Description Retrieve an asset category by its unique ID
// @Tags Asset Categories
// @Accept json
// @Produce json
// @Param id path string true "Asset Category ID"
// @Success 200 {object} dto.GetAssetCategoryByIDResponse
// @Router /v1/asset-categories/{id} [get]
func (ctrl *assetCategoryController) GetAssetCategoryByID(c *gin.Context) {
	id := c.Param("id")

	assetCategory, err := ctrl.service.GetAssetCategoryByID(id)
	if err != nil {
		if err.Error() == "asset category not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Asset category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve asset category: " + err.Error()})
		return
	}

	res := dto.GetAssetCategoryByIDResponse{
		Message: "Asset category retrieved successfully",
		AssetCategory: dto.AssetCategoryDTO{
			ID:   assetCategory.ID,
			Name: assetCategory.Name,
		},
	}
	c.JSON(http.StatusOK, res)
}

// UpdateAssetCategory godoc
// @Summary Update an asset category
// @Description Update an existing asset category by ID
// @Tags Asset Categories
// @Accept json
// @Produce json
// @Param id path string true "Asset Category ID"
// @Param assetCategoryRequest body dto.UpdateAssetCategoryRequest true "Update asset category request"
// @Success 200 {object} dto.UpdateAssetCategoryResponse
// @Router /v1/asset-categories/{id} [put]
func (ctrl *assetCategoryController) UpdateAssetCategory(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateAssetCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	assetCategory := models.AssetCategory{
		ID:   id,
		Name: req.Name,
	}

	if err := ctrl.service.UpdateAssetCategory(&assetCategory); err != nil {
		if err.Error() == "asset category not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update asset category: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.UpdateAssetCategoryResponse{
		Message: "Asset category updated successfully",
	})
}

// DeleteAssetCategory godoc
// @Summary Delete an asset category
// @Description Delete an asset category by its unique ID
// @Tags Asset Categories
// @Accept json
// @Produce json
// @Param id path string true "Asset Category ID"
// @Success 200 {object} dto.DeleteAssetCategoryResponse
// @Router /v1/asset-categories/{id} [delete]
func (ctrl *assetCategoryController) DeleteAssetCategory(c *gin.Context) {
	id := c.Param("id")

	if err := ctrl.service.DeleteAssetCategory(id); err != nil {
		if err.Error() == "asset category not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "cannot delete category as it is still used by assets" { // Example if service validates dependencies
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete asset category: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.DeleteAssetCategoryResponse{
		Message: "Asset category deleted successfully",
	})
}
