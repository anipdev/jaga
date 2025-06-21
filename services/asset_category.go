package services

import (
	"jaga/models"
	"jaga/repositories"
	"jaga/utils"
)

type AssetCategoryService interface {
	CreateAssetCategory(assetCategory *models.AssetCategory) error
	GetAssetCategories() ([]models.AssetCategory, error)
	GetAssetCategoryByID(id string) (*models.AssetCategory, error)
	UpdateAssetCategory(assetCategory *models.AssetCategory) error
	DeleteAssetCategory(id string) error
}

type assetCategoryService struct {
	repo repositories.AssetCategoryRepository
}

func NewAssetCategoryService(repo repositories.AssetCategoryRepository) AssetCategoryService {
	return &assetCategoryService{repo: repo}
}

func (s *assetCategoryService) CreateAssetCategory(assetCategory *models.AssetCategory) error {
	if assetCategory.ID == "" {
		assetCategory.ID = utils.GenerateUUID()
	}

	return s.repo.CreateAssetCategory(assetCategory)
}

func (s *assetCategoryService) GetAssetCategories() ([]models.AssetCategory, error) {
	return s.repo.GetAssetCategories()
}

func (s *assetCategoryService) GetAssetCategoryByID(id string) (*models.AssetCategory, error) {
	return s.repo.GetAssetCategoryByID(id)
}

func (s *assetCategoryService) UpdateAssetCategory(assetCategory *models.AssetCategory) error {
	return s.repo.UpdateAssetCategory(assetCategory)
}

func (s *assetCategoryService) DeleteAssetCategory(id string) error {
	return s.repo.DeleteAssetCategory(id)
}
