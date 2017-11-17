package db

// Upload represents a file waiting to be stored in a bucket.
type Upload struct {
	ID       int64
	Filename string `gorm:"not null"`
}
