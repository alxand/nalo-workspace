package main

import (
	"fmt"
	"log"

	"github.com/alxand/nalo-workspace/internal/config"
	"github.com/alxand/nalo-workspace/internal/domain/models"
	"github.com/alxand/nalo-workspace/internal/pkg/logger"
	postgresRepo "github.com/alxand/nalo-workspace/internal/repository/postgres"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize logger
	if err := logger.Init(cfg.Log.Level, cfg.Log.Format); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	log := logger.Get()

	// Initialize database
	db, err := initDatabase(cfg.Database)
	if err != nil {
		log.Fatal("Failed to initialize database", zap.Error(err))
	}

	// Run migrations
	if err := config.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations", zap.Error(err))
	}

	// Seed continents
	log.Info("Seeding continents...")
	if err := seedContinents(db); err != nil {
		log.Fatal("Failed to seed continents", zap.Error(err))
	}

	// Seed countries
	log.Info("Seeding countries...")
	if err := seedCountries(db); err != nil {
		log.Fatal("Failed to seed countries", zap.Error(err))
	}

	// Seed companies
	log.Info("Seeding companies...")
	if err := seedCompanies(db); err != nil {
		log.Fatal("Failed to seed companies", zap.Error(err))
	}

	// Create user repository
	userRepo := postgresRepo.NewUserRepository(db)

	// Check if admin user already exists
	exists, err := userRepo.ExistsByEmail("admin@nalo-workspace.com")
	if err != nil {
		log.Fatal("Failed to check if admin exists", zap.Error(err))
	}

	if exists {
		log.Info("Admin user already exists")
		return
	}

	// Create admin user
	adminUser := &models.User{
		Email:     "admin@nalo-workspace.com",
		Username:  "admin",
		Password:  "admin123456", // Will be hashed by GORM hook
		FirstName: "Admin",
		LastName:  "User",
		Role:      models.RoleAdmin,
		IsActive:  true,
	}

	if err := userRepo.Create(adminUser); err != nil {
		log.Fatal("Failed to create admin user", zap.Error(err))
	}

	log.Info("Admin user created successfully",
		zap.String("email", adminUser.Email),
		zap.String("username", adminUser.Username),
		zap.Int64("user_id", adminUser.ID),
	)

	fmt.Println("✅ Admin user created successfully!")
	fmt.Printf("Email: %s\n", adminUser.Email)
	fmt.Printf("Username: %s\n", adminUser.Username)
	fmt.Printf("Password: admin123456\n")
	fmt.Printf("Role: %s\n", adminUser.Role)
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
