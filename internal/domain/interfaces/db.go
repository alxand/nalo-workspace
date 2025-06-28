package interfaces

import "gorm.io/gorm"

type DBConnectionInterface interface {
	Connect() error
	GetDB() *gorm.DB
}
