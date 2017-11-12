package queue

import (
	"github.com/minio/minio-go"
	"github.com/nathan-osman/sprise/db"
)

// Queue manages uploads to S3 / Spaces.
type Queue struct {
	client *minio.Client
	conn   *db.Conn
}

// New creates a new upload queue.
func New(cfg *Config) (*Queue, error) {
	c, err := minio.New(cfg.Endpoint, cfg.AccessKey, cfg.SecretAccessKey, true)
	if err != nil {
		return nil, err
	}
	q := &Queue{
		client: c,
		conn:   cfg.Conn,
	}
	return q, nil
}

// Close shuts down the queue.
func (q *Queue) Close() {
	//...
}
