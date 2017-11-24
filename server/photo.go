package server

import (
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
	"github.com/nathan-osman/sprise/db"
)

// viewPhoto displays an individual photo.
func (s *Server) viewPhoto(w http.ResponseWriter, r *http.Request) {
	var (
		id = mux.Vars(r)["id"]
		p  = &db.Photo{}
	)
	if err := s.conn.Find(p, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.render(w, r, "photos/view.html", pongo2.Context{
		"icon":     "photo",
		"title":    p.Filename,
		"subtitle": "View photograph",
		"photo":    p,
	})
}
