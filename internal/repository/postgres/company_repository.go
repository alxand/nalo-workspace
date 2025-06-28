package repository

import (
	"errors"

	"github.com/alxand/nalo-workspace/internal/domain/interfaces"
	"github.com/alxand/nalo-workspace/internal/domain/models"
	"gorm.io/gorm"
)

type CompanyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) interfaces.CompanyInterface {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) Create(company *models.Company) error {
	return r.db.Create(company).Error
}

func (r *CompanyRepository) GetByID(id int64) (*models.Company, error) {
	var company models.Company
	err := r.db.Preload("Country").Preload("Users").First(&company, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &company, nil
}

func (r *CompanyRepository) GetByCode(code string) (*models.Company, error) {
	var company models.Company
	err := r.db.Preload("Country").Preload("Users").Where("code = ?", code).First(&company).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &company, nil
}

func (r *CompanyRepository) GetByCountry(countryID int64) ([]models.Company, error) {
	var companies []models.Company
	err := r.db.Preload("Country").Preload("Users").Where("country_id = ?", countryID).Find(&companies).Error
	return companies, err
}

func (r *CompanyRepository) GetByIndustry(industry string) ([]models.Company, error) {
	var companies []models.Company
	err := r.db.Preload("Country").Preload("Users").Where("industry = ?", industry).Find(&companies).Error
	return companies, err
}

func (r *CompanyRepository) GetAll() ([]models.Company, error) {
	var companies []models.Company
	err := r.db.Preload("Country").Preload("Users").Find(&companies).Error
	return companies, err
}

func (r *CompanyRepository) Update(company *models.Company) (*models.Company, error) {
	err := r.db.Save(company).Error
	if err != nil {
		return nil, err
	}
	return r.GetByID(company.ID)
}

func (r *CompanyRepository) Delete(id int64) error {
	return r.db.Delete(&models.Company{}, id).Error
}
