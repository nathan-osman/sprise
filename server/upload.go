package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/flosch/pongo2"
)

// upload displays the page for uploading photos.
func (s *Server) upload(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "upload.html", pongo2.Context{
		"title": "Upload",
	})
}

// uploadAjax processes file uploads from the uploader.
func (s *Server) uploadAjax(w http.ResponseWriter, r *http.Request) {
	var errorMessage string
	for {
		if r.Method != http.MethodPost {
			errorMessage = "request method must be POST"
			break
		}
		if err := r.ParseMultipartForm(2 << 24); err != nil {
			errorMessage = err.Error()
			break
		}
		f, handler, err := r.FormFile("qqfile")
		if err != nil {
			errorMessage = err.Error()
			break
		}
		defer f.Close()
		_ = handler
		break
	}
	o := make(map[string]interface{})
	if len(errorMessage) != 0 {
		o["error"] = errorMessage
	} else {
		o["success"] = true
	}
	b, err := json.Marshal(o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
