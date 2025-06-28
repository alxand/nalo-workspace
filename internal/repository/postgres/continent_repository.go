package repository

import (
	"errors"

	"github.com/alxand/nalo-workspace/internal/domain/interfaces"
	"github.com/alxand/nalo-workspace/internal/domain/models"
	"gorm.io/gorm"
)

type ContinentRepository struct {
	db *gorm.DB
}

func NewContinentRepository(db *gorm.DB) interfaces.ContinentInterface {
	return &ContinentRepository{db: db}
}

func (r *ContinentRepository) Create(continent *models.Continent) error {
	return r.db.Create(continent).Error
}

func (r *ContinentRepository) GetByID(id int64) (*models.Continent, error) {
	var continent models.Continent
	err := r.db.Preload("Countries").First(&continent, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &continent, nil
}

func (r *ContinentRepository) GetByCode(code string) (*models.Continent, error) {
	var continent models.Continent
	err := r.db.Preload("Countries").Where("code = ?", code).First(&continent).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &continent, nil
}

func (r *ContinentRepository) GetAll() ([]models.Continent, error) {
	var continents []models.Continent
	err := r.db.Preload("Countries").Find(&continents).Error
	return continents, err
}

func (r *ContinentRepository) Update(continent *models.Continent) (*models.Continent, error) {
	err := r.db.Save(continent).Error
	if err != nil {
		return nil, err
	}
	return r.GetByID(continent.ID)
}

func (r *ContinentRepository) Delete(id int64) error {
	return r.db.Delete(&models.Continent{}, id).Error
}
