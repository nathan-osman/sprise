package queue

import (
	"github.com/nathan-osman/sprise/db"
)

// Config provides parameters for the upload queue.
type Config struct {
	Endpoint        string
	AccessKey       string
	SecretAccessKey string
	Bucket          string
	Conn            *db.Conn
}
