package config

import (
	"errors"
	"fmt"
	"log" // Keep log for critical errors in main application startup

	// For os.Getenv and os.LookupEnv
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Error definitions
var (
	ErrUnsupportedDatabase = errors.New("unsupported database driver")
)

// DBConnector defines the interface for database connection types.
// This allows us to treat Postgres and SQLite connections polymorphically.
type DBConnector interface {
	Connect() (*gorm.DB, error)
}

// PostgresConfig holds configuration for a PostgreSQL database connection.
type PostgresConfig struct {
	DSN string
}

// NewPostgresConfig creates a new PostgresConfig instance.
func NewPostgresConfig(dsn string) *PostgresConfig {
	return &PostgresConfig{
		DSN: dsn,
	}
}

// Connect establishes a connection to the PostgreSQL database.
func (p *PostgresConfig) Connect() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(p.DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres database: %w", err)
	}

	if err := RunMigrations(db); err != nil {
		return nil, fmt.Errorf("gorm postgres failed to run migration: %w", err)
	}

	return db, nil
}

// SQLiteConfig holds configuration for a SQLite database connection.
type SQLiteConfig struct {
	DSN string
}

// NewSQLiteConfig creates a new SQLiteConfig instance, typically for testing.
func NewSQLiteConfig(dsn string) *SQLiteConfig {
	return &SQLiteConfig{
		DSN: dsn,
	}
}

// Connect establishes a connection to the SQLite database.
func (s *SQLiteConfig) Connect() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(s.DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to sqlite database: %w", err)
	}

	if err := RunMigrations(db); err != nil {
		return nil, fmt.Errorf("gorm sqlite failed to run migration: %w", err)
	}

	return db, nil
}

// // --- Helper Functions (assuming they exist or need to be defined) ---

// // GetDSN retrieves the PostgreSQL DSN. In a real application, this would
// // likely read from environment variables or a configuration file.
// func GetDSN() string {
// 	// Example: read from environment variable. Adjust as per your config strategy.
// 	dsn := os.Getenv("DATABASE_URL")
// 	if dsn == "" {
// 		// Fallback for development if not set, or panic/error depending on strictness
// 		log.Println("DATABASE_URL environment variable not set, using default for development. Please set it for production.")
// 		return "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Shanghai"
// 	}
// 	return dsn
// }

// // GetTestDSN retrieves the SQLite DSN, typically for testing.
// func GetTestDSN() string {
// 	// Example: always return in-memory for tests, or allow override via env var
// 	if val, ok := os.LookupEnv("TEST_DATABASE_URL"); ok {
// 		return val
// 	}
// 	return ":memory:" // Default to in-memory for most tests
// }

// InitDB is a utility function to connect to a database based on the provided connector.
// It's meant for application startup. If a connection fails, it will log.Fatal.
func InitDB(connector DBConnector) *gorm.DB {
	db, err := connector.Connect()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	return db
}

// Example usage in your main function (not part of config package, but demonstrates use)
/*
func main() {
    // For production (Postgres)
    pgConfig := NewPostgresConfig()
    appDB := InitDB(pgConfig) // This will panic if connection fails

    // For testing (SQLite)
    // sqConfig := NewSQLiteConfig()
    // testDB := InitDB(sqConfig)

    // Use appDB throughout your application
    // defer func() {
    //     sqlDB, _ := appDB.DB()
    //     sqlDB.Close()
    // }()
}
*/
