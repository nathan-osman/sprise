package db

import (
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Conn maintains a connection to the database.
type Conn struct {
	*gorm.DB
}

// New creates a new connection to the database.
func New(cfg *Config) (*Conn, error) {
	g, err := gorm.Open(cfg.Driver, cfg.Args)
	if err != nil {
		return nil, err
	}
	c := &Conn{g}
	return c, nil
}
