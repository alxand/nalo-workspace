// internal/repository/gorm_repository.go
package repository

import (
	"time"

	"github.com/alxand/nalo-workspace/internal/domain/models"
	"gorm.io/gorm"
)

// Add GORM model annotations

type TaskMigrator struct{}

func (m *TaskMigrator) Migrate(db *gorm.DB) error {
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

func (r *DailyTaskRepository) Update(log *models.DailyTask) (*models.DailyTask, error) {
	err := r.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(log).Error
	if err != nil {
		return nil, err
	}

	// Reload the updated task with all associations
	var updated models.DailyTask
	if err := r.DB.Preload("Deliverables").
		Preload("Activities").
		Preload("ProductFocus").
		Preload("NextSteps").
		Preload("Challenges").
		Preload("Notes").
		Preload("Comments").
		First(&updated, log.ID).Error; err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *DailyTaskRepository) Delete(id int64) error {
	return r.DB.Delete(&models.DailyTask{}, id).Error
}

func (r *DailyTaskRepository) GetByID(id int64) (*models.DailyTask, error) {
	var task models.DailyTask
	err := r.DB.First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *DailyTaskRepository) GetByDateAndUser(date string, userID int64) ([]models.DailyTask, error) {
	var tasks []models.DailyTask
	err := r.DB.Where("date = ? AND user_id = ?", date, userID).Find(&tasks).Error
	return tasks, err
}

func (r *DailyTaskRepository) List(limit, offset int) ([]models.DailyTask, error) {
	var tasks []models.DailyTask
	err := r.DB.Limit(limit).Offset(offset).Find(&tasks).Error
	return tasks, err
}
