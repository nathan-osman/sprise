package server

import (
	"net/http"

	"github.com/flosch/pongo2"
)

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	photos, err := s.conn.Photos(12)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.render(w, r, "index.html", pongo2.Context{
		"title":  "Home",
		"photos": photos,
	})
}
