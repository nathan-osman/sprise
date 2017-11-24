package db

// Bucket stores the information necessary to upload and retrieve data from an
// Amazon S3 bucket or Digital Ocean Spaces.
type Bucket struct {
	ID              int64
	Name            string `gorm:"not null"`
	Region          string `gorm:"not null"`
	Endpoint        string `gorm:"not null"`
	AccessKey       string `gorm:"not null"`
	SecretAccessKey string `gorm:"not null"`
}

// Buckets retrieves an ordered slice of all buckets in the database.
func (c *Conn) Buckets() ([]*Bucket, error) {
	var buckets []*Bucket
	if err := c.Order("name").Find(&buckets).Error; err != nil {
		return nil, err
	}
	return buckets, nil
}

// String returns a string representation of the bucket.
func (b *Bucket) String() string {
	return b.Name
}
