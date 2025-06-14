package services

import (
	"fmt"
	"jaga/dto"
	"jaga/models"
	"jaga/repositories"
	"time"

	"github.com/google/uuid"
)

func CreateAsset(input dto.CreateAsset, user models.User) error {
	if user.Role != "super_user" && user.Role != "admin" {
		return fmt.Errorf("%s is not allowed to create asset", user.Role)
	}

	asset := models.Asset{
		ID:           uuid.New().String(),
		Name:         input.Name,
		CategoryID:   input.CategoryID,
		Location:     input.Location,
		PurchaseDate: input.PurchaseDate,
		Condition:    input.Condition,
		Status:       input.Status,
		AddedBy:      user.ID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return repositories.CreateAsset(&asset)
}
