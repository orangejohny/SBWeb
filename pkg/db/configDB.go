// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPL
// license that can be found in the LICENSE file.

/*
Package db implements model.db interface. It contains everything that needed by
API to interact with database.
*/
package db

// Config is stuct for database configuration.
type Config struct {
	DBAddress    string `json:"DBAddress,"`
	MaxOpenConns int    `json:"MaxOpenConns,int"`
}
