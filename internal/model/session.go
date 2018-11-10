package model

// Session is object that represent session data
type Session struct {
	ID        int64
	Login     string
	UserAgent string
}

// SessionID is used as identificator of user session
type SessionID struct {
	ID string
}
