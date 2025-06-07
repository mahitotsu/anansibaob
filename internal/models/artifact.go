package models

import "time"

type Artifact struct {
	ID          string    `gorm:"type:char(36);primaryKey"`
	Name        string    `gorm:"not null;uniqueindex:idx_name_version"`
	Version     string    `gorm:"not null;uniqueindex:idx_name_version"`
	ReleaseDate time.Time `gorm:"type:date;not null"`
	Description string    `gorm:"type:text"`
}
