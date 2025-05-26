package models

type Deliverable struct {
	ID     int64  `gorm:"primaryKey" json:"id"`
	TaskID int64  `json:"task_id"`
	Item   string `json:"item"`
}
