package model

// SearchParams is a struct that has information about filtering ads for client
type SearchParams struct {
	Query  string `db:"query" schema:"query,optional"`
	Limit  int    `db:"limit" schema:"limit,optional"`
	Offset int    `db:"offset" schema:"offset,optional"`
}
