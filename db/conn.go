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
		&Upload{},
		&Tag{},
		&Photo{},
		&Comment{},
		&Album{},
	).Error
}

// Transaction executes the provided function as a single transaction. If an
// error is returned, the transaction is rolled back. If not, it is committed.
func (c *Conn) Transaction(fn func(*Conn) error) error {
	conn := c.Begin()
	if err := fn(&Conn{DB: conn, log: c.log}); err != nil {
		conn.Rollback()
		return err
	}
	conn.Commit()
	return nil
}
