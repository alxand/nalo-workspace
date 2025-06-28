package interfaces

import "github.com/alxand/nalo-workspace/internal/domain/models"

type DailyTaskInterface interface {
	Create(task *models.DailyTask) error
	GetByID(id int64) (*models.DailyTask, error)
	GetByDate(date string) ([]models.DailyTask, error)
	GetByDateAndUser(date string, userID int64) ([]models.DailyTask, error)
	Update(task *models.DailyTask) (*models.DailyTask, error)
	Delete(id int64) error
	List(limit, offset int) ([]models.DailyTask, error)
}
