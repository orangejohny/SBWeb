package sessionmanager

import (
	"encoding/json"

	"bmstu.codes/developers34/SBWeb/internal/model"
	"github.com/garyburd/redigo/redis"
)

// CreateSession creates new session
func (sm *SessionManager) CreateSession(in *model.Session) (*model.SessionID, error) {
	tocken, err := generateRandomString(sm.tockenLength)
	if err != nil {
		return nil, err
	}

	id := model.SessionID{ID: tocken}
	dataSerialized, _ := json.Marshal(in)
	mkey := "sessions:" + id.ID
	_, err = redis.String(sm.redisConn.Do("SET", mkey, dataSerialized, "EX", sm.expirationTime))
	if err != nil {
		return nil, err
	}

	return &id, nil
}

// CheckSession checks if session with such ID exists
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

// DeleteSession deletes session with such ID
func (sm *SessionManager) DeleteSession(in *model.SessionID) error {
	mkey := "sessions:" + in.ID
	_, err := redis.Int(sm.redisConn.Do("DEL", mkey))
	return err
}
