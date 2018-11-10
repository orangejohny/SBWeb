package db

// Config is stuct for database configuration
type Config struct {
	DBAddress    string `json:"DBAddress,"`
	MaxOpenConns int    `json:"MaxOpenConns,int"`
}
