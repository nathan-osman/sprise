package db

import (
	"time"
)

// Photo represents an individual photograph and its metadata.
type Photo struct {
	ID       int64
	Date     time.Time `gorm:"not null"`
	Filename string    `gorm:"not null"`
	Width    int       `gorm:"not null"`
	Height   int       `gorm:"not null"`
	Camera   string    `gorm:"not null"`
	Tags     []*Tag    `gorm:"many2many:photo_tags;"`
}
