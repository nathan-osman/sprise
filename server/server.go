package server

import (
	"net"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/nathan-osman/sprise/db"
	"github.com/sirupsen/logrus"
)

// Server provides the web interface for the application.
type Server struct {
	listener    net.Listener
	router      *mux.Router
	store       *sessions.CookieStore
	templateSet *pongo2.TemplateSet
	log         *logrus.Entry
	conn        *db.Conn
	stoppedCh   chan bool
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
			listener:    l,
			router:      mux.NewRouter(),
			store:       sessions.NewCookieStore([]byte(cfg.SecretKey)),
			templateSet: pongo2.NewSet("", &b0xLoader{}),
			log:         logrus.WithField("context", "server"),
			conn:        cfg.Conn,
			stoppedCh:   make(chan bool),
		}
		server = http.Server{
			Handler: r,
		}
	)
	s.router.HandleFunc("/", s.index)
	s.router.HandleFunc("/login", s.login)
	s.router.HandleFunc("/logout", s.requireUser(s.logout))
	s.router.HandleFunc("/upload", s.requireUser(s.upload))
	r.PathPrefix("/static").Handler(http.FileServer(HTTP))
	r.PathPrefix("/").Handler(s)
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

// ServeHTTP does preparatory work for the handlers, attempting to load the
// current user from the database.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	r = s.loadUser(r)
	s.router.ServeHTTP(w, r)
}

// Close shuts down the web server.
func (s *Server) Close() {
	s.listener.Close()
	<-s.stoppedCh
}
