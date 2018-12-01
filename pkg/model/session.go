// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPL
// license that can be found in the LICENSE file.

package model

// Session is object which represents session data
type Session struct {
	ID        int64
	Login     string
	UserAgent string
}

// SessionID is used as identificator of user's session
type SessionID struct {
	ID string
}
