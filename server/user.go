package server

import (
	"context"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/nathan-osman/sprise/db"
)

const (
	sessionName   = "session"
	sessionUserID = "userID"

	contextUser = "user"

	errInvalidCredentials = "invalid username or password"
)

// loadUser checks the session for the ID of the current user and attempts to
// load the model and add it to the request context.
func (s *Server) loadUser(req *http.Request) *http.Request {
	session, _ := s.store.Get(req, sessionName)
	if v, _ := session.Values[sessionUserID]; v != nil {
		u := &db.User{}
		if err := s.conn.First(u, v).Error; err == nil {
			req = req.WithContext(context.WithValue(req.Context(), contextUser, u))
		}
	}
	return req
}

type loginForm struct {
	Username string `form:"required"`
	Password string `form:"required"`
}

// login displays the login form and verifies the supplied credentials upon
// submission.
func (s *Server) login(w http.ResponseWriter, req *http.Request) {
	var (
		form = &loginForm{}
		ctx  = pongo2.Context{
			"title": "Login",
			"form":  form,
		}
	)
	if req.Method == http.MethodPost {
		for {
			var (
				fieldErrors = parseForm(req.Form, form)
				u           = &db.User{}
			)
			if len(fieldErrors) != 0 {
				ctx["field_errors"] = fieldErrors
				break
			}
			if err := s.conn.
				Where("username = ?", form.Username).
				First(u).Error; err != nil {
				ctx["error"] = errInvalidCredentials
				break
			}
			if err := u.Authenticate(form.Password); err != nil {
				ctx["error"] = errInvalidCredentials
				break
			}
			session, _ := s.store.Get(req, sessionName)
			session.Values[sessionUserID] = u.ID
			session.Save(req, w)
			http.Redirect(w, req, "/", http.StatusFound)
			return
		}
	}
	s.render(w, req, "users/login.html", ctx)
}
