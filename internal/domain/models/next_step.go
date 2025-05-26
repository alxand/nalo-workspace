package models

type NextStep struct {
	ID     int64  `gorm:"primaryKey" json:"id"`
	TaskID int64  `json:"task_id"`
	Step   string `json:"step"`
}
