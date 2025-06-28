package config

import (
	"github.com/alxand/nalo-workspace/internal/domain/models"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	// Auto-migrate all models
	err := db.AutoMigrate(
		&models.Continent{},
		&models.Country{},
		&models.Company{},
		&models.User{},
		&models.DailyTask{},
		&models.Deliverable{},
		&models.Activity{},
		&models.ProductFocus{},
		&models.NextStep{},
		&models.Challenge{},
		&models.Note{},
		&models.Comment{},
	)

	if err != nil {
		return err
	}

	return nil
}
