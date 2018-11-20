package model

// sm describes interface of session manager
type sm interface {
	CreateSession(in *Session, expires bool) (*SessionID, error)
	CheckSession(in *SessionID) (*Session, error)
	DeleteSession(in *SessionID) error
}
