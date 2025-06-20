package repositories

import (
	"jaga/models"

	"gorm.io/gorm"
)

type AssetCategoryRepository interface {
	CreateAssetCategory(assetCategory *models.AssetCategory) error
	GetAssetCategories() ([]models.AssetCategory, error)
	GetAssetCategoryByID(id string) (*models.AssetCategory, error)
	UpdateAssetCategory(assetCategory *models.AssetCategory) error
	DeleteAssetCategory(id string) error
}

type assetCategoryRepository struct {
	db *gorm.DB
}

func NewAssetCategoryRepository(db *gorm.DB) AssetCategoryRepository {
	return &assetCategoryRepository{db: db}
}

func (r *assetCategoryRepository) CreateAssetCategory(assetCategory *models.AssetCategory) error {
	return r.db.Create(assetCategory).Error
}

func (r *assetCategoryRepository) GetAssetCategories() ([]models.AssetCategory, error) {
	var assetCategories []models.AssetCategory
	err := r.db.Find(&assetCategories).Error
	return assetCategories, err
}

func (r *assetCategoryRepository) GetAssetCategoryByID(id string) (*models.AssetCategory, error) {
	var assetCategory models.AssetCategory
	err := r.db.First(&assetCategory, "id = ?", id).Error
	return &assetCategory, err
}

func (r *assetCategoryRepository) UpdateAssetCategory(assetCategory *models.AssetCategory) error {
	return r.db.Save(assetCategory).Error
}

func (r *assetCategoryRepository) DeleteAssetCategory(id string) error {
	return r.db.Delete(&models.AssetCategory{}, "id = ?", id).Error
}
