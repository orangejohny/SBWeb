package model

// Session is object that represent session data
type Session struct {
	Login     string
	UserAgent string
}

// SessionID is used as identificator of user session
type SessionID struct {
	ID string
}
