package db

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Conn maintains a connection to the database.
type Conn struct {
	*gorm.DB
	log *logrus.Entry
}

// New creates a new connection to the database.
func New(cfg *Config) (*Conn, error) {
	g, err := gorm.Open(cfg.Driver, cfg.Args)
	if err != nil {
		return nil, err
	}
	c := &Conn{
		DB:  g,
		log: logrus.WithField("context", "db"),
	}
	return c, nil
}

// Migrate performs a database migration, ensuring all required tables and
// columns exist.
func (c *Conn) Migrate() error {
	c.log.Info("performing migrations...")
	return c.AutoMigrate(
		&User{},
		&Bucket{},
		&Tag{},
		&Photo{},
		&Comment{},
		&Album{},
	).Error
}
