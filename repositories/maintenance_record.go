package repositories

import (
	"jaga/models"

	"gorm.io/gorm"
)

type MaintenanceRecordRepository interface {
	CreateMaintenanceRecord(record *models.MaintenanceRecord) error
	GetMaintenanceRecordByID(recordID string) (*models.MaintenanceRecord, error)
	GetMaintenanceRecords(
		page, itemsPerPage int,
		sortBy, sortDir, assetID, status string,
		scheduleIDs ...string) ([]models.MaintenanceRecord, int64, error)
	UpdateMaintenanceRecord(record *models.MaintenanceRecord) error
	DeleteMaintenanceRecord(recordID string) error
}

type maintenanceRecordRepository struct {
	db *gorm.DB
}

func NewMaintenanceRecordRepository(db *gorm.DB) MaintenanceRecordRepository {
	return &maintenanceRecordRepository{db: db}
}

func (r *maintenanceRecordRepository) CreateMaintenanceRecord(record *models.MaintenanceRecord) error {
	return r.db.Create(record).Error
}

func (r *maintenanceRecordRepository) GetMaintenanceRecordByID(recordID string) (*models.MaintenanceRecord, error) {
	var record models.MaintenanceRecord
	err := r.db.Preload("Asset").Where("id = ?", recordID).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *maintenanceRecordRepository) GetMaintenanceRecords(
	page, itemsPerPage int,
	sortBy, sortDir, assetID, status string,
	scheduleIDs ...string,
) ([]models.MaintenanceRecord, int64, error) {

	var records []models.MaintenanceRecord
	var totalItems int64

	query := r.db.Model(&models.MaintenanceRecord{}).Preload("Asset")

	if assetID != "" {
		query = query.Where("asset_id = ?", assetID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if len(scheduleIDs) > 0 {
		query = query.Where("schedule_id IN (?)", scheduleIDs)
	}

	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	if sortBy != "" {
		order := sortBy
		if sortDir != "" {
			order = sortBy + " " + sortDir
		}
		query = query.Order(order)
	} else {
		query = query.Order("maintenance_date desc")
	}

	if page > 0 && itemsPerPage > 0 {
		offset := (page - 1) * itemsPerPage
		query = query.Limit(itemsPerPage).Offset(offset)
	}

	if err := query.Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, totalItems, nil
}

func (r *maintenanceRecordRepository) UpdateMaintenanceRecord(record *models.MaintenanceRecord) error {
	return r.db.Save(record).Error
}

func (r *maintenanceRecordRepository) DeleteMaintenanceRecord(recordID string) error {
	return r.db.Delete(&models.MaintenanceRecord{}, "id = ?", recordID).Error
}
