package model

// Model is a struct that contains db and sm interfaces. Such project model allows
// to use different database and session manager implementation without changing business-logic
type Model struct {
	DB
	SM
}

// New creates Model structure from object that implements db interface
func New(db DB, sm SM) *Model {
	return &Model{
		DB: db,
		SM: sm,
	}
}
