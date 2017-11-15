package db

// Bucket stores the information necessary to upload and retrieve data from an
// Amazon S3 bucket or Digital Ocean Spaces.
type Bucket struct {
	ID              int64
	Endpoint        string `gorm:"not null"`
	AccessKey       string `gorm:"not null"`
	SecretAccessKey string `gorm:"not null"`
	Name            string `gorm:"not null"`
}
