package server

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/flosch/pongo2"
	"github.com/nathan-osman/sprise/db"
)

// upload displays the page for uploading photos.
func (s *Server) upload(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "upload.html", pongo2.Context{
		"icon":     "upload",
		"title":    "Upload",
		"subtitle": "Transfer files to Sprise",
	})
}

// saveUpload creates an upload in the database and copies the contents to the
// upload queue directory.
func (s *Server) saveUpload(r io.Reader, filename string) error {
	return s.conn.Transaction(func(conn *db.Conn) error {
		u := &db.Upload{
			Date:     time.Now(),
			Filename: filename,
		}
		if err := s.conn.Save(u).Error; err != nil {
			return err
		}
		if err := os.MkdirAll(s.queueDir, 0755); err != nil {
			return err
		}
		f, err := os.Create(path.Join(s.queueDir, strconv.FormatInt(u.ID, 10)))
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(f, r)
		return err
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
		if err := s.saveUpload(f, handler.Filename); err != nil {
			errorMessage = err.Error()
			break
		}
		s.queue.Trigger()
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
