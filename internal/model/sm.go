package model

// SM describes interface of session manager
type SM interface {
	CreateSession(in *Session, expires bool) (*SessionID, error)
	CheckSession(in *SessionID) (*Session, error)
	DeleteSession(in *SessionID) error
}
