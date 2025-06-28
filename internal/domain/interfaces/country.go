package interfaces

import "github.com/alxand/nalo-workspace/internal/domain/models"

type CountryInterface interface {
	Create(country *models.Country) error
	GetByID(id int64) (*models.Country, error)
	GetByCode(code string) (*models.Country, error)
	GetByContinent(continentID int64) ([]models.Country, error)
	GetAll() ([]models.Country, error)
	Update(country *models.Country) (*models.Country, error)
	Delete(id int64) error
}
