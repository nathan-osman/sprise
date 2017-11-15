package db

// Tag acts as a classifier for photos, allowing subjects and objects in the
// photo to be identified.
type Tag struct {
	ID   int64
	Name string `gorm:"not null"`
}
