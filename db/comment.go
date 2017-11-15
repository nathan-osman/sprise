package db

import (
	"time"
)

// Comment represents feedback left by users regarding a photograph.
type Comment struct {
	ID      int64
	Date    time.Time `gorm:"not null"`
	Text    string    `gorm:"not null"`
	Photo   *User     `gorm:"ForeignKey:PhotoID"`
	PhotoID int64     `sql:"type:int REFERENCES photos(id) ON DELETE CASCADE"`
	User    *User     `gorm:"ForeignKey:UserID"`
	UserID  int64     `sql:"type:int REFERENCES users(id) ON DELETE CASCADE"`
}
