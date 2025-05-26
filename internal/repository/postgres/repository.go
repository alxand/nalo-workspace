// internal/repository/gorm_repository.go
package repository

import (
	"log"
	"time"

	"github.com/alxand/nalo-workspace/internal/domain/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Add GORM model annotations

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.DailyTask{},
		&models.Deliverable{},
		&models.Activity{},
		&models.ProductFocus{},
		&models.NextStep{},
		&models.Challenge{},
		&models.Note{},
		&models.Comment{},
	)
}

func InitDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	if err := migrate(db); err != nil {
		log.Fatalf("failed to run migration: %v", err)
	}
	return db
}

type DailyTaskRepository struct {
	DB *gorm.DB
}

func NewDailyTaskRepository(db *gorm.DB) *DailyTaskRepository {
	return &DailyTaskRepository{DB: db}
}

func (r *DailyTaskRepository) Create(log *models.DailyTask) error {
	return r.DB.Create(log).Error
}

func (r *DailyTaskRepository) GetByDate(date string) ([]models.DailyTask, error) {
	var logs []models.DailyTask
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}
	err = r.DB.Preload("Deliverables").Preload("Activities").Preload("ProductFocus").Preload("NextSteps").
		Preload("Challenges").Preload("Notes").Preload("Comments").
		Where("date = ?", parsedDate).Find(&logs).Error
	return logs, err
}

func (r *DailyTaskRepository) Update(log *models.DailyTask) error {
	return r.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(log).Error
}

func (r *DailyTaskRepository) Delete(id int64) error {
	return r.DB.Delete(&models.DailyTask{}, id).Error
}
