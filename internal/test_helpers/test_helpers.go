package test_helpers

import (
	"github.com/alxand/nalo-workspace/internal/config"
	"github.com/alxand/nalo-workspace/internal/domain/interfaces"
	"github.com/alxand/nalo-workspace/internal/domain/models"
	"github.com/alxand/nalo-workspace/internal/pkg/logger"
	"github.com/alxand/nalo-workspace/internal/pkg/middleware"
	"github.com/alxand/nalo-workspace/internal/pkg/validation"
	"github.com/alxand/nalo-workspace/internal/repository/sqlite"
	"github.com/gofiber/fiber/v2"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
)

// TestDB holds the test database and repositories
type TestDB struct {
	DB            *gorm.DB
	CompanyRepo   interfaces.CompanyInterface
	ContinentRepo interfaces.ContinentInterface
	CountryRepo   interfaces.CountryInterface
	UserRepo      interfaces.UserInterface
}

// SetupTestDB creates a new SQLite in-memory database for testing
func SetupTestDB() (*TestDB, error) {
	// Initialize validation
	validation.Init()

	// Create SQLite config for in-memory database
	sqliteConfig := config.NewSQLiteConfig(":memory:")

	// Connect to database
	db, err := sqliteConfig.Connect()
	if err != nil {
		return nil, err
	}

	// Create repositories
	companyRepo := sqlite.NewCompanyRepository(db)
	continentRepo := sqlite.NewContinentRepository(db)
	countryRepo := sqlite.NewCountryRepository(db)
	userRepo := sqlite.NewUserRepository(db)

	return &TestDB{
		DB:            db,
		CompanyRepo:   companyRepo,
		ContinentRepo: continentRepo,
		CountryRepo:   countryRepo,
		UserRepo:      userRepo,
	}, nil
}

// SetupTestApp creates a new Fiber app with test configuration
func SetupTestApp() *fiber.App {
	// Create logger
	log := logger.Get()

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler(log),
	})

	// Add request logger middleware
	app.Use(fiberlogger.New())

	return app
}

// CleanupTestDB cleans up the test database
func CleanupTestDB(testDB *TestDB) error {
	if testDB.DB != nil {
		sqlDB, err := testDB.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// CreateTestContinent creates a test continent in the database
func CreateTestContinent(testDB *TestDB, name, code string) (*models.Continent, error) {
	continent := &models.Continent{
		Name:        name,
		Code:        code,
		Description: "Test continent",
	}

	err := testDB.ContinentRepo.Create(continent)
	if err != nil {
		return nil, err
	}

	return continent, nil
}

// CreateTestCountry creates a test country in the database
func CreateTestCountry(testDB *TestDB, name, code string, continentID int64) (*models.Country, error) {
	country := &models.Country{
		Name:        name,
		Code:        code,
		ContinentID: continentID,
		Description: "Test country",
	}

	err := testDB.CountryRepo.Create(country)
	if err != nil {
		return nil, err
	}

	return country, nil
}

// CreateTestCompany creates a test company in the database
func CreateTestCompany(testDB *TestDB, name, code string, countryID int64) (*models.Company, error) {
	company := &models.Company{
		Name:        name,
		Code:        code,
		CountryID:   countryID,
		Description: "Test company",
		Industry:    "Technology",
		Size:        "medium",
	}

	err := testDB.CompanyRepo.Create(company)
	if err != nil {
		return nil, err
	}

	return company, nil
}
