package repositories

import (
	"jaga/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUsers(page, itemsPerPage int, sortBy, sortDir string) ([]models.User, int64, error)
	GetUserByID(userID string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByRole(role string) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(userID string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUsers(page, itemsPerPage int, sortBy, sortDir string) ([]models.User, int64, error) {
	var users []models.User
	var totalItems int64

	if err := r.db.Model(&models.User{}).Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * itemsPerPage
	query := r.db.Offset(offset).Limit(itemsPerPage)

	if sortBy != "" && sortDir != "" {
		query = query.Order(sortBy + " " + sortDir)
	} else {
		query = query.Order("created_at desc")
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, totalItems, nil
}

func (r *userRepository) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByRole(role string) (*models.User, error) {
	var user models.User

	if err := r.db.Where("role = ?", role).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) UpdateUser(user *models.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) DeleteUser(userID string) error {
	if err := r.db.Delete(&models.User{}, "id = ?", userID).Error; err != nil {
		return err
	}
	return nil
}
