package db

import (
	"fmt"
	"path"
	"strconv"
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

// Path returns the path to the upload on disk.
func (u *Upload) Path(queueDir string) string {
	return path.Join(queueDir, strconv.FormatInt(u.ID, 10))
}

// ThumbPath returns the path to the thumbnail on disk.
func (u *Upload) ThumbPath(queueDir string) string {
	return fmt.Sprintf("%s-thumb", u.Path(queueDir))
}
