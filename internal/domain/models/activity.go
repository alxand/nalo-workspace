package models

type Activity struct {
	ID     int64  `gorm:"primaryKey" json:"id"`
	TaskID int64  `json:"task_id"`
	Name   string `json:"name"`
}
