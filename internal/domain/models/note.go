package models

type Note struct {
	ID     int64  `gorm:"primaryKey" json:"id"`
	TaskID int64  `json:"task_id"`
	Text   string `json:"text"`
}
