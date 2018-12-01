// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPL
// license that can be found in the LICENSE file.

/*
Package sessionmanager is used to implement model.SM interface. Its purpose
is controlling sessions of clients that are used for authentification.

Session manager uses key-value storage Redis for sessions.
*/
package sessionmanager

// Config is a struct for configuring session manager.
type Config struct {
	DBAddress      string `json:"DBAddress,"`
	TockenLength   int    `json:"TockenLength,int"`
	ExpirationTime int    `json:"ExpirationTime,int"`
}
