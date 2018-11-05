package model

// Model is a struct that contains db and sm interfaces. Such project model allows
// to use different database and session manager implementation without changing business-logic
type Model struct {
	db
	sm
}

// New creates Model structure from object that implements db interface
func New(db db, sm sm) *Model {
	return &Model{
		db: db,
		sm: sm,
	}
}
