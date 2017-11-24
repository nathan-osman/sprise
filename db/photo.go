package db

import (
	"fmt"
	"strconv"
	"time"
)

// Photo represents an individual photograph and its metadata.
type Photo struct {
	ID       int64
	Date     time.Time `gorm:"not null"`
	Filename string    `gorm:"not null"`
	URL      string    `gorm:"not null"`
	ThumbURL string    `gorm:"not null"`
	Width    int       `gorm:"not null"`
	Height   int       `gorm:"not null"`
	Camera   string    `gorm:"not null"`
	Bucket   *Bucket   `gorm:"ForeignKey:BucketID"`
	BucketID int64     `sql:"type:int REFERENCES buckets(id) ON DELETE CASCADE"`
	Tags     []*Tag    `gorm:"many2many:photo_tags;"`
}

// Name returns the object name used for the photo in the bucket.
func (p *Photo) Name() string {
	return strconv.FormatInt(p.ID, 10)
}

// ThumbName returns the object name used for the photo's thumbnail in the
// bucket.
func (p *Photo) ThumbName() string {
	return fmt.Sprintf("%s-thumb", p.Name())
}
