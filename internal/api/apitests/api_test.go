package apitests

import (
	"errors"

	"bmstu.codes/developers34/SBWeb/internal/model"
)

type smTestCase struct {
	session   model.Session
	sessionID model.SessionID
	err       error
}

func (sm *smTestCase) CreateSession(in *model.Session, expires bool) (*model.SessionID, error) {
	return &sm.sessionID, sm.err
}

func (sm *smTestCase) CheckSession(in *model.SessionID) (*model.Session, error) {
	return &sm.session, sm.err
}

func (sm *smTestCase) DeleteSession(in *model.SessionID) error {
	return sm.err
}

type dbTestCase struct {
	ad   model.AdItem
	ads  []*model.AdItem
	user model.User
	id   int64
	err  error
}

func (db *dbTestCase) GetAd(adID int64) (*model.AdItem, error) {
	if adID != db.ad.ID {
		return &db.ad, errors.New("No ad with such id")
	}
	return &db.ad, db.err
}

func (db *dbTestCase) GetAds(limit int, offset int) ([]*model.AdItem, error) {
	return db.ads, db.err
}

func (db *dbTestCase) GetAdsOfUser(userID int64) ([]*model.AdItem, error) {
	return db.ads, db.err
}

func (db *dbTestCase) GetUserWithID(userID int64) (*model.User, error) {
	return &db.user, db.err
}

func (db *dbTestCase) GetUserWithEmail(email string) (*model.User, error) {
	return &db.user, db.err
}

func (db *dbTestCase) NewUser(user *model.User) (int64, error) {
	return db.id, db.err
}

func (db *dbTestCase) NewAd(ad *model.AdItem) (int64, error) {
	return db.id, db.err
}

func (db *dbTestCase) EditUser(user *model.User) (int64, error) {
	return db.id, db.err
}

func (db *dbTestCase) EditAd(ad *model.AdItem) (int64, error) {
	return db.id, db.err
}

func (db *dbTestCase) RemoveUser(userID int64) (int64, error) {
	return db.id, db.err
}

func (db *dbTestCase) RemoveAd(adID int64) (int64, error) {
	return db.id, db.err
}

type testCase struct {
	smTestCase
	dbTestCase
}
