package db

// Album represents a logical grouping of photos.
type Album struct {
	ID     int64
	Title  string   `gorm:"not null"`
	Photos []*Photo `gorm:"many2many:album_photos;"`
}
