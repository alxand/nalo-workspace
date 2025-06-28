package repository

import (
	"errors"

	"github.com/alxand/nalo-workspace/internal/domain/interfaces"
	"github.com/alxand/nalo-workspace/internal/domain/models"
	"gorm.io/gorm"
)

type CountryRepository struct {
	db *gorm.DB
}

func NewCountryRepository(db *gorm.DB) interfaces.CountryInterface {
	return &CountryRepository{db: db}
}

func (r *CountryRepository) Create(country *models.Country) error {
	return r.db.Create(country).Error
}

func (r *CountryRepository) GetByID(id int64) (*models.Country, error) {
	var country models.Country
	err := r.db.Preload("Continent").Preload("Companies").Preload("Users").First(&country, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &country, nil
}

func (r *CountryRepository) GetByCode(code string) (*models.Country, error) {
	var country models.Country
	err := r.db.Preload("Continent").Preload("Companies").Preload("Users").Where("code = ?", code).First(&country).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &country, nil
}

func (r *CountryRepository) GetByContinent(continentID int64) ([]models.Country, error) {
	var countries []models.Country
	err := r.db.Preload("Continent").Preload("Companies").Where("continent_id = ?", continentID).Find(&countries).Error
	return countries, err
}

func (r *CountryRepository) GetAll() ([]models.Country, error) {
	var countries []models.Country
	err := r.db.Preload("Continent").Preload("Companies").Find(&countries).Error
	return countries, err
}

func (r *CountryRepository) Update(country *models.Country) (*models.Country, error) {
	err := r.db.Save(country).Error
	if err != nil {
		return nil, err
	}
	return r.GetByID(country.ID)
}

func (r *CountryRepository) Delete(id int64) error {
	return r.db.Delete(&models.Country{}, id).Error
}
