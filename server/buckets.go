package server

import (
	"net/http"

	"github.com/flosch/pongo2"
)

// buckets displays a list of all currently registered buckets and provides a
// button for adding new ones.
func (s *Server) buckets(w http.ResponseWriter, r *http.Request) {
	buckets, err := s.conn.Buckets()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.render(w, r, "buckets/index.html", pongo2.Context{
		"icon":     "cube",
		"title":    "Buckets",
		"subtitle": "Manage photo storage",
		"buckets":  buckets,
	})
}
