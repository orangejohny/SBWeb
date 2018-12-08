// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPL
// license that can be found in the LICENSE file.

// middleware.go contains different middlewares which can be used with API server.

package api

import (
	"log"
	"net/http"
	"net/http/httputil"

	"bmstu.codes/developers34/SBWeb/pkg/model"
)

// TODO: add middleware that checks connection to DB and SM

// checkCookieMiddleware checks authentification of user.
func checkCookieMiddleware(m *model.Model, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		cookieSession, err := r.Cookie("session_id")
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(apiErrorHandle(requiredCookie, noCookieError, err, noCookieMsg))
			return
		}

		_, err = m.CheckSession(&model.SessionID{
			ID: cookieSession.Value,
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(apiErrorHandle(badCookie, badCookieErr, err, badCookieMsg))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// logRequestMiddleware logs incoming request for debugging.
func logRequestMiddleware(m *model.Model, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		byts, _ := httputil.DumpRequest(r, true)
		log.Println(string(byts))

		next.ServeHTTP(w, r)
	})
}

// checkConnSM checks connection to SM and trying to reconnect if needed
func checkConnSM(m *model.Model, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		if !m.IsConnected() {
			err = m.TryReconnect()
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(apiErrorHandle(connectProvider, "ConnSMErr", err, "Can't connect with SM"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
