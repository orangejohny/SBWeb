package sessionmanager

// Config is struct for configuring of session manager
type Config struct {
	DBAddress      string
	TockenLength   int
	ExpirationTime int
}
