package repositories

import (
	"errors"
	"jaga/config"
	"jaga/models"
)

func CreateAsset(asset *models.Asset) error {
	if err := config.DB.Create(asset).Error; err != nil {
		return errors.New("failed to create asset")
	}
	return nil
}
