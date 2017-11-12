package server

import (
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nathan-osman/sprise/db"
)

// Server provides the web interface for the application.
type Server struct {
	listener  net.Listener
	router    *mux.Router
	conn      *db.Conn
	stoppedCh chan bool
}

// New creates a new server instance.
func New(cfg *Config) (*Server, error) {
	l, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, err
	}
	var (
		r = mux.NewRouter()
		s = &Server{
			listener:  l,
			router:    r,
			conn:      cfg.Conn,
			stoppedCh: make(chan bool),
		}
		server = http.Server{
			Handler: r,
		}
	)
	go func() {
		defer close(s.stoppedCh)
		if err := server.Serve(l); err != nil {
			//...
		}
	}()
	return s, nil
}

// Close shuts down the web server.
func (s *Server) Close() {
	s.listener.Close()
	<-s.stoppedCh
}
