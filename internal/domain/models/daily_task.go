package models

import "time"

type DailyTask struct {
	ID                int64     `json:"id"`
	Day               string    `json:"day"`
	Date              time.Time `json:"date"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
	Status            string    `json:"status"`
	Score             int       `json:"score"`
	ProductivityScore int       `json:"productivity_score"`

	Deliverables []Deliverable
	Activities   []Activity
	ProductFocus []ProductFocus
	NextSteps    []NextStep
	Challenges   []Challenge
	Notes        []Note
	Comments     []Comment
}
