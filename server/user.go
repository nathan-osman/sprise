package server

import (
	"context"
	"net/http"

	"github.com/nathan-osman/sprise/db"
)

const (
	sessionName   = "session"
	sessionUserID = "userID"

	contextUser = "user"
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
