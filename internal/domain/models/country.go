package models

import (
	"time"
)

type Country struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name" validate:"required"`
	Code        string    `gorm:"not null;unique;size:3" json:"code" validate:"required,len=3"`
	ContinentID int64     `gorm:"not null" json:"continent_id" validate:"required"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Continent Continent `gorm:"foreignKey:ContinentID" json:"continent,omitempty" validate:"-"`
	Companies []Company `gorm:"foreignKey:CountryID" json:"companies,omitempty" validate:"-"`
	Users     []User    `gorm:"foreignKey:CountryID" json:"users,omitempty" validate:"-"`
}
