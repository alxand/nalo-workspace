package models

type ProductFocus struct {
	ID     int64  `gorm:"primaryKey" json:"id"`
	TaskID int64  `json:"task_id"`
	Area   string `json:"area"`
}
