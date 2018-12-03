// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPL
// license that can be found in the LICENSE file.

package model

// SM describes interface of session manager.
type SM interface {
	CreateSession(in *Session, expires bool) (*SessionID, error)
	CheckSession(in *SessionID) (*Session, error)
	DeleteSession(in *SessionID) error

	TryReconnect() error
	IsConnected() bool
}
