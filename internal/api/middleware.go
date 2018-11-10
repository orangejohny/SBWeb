package api

import (
	"net/http"

	"github.com/orangejohny/SBWeb/internal/model"
)

// checkCookieMiddleware checks authentification of user
func checkCookieMiddleware(m *model.Model, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		cookieSession, err := r.Cookie("session_id")
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(apiErrorHandle("Can't update profile", "NoCookieError", err))
			return
		}

		_, err = m.CheckSession(&model.SessionID{
			ID: cookieSession.Value,
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(apiErrorHandle("Can't update profile", "No such session", err))
			return
		}

		next.ServeHTTP(w, r)
	})
}
