package server

import (
	"context"
	"net/http"
	"net/url"

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
func (s *Server) loadUser(r *http.Request) *http.Request {
	session, _ := s.store.Get(r, sessionName)
	if v, _ := session.Values[sessionUserID]; v != nil {
		u := &db.User{}
		if err := s.conn.First(u, v).Error; err == nil {
			r = r.WithContext(context.WithValue(r.Context(), contextUser, u))
		}
	}
	return r
}

// requireUser prevents a visitor from accessing a page before logging in.
func (s *Server) requireUser(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(contextUser) == nil {
			u := url.QueryEscape(r.URL.RequestURI())
			http.Redirect(w, r, "/login?url="+u, http.StatusFound)
		} else {
			fn(w, r)
		}
	}
}

// requireAdmin prevents ordinary users from accessing an admin page.
func (s *Server) requireAdmin(fn http.HandlerFunc) http.HandlerFunc {
	return s.requireUser(func(w http.ResponseWriter, r *http.Request) {
		if !r.Context().Value(contextUser).(*db.User).IsAdmin {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		} else {
			fn(w, r)
		}
	})
}

type loginForm struct {
	Username string `form:"required"`
	Password string `form:"required"`
}

// login displays the login form and verifies the supplied credentials upon
// submission.
func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var (
		form = &loginForm{}
		ctx  = pongo2.Context{
			"title": "Login",
			"form":  form,
		}
	)
	if r.Method == http.MethodPost {
		for {
			var (
				fieldErrors = parseForm(r.Form, form)
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
			session, _ := s.store.Get(r, sessionName)
			session.Values[sessionUserID] = u.ID
			session.Save(r, w)
			redir := r.URL.Query().Get("url")
			if redir == "" {
				redir = "/"
			}
			http.Redirect(w, r, redir, http.StatusFound)
			return
		}
	}
	s.render(w, r, "users/login.html", ctx)
}

// logout logs out the current user.
func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, sessionName)
	session.Values[sessionUserID] = nil
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}
