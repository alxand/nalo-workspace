package models

type Challenge struct {
	ID     int64  `gorm:"primaryKey" json:"id"`
	TaskID int64  `json:"task_id"`
	Issue  string `json:"issue"`
}
