package interfaces

import "gorm.io/gorm"

type MigrationInterface interface {
	Migrate(db *gorm.DB) error
}
