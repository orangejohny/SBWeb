package sessionmanager

import "github.com/garyburd/redigo/redis"

// SessionManager stores connection to redis database
type SessionManager struct {
	redisConn      redis.Conn
	tockenLength   int
	expirationTime int
}
