package server

import (
	"github.com/nathan-osman/sprise/db"
	"github.com/nathan-osman/sprise/queue"
)

// Config provides parameters for the web server.
type Config struct {
	Addr      string
	QueueDir  string
	SecretKey string
	Conn      *db.Conn
	Queue     *queue.Queue
}
