package models

import "time"

type TaskStatus string

const (
	TaskStatusCreated    TaskStatus = "created"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusReviewable TaskStatus = "reviewable"
	TaskStatusCompleted  TaskStatus = "completed"
)

type Task struct {
	ID          string     `gorm:"type:char(36);primaryKey"`
	Title       string     `gorm:"not null"`
	Status      TaskStatus `gorm:"not null"`
	DueDate     time.Time  `gorm:"type:date;not null"`
	Description string     `gorm:"type:text"`
}
