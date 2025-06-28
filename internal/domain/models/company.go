package models

import (
	"time"
)

type Company struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name" validate:"required"`
	Code        string    `gorm:"unique" json:"code"`
	CountryID   int64     `gorm:"not null" json:"country_id" validate:"required"`
	Description string    `json:"description"`
	Website     string    `json:"website"`
	Industry    string    `json:"industry"`
	Size        string    `json:"size" validate:"company_size"` // small, medium, large, enterprise
	Founded     int       `json:"founded"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Country Country `gorm:"foreignKey:CountryID" json:"country,omitempty" validate:"-"`
	Users   []User  `gorm:"foreignKey:CompanyID" json:"users,omitempty" validate:"-"`
}
