package sessionmanager

// Config is struct for configuring of session manager
type Config struct {
	DBAddress      string `json:"DBAddress,"`
	TockenLength   int    `json:"TockenLength,int"`
	ExpirationTime int    `json:"ExpirationTime,int"`
}
