package queue

import (
	"github.com/nathan-osman/sprise/db"
)

// Queue manages uploads to S3 / Spaces.
type Queue struct {
	conn *db.Conn
}

// New creates a new upload queue.
func New(cfg *Config) (*Queue, error) {
	q := &Queue{
		conn: cfg.Conn,
	}
	return q, nil
}

// Close shuts down the queue.
func (q *Queue) Close() {
	//...
}
