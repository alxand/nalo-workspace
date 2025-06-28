package repository

import (
	"time"

	"github.com/alxand/nalo-workspace/internal/domain/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	var user models.User
	err := r.DB.Preload("Country").Preload("Company").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Preload("Country").Preload("Company").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.DB.Preload("Country").Preload("Company").Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByCountry(countryID int64) ([]models.User, error) {
	var users []models.User
	err := r.DB.Preload("Country").Preload("Company").Where("country_id = ?", countryID).Find(&users).Error
	return users, err
}

func (r *UserRepository) GetByCompany(companyID int64) ([]models.User, error) {
	var users []models.User
	err := r.DB.Preload("Country").Preload("Company").Where("company_id = ?", companyID).Find(&users).Error
	return users, err
}

func (r *UserRepository) Update(user *models.User) error {
	return r.DB.Save(user).Error
}

func (r *UserRepository) Delete(id int64) error {
	return r.DB.Delete(&models.User{}, id).Error
}

func (r *UserRepository) List(limit, offset int) ([]models.User, error) {
	var users []models.User
	err := r.DB.Preload("Country").Preload("Company").Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

func (r *UserRepository) UpdateLastLogin(id int64) error {
	now := time.Now()
	return r.DB.Model(&models.User{}).Where("id = ?", id).Update("last_login", now).Error
}

func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.DB.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.DB.Model(&models.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}
