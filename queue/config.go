package queue

import (
	"github.com/nathan-osman/sprise/db"
)

// Config provides parameters for the upload queue.
type Config struct {
	QueueDir string
	Conn     *db.Conn
}
