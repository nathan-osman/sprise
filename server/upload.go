package server

import (
	"net/http"

	"github.com/flosch/pongo2"
)

// upload displays the page for uploading photos.
func (s *Server) upload(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "upload.html", pongo2.Context{
		"title": "Upload",
	})
}
