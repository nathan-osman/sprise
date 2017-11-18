package server

import (
	"net"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/nathan-osman/sprise/db"
	"github.com/nathan-osman/sprise/queue"
	"github.com/sirupsen/logrus"
)

// Server provides the web interface for the application.
type Server struct {
	queueDir    string
	listener    net.Listener
	router      *mux.Router
	store       *sessions.CookieStore
	templateSet *pongo2.TemplateSet
	log         *logrus.Entry
	conn        *db.Conn
	queue       *queue.Queue
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
			queueDir:    cfg.QueueDir,
			listener:    l,
			router:      mux.NewRouter(),
			store:       sessions.NewCookieStore([]byte(cfg.SecretKey)),
			templateSet: pongo2.NewSet("", &b0xLoader{}),
			log:         logrus.WithField("context", "server"),
			conn:        cfg.Conn,
			queue:       cfg.Queue,
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
	s.router.HandleFunc("/upload/ajax", s.requireUser(s.uploadAjax))
	s.router.HandleFunc("/admin/buckets", s.requireAdmin(s.buckets))
	s.router.HandleFunc("/admin/buckets/new", s.requireAdmin(s.editBucket))
	s.router.HandleFunc("/admin/buckets/{id:[0-9]+}/edit", s.requireAdmin(s.editBucket))
	s.router.HandleFunc("/admin/buckets/{id:[0-9]+}/delete", s.requireAdmin(s.deleteBucket))
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
