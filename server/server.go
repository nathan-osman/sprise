package server

import (
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nathan-osman/sprise/db"
	"github.com/sirupsen/logrus"
)

// Server provides the web interface for the application.
type Server struct {
	listener  net.Listener
	router    *mux.Router
	log       *logrus.Entry
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
			log:       logrus.WithField("context", "server"),
			conn:      cfg.Conn,
			stoppedCh: make(chan bool),
		}
		server = http.Server{
			Handler: r,
		}
	)
	go func() {
		defer close(s.stoppedCh)
		defer s.log.Info("server has stopped")
		s.log.Info("starting server...")
		if err := server.Serve(l); err != nil {
			s.log.Error(err.Error())
		}
	}()
	return s, nil
}

// Close shuts down the web server.
func (s *Server) Close() {
	s.listener.Close()
	<-s.stoppedCh
}
