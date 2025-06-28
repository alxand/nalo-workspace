package interfaces

import "github.com/alxand/nalo-workspace/internal/domain/models"

type UserInterface interface {
	Create(user *models.User) error
	GetByID(id int64) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByCountry(countryID int64) ([]models.User, error)
	GetByCompany(companyID int64) ([]models.User, error)
	Update(user *models.User) error
	Delete(id int64) error
	List(limit, offset int) ([]models.User, error)
	UpdateLastLogin(id int64) error
	ExistsByEmail(email string) (bool, error)
	ExistsByUsername(username string) (bool, error)
}
