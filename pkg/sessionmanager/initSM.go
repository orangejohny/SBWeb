// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPL
// license that can be found in the LICENSE file.

package sessionmanager

import (
	"log"

	"github.com/garyburd/redigo/redis"
)

// InitConnSM initiates connection to redis database.
// It returns struct that implements model.SM interface used
// by API to interact with sessions.
func InitConnSM(cfg Config) (*SessionManager, error) {
	redisConn, err := redis.DialURL(cfg.DBAddress)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	sessManager := &SessionManager{
		redisAddr:      cfg.DBAddress,
		redisConn:      redisConn,
		tockenLength:   cfg.TockenLength,
		expirationTime: cfg.ExpirationTime,
	}

	return sessManager, nil
}
