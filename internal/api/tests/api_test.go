package api_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"bmstu.codes/developers34/SBWeb/internal/api"

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
	if adID != 15 {
		fmt.Println(adID)
		return nil, errors.New("No ad with such id")
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

func TestApi(t *testing.T) {
	testDB := dbTestCase{
		ad: model.AdItem{
			ID:    15,
			Title: "Building",
			City:  "Moscow",
			User: model.User{
				ID:        12,
				FirstName: "Ivan",
				LastName:  "Ivanov",
				Email:     "Ivan@ivanov.iva",
			},
			Description: "Awesome",
		},
		err: nil,
	}

	testSM := smTestCase{
		sessionID: model.SessionID{ID: "dwedqf1234ewdw"},
		err:       nil,
	}

	m := model.New(&testDB, &testSM)

	go api.StartServer(api.Config{
		Address:      "localhost:54000",
		ReadTimeout:  "10s",
		WriteTimeout: "10s",
		IdleTimeout:  "10s",
	}, m)

	client := http.Client{
		Timeout: time.Second * 5,
	}
	result, err := client.Get("http://localhost:54000/ads/15")

	if err != nil {
		t.Fatal("Expected no error while request\nGot: ", err.Error())
	}

	fmt.Println("eeeeeeeeee")

	adData, _ := json.Marshal(model.AdItem{
		ID:    15,
		Title: "Building",
		City:  "Moscow",
		User: model.User{
			ID:        12,
			FirstName: "Ivan",
			LastName:  "Ivanov",
			Email:     "Ivan@ivanov.iva",
		},
		Description: "Awesome"})

	defer result.Body.Close()
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatal("Expected no error while reading body")
	}

	if string(body) != string(adData) {
		t.Error("Expected equal ads:\nExpected:\n", string(adData), "Received:\n", string(body))
	}
}
