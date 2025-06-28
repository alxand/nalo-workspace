package interfaces

import "github.com/alxand/nalo-workspace/internal/domain/models"

type ContinentInterface interface {
	Create(continent *models.Continent) error
	GetByID(id int64) (*models.Continent, error)
	GetByCode(code string) (*models.Continent, error)
	GetAll() ([]models.Continent, error)
	Update(continent *models.Continent) (*models.Continent, error)
	Delete(id int64) error
}
