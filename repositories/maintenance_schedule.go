package repositories

import (
	"jaga/models"
	"time"

	"gorm.io/gorm"
)

type MaintenanceScheduleRepository interface {
	CreateMaintenanceSchedule(schedule *models.MaintenanceSchedule) error
	GetMaintenanceScheduleByID(scheduleID string) (*models.MaintenanceSchedule, error)
	GetMaintenanceSchedules(page, itemsPerPage int, sortBy, sortDir, assetID, scheduleType string, startDate, endDate *time.Time) ([]models.MaintenanceSchedule, int64, error)
	UpdateMaintenanceSchedule(schedule *models.MaintenanceSchedule) error
	DeleteMaintenanceSchedule(scheduleID string) error
}

type maintenanceScheduleRepository struct {
	db *gorm.DB
}

func NewMaintenanceScheduleRepository(db *gorm.DB) MaintenanceScheduleRepository {
	return &maintenanceScheduleRepository{db: db}
}

func (r *maintenanceScheduleRepository) CreateMaintenanceSchedule(schedule *models.MaintenanceSchedule) error {
	return r.db.Create(schedule).Error
}

func (r *maintenanceScheduleRepository) GetMaintenanceScheduleByID(scheduleID string) (*models.MaintenanceSchedule, error) {
	var schedule models.MaintenanceSchedule
	err := r.db.Preload("Asset").Where("id = ?", scheduleID).First(&schedule).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *maintenanceScheduleRepository) GetMaintenanceSchedules(page, itemsPerPage int, sortBy, sortDir, assetID, scheduleType string, startDate, endDate *time.Time) ([]models.MaintenanceSchedule, int64, error) {
	var schedules []models.MaintenanceSchedule
	var totalItems int64

	query := r.db.Model(&models.MaintenanceSchedule{}).Preload("Asset")

	if assetID != "" {
		query = query.Where("asset_id = ?", assetID)
	}
	if scheduleType != "" {
		query = query.Where("schedule_type = ?", scheduleType)
	}

	if startDate != nil && endDate != nil {
		query = query.Where("next_maintenance_date BETWEEN ? AND ?", *startDate, *endDate)
	} else if startDate != nil {
		query = query.Where("next_maintenance_date >= ?", *startDate)
	} else if endDate != nil {
		query = query.Where("next_maintenance_date <= ?", *endDate)
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
		query = query.Order("next_maintenance_date asc")
	}

	if page > 0 && itemsPerPage > 0 {
		offset := (page - 1) * itemsPerPage
		query = query.Limit(itemsPerPage).Offset(offset)
	}

	if err := query.Find(&schedules).Error; err != nil {
		return nil, 0, err
	}

	return schedules, totalItems, nil
}

func (r *maintenanceScheduleRepository) UpdateMaintenanceSchedule(schedule *models.MaintenanceSchedule) error {
	return r.db.Save(schedule).Error
}

func (r *maintenanceScheduleRepository) DeleteMaintenanceSchedule(scheduleID string) error {
	return r.db.Delete(&models.MaintenanceSchedule{}, "id = ?", scheduleID).Error
}
