package interfaces

import "github.com/alxand/nalo-workspace/internal/domain/models"

type CompanyInterface interface {
	Create(company *models.Company) error
	GetByID(id int64) (*models.Company, error)
	GetByCode(code string) (*models.Company, error)
	GetByCountry(countryID int64) ([]models.Company, error)
	GetByIndustry(industry string) ([]models.Company, error)
	GetAll() ([]models.Company, error)
	Update(company *models.Company) (*models.Company, error)
	Delete(id int64) error
}
