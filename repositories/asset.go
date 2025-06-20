package repositories

import (
	"jaga/models"

	"gorm.io/gorm"
)

type AssetRepository interface {
	CreateAsset(asset *models.Asset) (*models.Asset, error)
	GetAssetByID(assetID string) (*models.Asset, error)
	GetAssets(page, itemsPerPage int, sortBy, sortDir, search, categoryID, status string) ([]models.Asset, int64, error)
	UpdateAsset(asset *models.Asset) (*models.Asset, error)
	UpdateAssetStatus(assetID, status string) error
	DeleteAsset(assetID string) error
}

type assetRepository struct {
	db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) AssetRepository {
	return &assetRepository{db: db}
}

func (r *assetRepository) CreateAsset(asset *models.Asset) (*models.Asset, error) {
	if err := r.db.Create(asset).Error; err != nil {
		return nil, err
	}
	return asset, nil
}

func (r *assetRepository) GetAssetByID(assetID string) (*models.Asset, error) {
	var asset models.Asset
	if err := r.db.Preload("Category").Where("id = ?", assetID).First(&asset).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *assetRepository) GetAssets(page, itemsPerPage int, sortBy, sortDir, search, categoryID, status string) ([]models.Asset, int64, error) {
	var assets []models.Asset
	var totalItems int64

	query := r.db.Model(&models.Asset{}).Preload("Category")

	if search != "" {
		query = query.Where("name LIKE ? OR location LIKE ? OR `condition` LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&totalItems)

	if sortBy != "" {
		order := sortBy
		if sortDir != "" {
			order = sortBy + " " + sortDir
		}
		query = query.Order(order)
	} else {
		query = query.Order("created_at desc")
	}

	if page > 0 && itemsPerPage > 0 {
		offset := (page - 1) * itemsPerPage
		query = query.Limit(itemsPerPage).Offset(offset)
	}

	if err := query.Find(&assets).Error; err != nil {
		return nil, 0, err
	}

	return assets, totalItems, nil
}

func (r *assetRepository) UpdateAsset(asset *models.Asset) (*models.Asset, error) {
	if err := r.db.Save(asset).Error; err != nil {
		return nil, err
	}
	return asset, nil
}

func (r *assetRepository) UpdateAssetStatus(assetID, status string) error {
	return r.db.Model(&models.Asset{}).Where("id = ?", assetID).Update("status", status).Error
}

func (r *assetRepository) DeleteAsset(assetID string) error {
	return r.db.Delete(&models.Asset{}, "id = ?", assetID).Error
}
