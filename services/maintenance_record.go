package services

import (
	"errors"
	"jaga/models"
	"jaga/repositories"
	"jaga/utils"

	"gorm.io/gorm"
)

type MaintenanceRecordService interface {
	CreateMaintenanceRecord(record *models.MaintenanceRecord) error
	GetMaintenanceRecordByID(recordID string) (*models.MaintenanceRecord, error)
	GetMaintenanceRecords(page, itemsPerPage int, sortBy, sortDir, assetID, scheduleID, status string) ([]models.MaintenanceRecord, int64, error)
	UpdateMaintenanceRecord(record *models.MaintenanceRecord) error
	UpdateMaintenanceRecordStatus(recordID, status string) error
	DeleteMaintenanceRecord(recordID string) error
}

type maintenanceRecordService struct {
	repo         repositories.MaintenanceRecordRepository
	assetRepo    repositories.AssetRepository
	scheduleRepo repositories.MaintenanceScheduleRepository
	userRepo     repositories.UserRepository
}

func NewMaintenanceRecordService(
	repo repositories.MaintenanceRecordRepository,
	assetRepo repositories.AssetRepository,
	scheduleRepo repositories.MaintenanceScheduleRepository,
	userRepo repositories.UserRepository,
) MaintenanceRecordService {
	return &maintenanceRecordService{
		repo:         repo,
		assetRepo:    assetRepo,
		scheduleRepo: scheduleRepo,
		userRepo:     userRepo,
	}
}

func (s *maintenanceRecordService) CreateMaintenanceRecord(record *models.MaintenanceRecord) error {
	_, err := s.assetRepo.GetAssetByID(record.AssetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("asset not found")
		}
		return err
	}

	if record.ScheduleID != nil && *record.ScheduleID != "" {
		_, err := s.scheduleRepo.GetMaintenanceScheduleByID(*record.ScheduleID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("maintenance schedule not found")
			}
			return err
		}
	}

	if record.PerformedBy != nil && *record.PerformedBy != "" {
		_, err := s.userRepo.GetUserByID(*record.PerformedBy)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("performed by not found")
			}
			return err
		}
	}

	if record.ID == "" {
		record.ID = utils.GenerateUUID()
	}

	return s.repo.CreateMaintenanceRecord(record)
}

func (s *maintenanceRecordService) GetMaintenanceRecordByID(recordID string) (*models.MaintenanceRecord, error) {
	record, err := s.repo.GetMaintenanceRecordByID(recordID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("maintenance record not found")
		}
		return nil, err
	}
	return record, nil
}

func (s *maintenanceRecordService) GetMaintenanceRecords(page, itemsPerPage int, sortBy, sortDir, assetID, scheduleID, status string) ([]models.MaintenanceRecord, int64, error) {
	records, totalItems, err := s.repo.GetMaintenanceRecords(page, itemsPerPage, sortBy, sortDir, assetID, scheduleID, status)
	if err != nil {
		return nil, 0, err
	}
	return records, totalItems, nil
}

func (s *maintenanceRecordService) UpdateMaintenanceRecord(record *models.MaintenanceRecord) error {
	_, err := s.repo.GetMaintenanceRecordByID(record.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("maintenance record not found")
		}
		return err
	}

	if record.AssetID != "" {
		_, err := s.assetRepo.GetAssetByID(record.AssetID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("asset not found")
			}
			return err
		}
	}

	if record.ScheduleID != nil && *record.ScheduleID != "" {
		_, err := s.scheduleRepo.GetMaintenanceScheduleByID(*record.ScheduleID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("maintenance schedule not found")
			}
			return err
		}
	} else if record.ScheduleID != nil && *record.ScheduleID == "" {

	}

	if record.PerformedBy != nil && *record.PerformedBy != "" {
		_, err := s.userRepo.GetUserByID(*record.PerformedBy)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("performed by not found")
			}
			return err
		}
	} else if record.PerformedBy != nil && *record.PerformedBy == "" {

	}

	return s.repo.UpdateMaintenanceRecord(record)
}

func (s *maintenanceRecordService) UpdateMaintenanceRecordStatus(recordID, status string) error {
	record, err := s.repo.GetMaintenanceRecordByID(recordID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("maintenance record not found")
		}
		return err
	}

	record.Status = status

	return s.repo.UpdateMaintenanceRecord(record)
}

func (s *maintenanceRecordService) DeleteMaintenanceRecord(recordID string) error {
	_, err := s.repo.GetMaintenanceRecordByID(recordID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("maintenance record not found")
		}
		return err
	}

	return s.repo.DeleteMaintenanceRecord(recordID)
}
