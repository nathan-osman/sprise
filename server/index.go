package server

import (
	"net/http"

	"github.com/flosch/pongo2"
)

func (s *Server) index(w http.ResponseWriter, req *http.Request) {
	s.render(w, req, "index.html", pongo2.Context{
		"title": "Home",
	})
}
