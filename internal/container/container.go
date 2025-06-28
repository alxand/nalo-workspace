package container

import (
	"github.com/alxand/nalo-workspace/internal/api/auth"
	"github.com/alxand/nalo-workspace/internal/api/company"
	"github.com/alxand/nalo-workspace/internal/api/continent"
	"github.com/alxand/nalo-workspace/internal/api/country"
	"github.com/alxand/nalo-workspace/internal/api/dailytask"
	"github.com/alxand/nalo-workspace/internal/api/user"
	"github.com/alxand/nalo-workspace/internal/config"
	"github.com/alxand/nalo-workspace/internal/domain/interfaces"
	"github.com/alxand/nalo-workspace/internal/pkg/logger"
	postgresRepo "github.com/alxand/nalo-workspace/internal/repository/postgres"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Container holds all application dependencies
type Container struct {
	Config *config.Config
	Logger *zap.Logger
	DB     *gorm.DB

	// Repositories
	DailyTaskRepo interfaces.DailyTaskInterface
	UserRepo      interfaces.UserInterface
	ContinentRepo interfaces.ContinentInterface
	CountryRepo   interfaces.CountryInterface
	CompanyRepo   interfaces.CompanyInterface

	// Services
	AuthService *auth.Service

	// Handlers
	DailyTaskHandler *dailytask.TaskHandler
	AuthHandler      *auth.AuthHandler
	UserHandler      *user.UserHandler
	ContinentHandler *continent.ContinentHandler
	CountryHandler   *country.CountryHandler
	CompanyHandler   *company.CompanyHandler
}

// NewContainer creates a new container with all dependencies initialized
func NewContainer() (*Container, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// Initialize logger
	if err := logger.Init(cfg.Log.Level, cfg.Log.Format); err != nil {
		return nil, err
	}
	log := logger.Get()

	// Initialize database
	db, err := initDatabase(cfg.Database)
	if err != nil {
		return nil, err
	}

	// Auto-migrate models
	if err := config.RunMigrations(db); err != nil {
		return nil, err
	}

	// Initialize repositories
	dailyTaskRepo := postgresRepo.NewDailyTaskRepository(db)
	userRepo := postgresRepo.NewUserRepository(db)
	continentRepo := postgresRepo.NewContinentRepository(db)
	countryRepo := postgresRepo.NewCountryRepository(db)
	companyRepo := postgresRepo.NewCompanyRepository(db)

	// Initialize services
	authService := auth.NewService(cfg.JWT, userRepo, log)

	// Initialize handlers
	dailyTaskHandler := dailytask.NewTDailyTaskHandler(dailyTaskRepo, log)
	authHandler := auth.NewAuthHandler(authService, log)
	userHandler := user.NewUserHandler(userRepo, log)
	continentHandler := continent.NewContinentHandler(continentRepo, log)
	countryHandler := country.NewCountryHandler(countryRepo, log)
	companyHandler := company.NewCompanyHandler(companyRepo, log)

	return &Container{
		Config:           cfg,
		Logger:           log,
		DB:               db,
		DailyTaskRepo:    dailyTaskRepo,
		UserRepo:         userRepo,
		ContinentRepo:    continentRepo,
		CountryRepo:      countryRepo,
		CompanyRepo:      companyRepo,
		AuthService:      authService,
		DailyTaskHandler: dailyTaskHandler,
		AuthHandler:      authHandler,
		UserHandler:      userHandler,
		ContinentHandler: continentHandler,
		CountryHandler:   countryHandler,
		CompanyHandler:   companyHandler,
	}, nil
}

// Close closes all resources in the container
func (c *Container) Close() error {
	// Close database connection
	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// initDatabase initializes the database connection
func initDatabase(dbConfig config.DatabaseConfig) (*gorm.DB, error) {
	var connector config.DBConnector

	switch dbConfig.Driver {
	case "postgres":
		connector = config.NewPostgresConfig(dbConfig.DSN)
	case "sqlite":
		connector = config.NewSQLiteConfig(dbConfig.DSN)
	default:
		return nil, config.ErrUnsupportedDatabase
	}

	return connector.Connect()
}
