package db

import (
	"time"
)

// Upload represents a file waiting to be stored in a bucket.
type Upload struct {
	ID       int64
	Date     time.Time `gorm:"not null"`
	Filename string    `gorm:"not null"`
	Bucket   *Bucket   `gorm:"ForeignKey:BucketID"`
	BucketID int64     `sql:"type:int REFERENCES buckets(id) ON DELETE CASCADE"`
}
