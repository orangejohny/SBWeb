package api

import (
	"net/http"

	"github.com/orangejohny/SBWeb/internal/model"
)

// getIDfromCookie return ID of user with cookie from request
// this function must be used with checkSessionMiddleware because
// it doesn't handle any errors
func getIDfromCookie(m *model.Model, r *http.Request) int64 {
	cookieSession, _ := r.Cookie("session_id")
	session, _ := m.CheckSession(&model.SessionID{
		ID: cookieSession.Value,
	})

	return session.ID
}
