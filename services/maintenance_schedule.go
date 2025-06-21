package services

import (
	"errors"
	"jaga/models"
	"jaga/repositories"
	"jaga/utils"
	"time"

	"gorm.io/gorm"
)

type MaintenanceScheduleService interface {
	CreateMaintenanceSchedule(schedule *models.MaintenanceSchedule) error
	GetMaintenanceScheduleByID(scheduleID string) (*models.MaintenanceSchedule, error)
	GetMaintenanceSchedules(page, itemsPerPage int, sortBy, sortDir, assetID, scheduleType string, startDate, endDate *time.Time) ([]models.MaintenanceSchedule, int64, error)
	UpdateMaintenanceSchedule(schedule *models.MaintenanceSchedule) error
	DeleteMaintenanceSchedule(scheduleID string) error
}

type maintenanceScheduleService struct {
	repo      repositories.MaintenanceScheduleRepository
	assetRepo repositories.AssetRepository
}

func NewMaintenanceScheduleService(repo repositories.MaintenanceScheduleRepository, assetRepo repositories.AssetRepository) MaintenanceScheduleService {
	return &maintenanceScheduleService{repo: repo, assetRepo: assetRepo}
}

func (s *maintenanceScheduleService) CreateMaintenanceSchedule(schedule *models.MaintenanceSchedule) error {
	_, err := s.assetRepo.GetAssetByID(schedule.AssetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("asset not found")
		}
		return err
	}

	if schedule.ID == "" {
		schedule.ID = utils.GenerateUUID()
	}
	return s.repo.CreateMaintenanceSchedule(schedule)
}

func (s *maintenanceScheduleService) GetMaintenanceScheduleByID(scheduleID string) (*models.MaintenanceSchedule, error) {
	schedule, err := s.repo.GetMaintenanceScheduleByID(scheduleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("maintenance schedule not found")
		}
		return nil, err
	}
	return schedule, nil
}

func (s *maintenanceScheduleService) GetMaintenanceSchedules(page, itemsPerPage int, sortBy, sortDir, assetID, scheduleType string, startDate, endDate *time.Time) ([]models.MaintenanceSchedule, int64, error) {
	schedules, totalItems, err := s.repo.GetMaintenanceSchedules(
		page, itemsPerPage, sortBy, sortDir, assetID, scheduleType, nil, nil,
	)
	if err != nil {
		return nil, 0, err
	}
	return schedules, totalItems, nil
}

func (s *maintenanceScheduleService) UpdateMaintenanceSchedule(schedule *models.MaintenanceSchedule) error {

	_, err := s.repo.GetMaintenanceScheduleByID(schedule.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("maintenance schedule not found")
		}
		return err
	}

	if schedule.AssetID != "" {
		_, err := s.assetRepo.GetAssetByID(schedule.AssetID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("asset not found")
			}
			return err
		}
	}

	return s.repo.UpdateMaintenanceSchedule(schedule)
}

func (s *maintenanceScheduleService) DeleteMaintenanceSchedule(scheduleID string) error {

	_, err := s.repo.GetMaintenanceScheduleByID(scheduleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("maintenance schedule not found")
		}
		return err
	}
	return s.repo.DeleteMaintenanceSchedule(scheduleID)
}
