// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPL
// license that can be found in the LICENSE file.

package model

// SearchParams is a struct that has information about filtering ads for client.
type SearchParams struct {
	Query  string `db:"query" schema:"query,optional"`
	Limit  int    `db:"limit" schema:"limit,optional"`
	Offset int    `db:"offset" schema:"offset,optional"`
}
