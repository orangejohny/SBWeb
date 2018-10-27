package model

// Model is a struct that contains db interface. Such project modell allows
// to use different database implemetnation without changing business-logic
type Model struct {
	db
}

// New creates Model structure from object that implements db interface
func New(db db) *Model {
	return &Model{
		db: db,
	}
}
