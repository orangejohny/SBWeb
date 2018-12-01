// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPL
// license that can be found in the LICENSE file.

package api

// Config for api package. Address is a host with port (i.e. http://127.0.0.1:8080).
type Config struct {
	Address      string `json:"Address,"`
	ReadTimeout  string `json:"ReadTimeout,"`
	WriteTimeout string `json:"WriteTimeout,"`
	IdleTimeout  string `json:"IdleTimeout,"`
}
