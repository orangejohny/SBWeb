// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPL
// license that can be found in the LICENSE file.

package sessionmanager

import "github.com/garyburd/redigo/redis"

// SessionManager stores connection to redis database.
type SessionManager struct {
	redisConn      redis.Conn
	tockenLength   int
	expirationTime int
}
