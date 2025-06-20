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

func (ctrl *assetCategoryController) CreateAssetCategory(c *gin.Context) {
	var req dto.CreateAssetCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assetCategory := models.AssetCategory{
		Name: req.Name,
	}

	if err := ctrl.service.CreateAssetCategory(&assetCategory); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := dto.CreateAssetCategoryResponse{
		Message: "Asset category created successfully",
	}
	c.JSON(http.StatusCreated, res)
}

func (ctrl *assetCategoryController) GetAssetCategories(c *gin.Context) {

	assetCategories, err := ctrl.service.GetAssetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var assetCategoryDTOs []dto.AssetCategoryDTO
	for _, ac := range assetCategories {
		assetCategoryDTOs = append(assetCategoryDTOs, dto.AssetCategoryDTO{
			ID:   ac.ID,
			Name: ac.Name,
		})
	}

	res := dto.GetAssetCategoriesResponse{
		Message: "Asset categories retrieved successfully",
	}
	c.JSON(http.StatusOK, res)
}

func (ctrl *assetCategoryController) GetAssetCategoryByID(c *gin.Context) {
	id := c.Param("id")

	assetCategory, err := ctrl.service.GetAssetCategoryByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset category not found"})
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

func (ctrl *assetCategoryController) UpdateAssetCategory(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateAssetCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = id

	assetCategory := models.AssetCategory{
		ID:   req.ID,
		Name: req.Name,
	}

	if err := ctrl.service.UpdateAssetCategory(&assetCategory); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := dto.UpdateAssetCategoryResponse{
		Message: "Asset category updated successfully",
	}
	c.JSON(http.StatusOK, res)
}

func (ctrl *assetCategoryController) DeleteAssetCategory(c *gin.Context) {
	id := c.Param("id")

	if err := ctrl.service.DeleteAssetCategory(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := dto.DeleteAssetCategoryResponse{
		Message: "Asset category deleted successfully",
	}
	c.JSON(http.StatusOK, res)
}
