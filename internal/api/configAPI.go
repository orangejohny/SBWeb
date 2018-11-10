package api

// Config for api package. Address is domain name of server
type Config struct {
	Address      string `json:"Address,"`
	ReadTimeout  string `json:"ReadTimeout,"`
	WriteTimeout string `json:"WriteTimeout,"`
	IdleTimeout  string `json:"IdleTimeout,"`
}
