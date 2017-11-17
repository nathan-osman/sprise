package server

import (
	"fmt"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
	"github.com/nathan-osman/sprise/db"
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

type editBucketForm struct {
	Name            string
	Endpoint        string
	AccessKey       string
	SecretAccessKey string
}

// editBucket provides a form for creating or editing buckets.
func (s *Server) editBucket(w http.ResponseWriter, r *http.Request) {
	s.conn.Transaction(func(conn *db.Conn) error {
		var (
			id   = mux.Vars(r)["id"]
			b    = &db.Bucket{}
			form = &editBucketForm{}
			ctx  = pongo2.Context{
				"form": form,
			}
		)
		if len(id) != 0 {
			if err := conn.Find(b, id).Error; err != nil {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return nil
			}
			copyStruct(b, form)
			ctx["icon"] = "write"
			ctx["title"] = fmt.Sprintf("Edit %s", b.Name)
			ctx["subtitle"] = "Modify an existing bucket"
			ctx["action"] = "Save"
		} else {
			ctx["icon"] = "plus"
			ctx["title"] = "New Bucket"
			ctx["subtitle"] = "Create a new bucket"
			ctx["action"] = "Create"
		}
		if r.Method == http.MethodPost {
			parseForm(r.Form, form)
			copyStruct(form, b)
			if err := conn.Save(b).Error; err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return nil
			}
			http.Redirect(w, r, "/admin/buckets", http.StatusFound)
			return nil
		}
		s.render(w, r, "buckets/edit.html", ctx)
		return nil
	})
}
