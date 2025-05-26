package models

type Comment struct {
	ID      int64  `gorm:"primaryKey" json:"id"`
	TaskID  int64  `json:"task_id"`
	Author  string `json:"author"` // Manager, BD, MD, etc.
	Content string `json:"content"`
}
