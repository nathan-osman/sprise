package server

import (
	"github.com/nathan-osman/sprise/db"
)

// Config provides parameters for the web server.
type Config struct {
	Addr      string
	SecretKey string
	Conn      *db.Conn
}
