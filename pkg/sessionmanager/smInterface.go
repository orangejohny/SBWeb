// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPL
// license that can be found in the LICENSE file.

package sessionmanager

import (
	"encoding/json"
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/orangejohny/SBWeb/pkg/model"
)

// CreateSession creates new session in database.
func (sm *SessionManager) CreateSession(in *model.Session, expires bool) (*model.SessionID, error) {
	tocken, err := generateRandomString(sm.tockenLength)
	if err != nil {
		return nil, err
	}

	id := model.SessionID{ID: tocken}
	dataSerialized, _ := json.Marshal(in)
	mkey := "sessions:" + id.ID
	if expires {
		_, err = redis.String(sm.redisConn.Do("SET", mkey, dataSerialized, "EX", sm.expirationTime))
	} else {
		_, err = redis.String(sm.redisConn.Do("SET", mkey, dataSerialized))
	}
	if err != nil {
		return nil, err
	}

	return &id, nil
}

// CheckSession checks if session with such ID exists in database.
func (sm *SessionManager) CheckSession(in *model.SessionID) (*model.Session, error) {
	mkey := "sessions:" + in.ID
	data, err := redis.Bytes(sm.redisConn.Do("GET", mkey))
	if err != nil {
		return nil, err
	}

	sess := &model.Session{}
	err = json.Unmarshal(data, sess)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

// DeleteSession deletes session with such ID.
func (sm *SessionManager) DeleteSession(in *model.SessionID) error {
	mkey := "sessions:" + in.ID
	_, err := redis.Int(sm.redisConn.Do("DEL", mkey))
	return err
}

// TryReconnect reconnects to redis
func (sm *SessionManager) TryReconnect() error {
	conn, err := redis.DialURL(sm.redisAddr)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	sm.redisConn = conn
	return nil
}

// IsConnected checks if connection is active
func (sm *SessionManager) IsConnected() bool {
	_, err := sm.redisConn.Do("PING")
	if err != nil {
		return false
	}
	return true
}
