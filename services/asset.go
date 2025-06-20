package services

import (
	"errors"

	"jaga/models"
	"jaga/repositories"

	"gorm.io/gorm"
)

type AssetService interface {
	CreateAsset(asset *models.Asset) (*models.Asset, error)
	GetAssetByID(assetID string) (*models.Asset, error)
	GetAssets(page, itemsPerPage int, sortBy, sortDir, search, categoryID, status string) ([]models.Asset, int64, error)
	UpdateAsset(asset *models.Asset) (*models.Asset, error)
	UpdateAssetStatus(assetID, status string) error
	DeleteAsset(assetID string) error
}

type assetService struct {
	assetRepo    repositories.AssetRepository
	categoryRepo repositories.AssetCategoryRepository
}

func NewAssetService(assetRepo repositories.AssetRepository, categoryRepo repositories.AssetCategoryRepository) AssetService {
	return &assetService{
		assetRepo:    assetRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *assetService) CreateAsset(asset *models.Asset) (*models.Asset, error) {

	_, err := s.categoryRepo.GetAssetCategoryByID(asset.CategoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	createdAsset, err := s.assetRepo.CreateAsset(asset)
	if err != nil {
		return nil, err
	}

	return s.assetRepo.GetAssetByID(createdAsset.ID)
}

func (s *assetService) GetAssetByID(assetID string) (*models.Asset, error) {
	asset, err := s.assetRepo.GetAssetByID(assetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("asset not found")
		}
		return nil, err
	}
	return asset, nil
}

func (s *assetService) GetAssets(page, itemsPerPage int, sortBy, sortDir, search, categoryID, status string) ([]models.Asset, int64, error) {
	assets, totalItems, err := s.assetRepo.GetAssets(page, itemsPerPage, sortBy, sortDir, search, categoryID, status)
	if err != nil {
		return nil, 0, err
	}
	return assets, totalItems, nil
}

func (s *assetService) UpdateAsset(asset *models.Asset) (*models.Asset, error) {

	_, err := s.assetRepo.GetAssetByID(asset.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("asset not found")
		}
		return nil, err
	}

	if asset.CategoryID != "" {
		_, err := s.categoryRepo.GetAssetCategoryByID(asset.CategoryID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("category not found")
			}
			return nil, err
		}
	}

	updatedAsset, err := s.assetRepo.UpdateAsset(asset)
	if err != nil {
		return nil, err
	}

	return s.assetRepo.GetAssetByID(updatedAsset.ID)
}

func (s *assetService) UpdateAssetStatus(assetID, status string) error {

	_, err := s.assetRepo.GetAssetByID(assetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("asset not found")
		}
		return err
	}

	return s.assetRepo.UpdateAssetStatus(assetID, status)
}

func (s *assetService) DeleteAsset(assetID string) error {

	_, err := s.assetRepo.GetAssetByID(assetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("asset not found")
		}
		return err
	}

	return s.assetRepo.DeleteAsset(assetID)
}
