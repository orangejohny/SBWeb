package sessionmanager

import (
	"log"

	"github.com/garyburd/redigo/redis"
)

// InitConnSM initiates connection to sessions database
func InitConnSM(cfg Config) (*SessionManager, error) {
	redisConn, err := redis.DialURL(cfg.DBAddress)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	sessManager := &SessionManager{
		redisConn:      redisConn,
		tockenLength:   cfg.TockenLength,
		expirationTime: cfg.ExpirationTime,
	}

	return sessManager, nil
}
