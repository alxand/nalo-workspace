package models

import (
	"time"
)

type DailyTask struct {
	ID                int64     `gorm:"primaryKey" json:"id"`
	UserID            int64     `gorm:"not null;index" json:"user_id" validate:"required"`
	Day               string    `json:"day"`
	Date              time.Time `json:"date"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
	Status            string    `json:"status"`
	Score             int       `json:"score"`
	ProductivityScore int       `json:"productivity_score"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	// Relationships
	User         User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Deliverables []Deliverable  `gorm:"foreignKey:TaskID" json:"deliverables"`
	Activities   []Activity     `gorm:"foreignKey:TaskID" json:"activities"`
	ProductFocus []ProductFocus `gorm:"foreignKey:TaskID" json:"product_focus"`
	NextSteps    []NextStep     `gorm:"foreignKey:TaskID" json:"next_steps"`
	Challenges   []Challenge    `gorm:"foreignKey:TaskID" json:"challenges"`
	Notes        []Note         `gorm:"foreignKey:TaskID" json:"notes"`
	Comments     []Comment      `gorm:"foreignKey:TaskID" json:"comments"`
}
