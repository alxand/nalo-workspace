package models

import (
	"time"
)

type Continent struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null;unique" json:"name" validate:"required"`
	Code        string    `gorm:"not null;unique;size:2" json:"code" validate:"required,len=2"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Countries []Country `gorm:"foreignKey:ContinentID" json:"countries,omitempty" validate:"-"`
}
