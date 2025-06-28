package sqlite

import (
	"errors"

	"github.com/alxand/nalo-workspace/internal/domain/interfaces"
	"github.com/alxand/nalo-workspace/internal/domain/models"
	"gorm.io/gorm"
)

// TaskMigrator handles database migrations for SQLite
type TaskMigrator struct{}

func (m *TaskMigrator) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Continent{},
		&models.Country{},
		&models.Company{},
		&models.User{},
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

// CompanyRepository implements CompanyInterface for SQLite
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

// ContinentRepository implements ContinentInterface for SQLite
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

// CountryRepository implements CountryInterface for SQLite
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
	err := r.db.Preload("Continent").Preload("Companies").Preload("Users").Where("continent_id = ?", continentID).Find(&countries).Error
	return countries, err
}

func (r *CountryRepository) GetAll() ([]models.Country, error) {
	var countries []models.Country
	err := r.db.Preload("Continent").Preload("Companies").Preload("Users").Find(&countries).Error
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

// UserRepository implements UserInterface for SQLite
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaces.UserInterface {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Company").Preload("Country").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Company").Preload("Country").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Company").Preload("Country").Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByCountry(countryID int64) ([]models.User, error) {
	var users []models.User
	err := r.db.Preload("Company").Preload("Country").Where("country_id = ?", countryID).Find(&users).Error
	return users, err
}

func (r *UserRepository) GetByCompany(companyID int64) ([]models.User, error) {
	var users []models.User
	err := r.db.Preload("Company").Preload("Country").Where("company_id = ?", companyID).Find(&users).Error
	return users, err
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Preload("Company").Preload("Country").Find(&users).Error
	return users, err
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id int64) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *UserRepository) List(limit, offset int) ([]models.User, error) {
	var users []models.User
	err := r.db.Preload("Company").Preload("Country").Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

func (r *UserRepository) UpdateLastLogin(id int64) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("last_login", gorm.Expr("NOW()")).Error
}

func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}
