// Copyright 2018 Dmitry Kargashin <dkargashin3@gmail.com>
// Use of this source code is governed by GNU LGPL
// license that can be found in the LICENSE file.

// file api_test.go contains unit tests for api package.

package api_test

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"bmstu.codes/developers34/SBWeb/pkg/api"
	"bmstu.codes/developers34/SBWeb/pkg/api/mock_model"

	"github.com/golang/mock/gomock"
	jsoniter "github.com/json-iterator/go"

	"bmstu.codes/developers34/SBWeb/pkg/model"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	domain = "http://localhost:49123"
)

func TestStartWithBadConfig(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	_, ch := api.StartServer(api.Config{
		Address:      "localhost:54000",
		ReadTimeout:  "10eregeurgi",
		WriteTimeout: "10s",
		IdleTimeout:  "10s",
	}, nil)
	err := <-ch
	if err == http.ErrServerClosed {
		t.Error("Expected other error")
	}
}

var adsInDB = map[int64]*model.AdItem{
	15: {
		ID:          0,
		Title:       "Building",
		City:        "Moscow",
		UserID:      12,
		Description: "Awesome",
	},
	16: {
		ID:     15,
		Title:  "Building",
		City:   "Moscow",
		UserID: 12,
		User: model.User{
			ID:        12,
			FirstName: "Ivan",
			LastName:  "Ivanov",
			Email:     "Ivan@ivanov.com",
		},
		Description: "Awesome",
	},
	17: {
		ID:          15,
		Title:       "Building",
		City:        "Moscow",
		Description: "Awesome",
		User: model.User{
			ID: 12,
		},
	},
	21: {
		ID:    21,
		Title: "Broking",
		City:  "Vegas",
		User: model.User{
			ID:        18,
			FirstName: "John",
			LastName:  "Johnov",
			Email:     "john@johnov.com",
		},
		Description: "can broke everything",
	},
	54: {
		ID:    54,
		Title: "House",
		City:  "New York",
		User: model.User{
			ID:        17,
			FirstName: "Petya",
			LastName:  "Succeeded",
			Email:     "pet@animal.com",
		},
		Description: "very good house",
	},
}

var usersInDB = map[int64]*model.User{
	12: {
		ID:        12,
		FirstName: "Ivan",
		LastName:  "Ivanov",
		Email:     "Ivan@ivanov.com",
		Password:  "$2a$10$kCw55AyZuzo3u6GicdIjg.RKscZZ0IxvZMSKIiuY.fQ7R2F9OPbba",
	},
	17: {
		ID:        17,
		FirstName: "Ivan",
		LastName:  "Ivanov",
		Email:     "pet@animal.com",
		Password:  "$2a$10$kCw55AyZuzo3u6GicdIjg.RKscZZ0IxvZMSKIiuY.fQ7R2F9OPbba",
	},
	18: {
		ID:        18,
		FirstName: "John",
		LastName:  "Johnov",
		Email:     "john@johnov.com",
	},
}

// apiError is a struct that returned by
// API in case of error.
type apiError struct {
	Description string `json:"description"`
	Message     string `json:"message"`
	ErrorCode   string `json:"error"`
}

// createUserResp is a struct that returned
// by API in case of succeed creation.
type createUserResp struct {
	ID  int64
	Ref string
}

// loginUserResp is a struct that returned to
// android app in case of succeed login.
type loginUserResp struct {
	ID        int64  `json:"id,"`
	FirstName string `json:"first_name,"`
	LastName  string `json:"last_name,"`
}

// mockDataDB is a struct that contains information to
// mock functions of interface model.DB.
type mockDataDB struct {
	inputID     int64
	inputLimit  int
	inputOffset int
	inputQuery  string
	inputEmail  string
	inputUser   *model.User
	inputAd     *model.AdItem

	inputIDimg int64

	outputErrorImg error
	outputUserImg  *model.User

	outputUser       *model.User
	outputUserCreate *createUserResp
	outputAd         *model.AdItem
	outputAds        []*model.AdItem
	outputError      error
	outputID         int64

	outputErrorSec error
}

// mockDataSM is a struct that contains information to
// mock functions from interface model.SM.
type mockDataSM struct {
	inputSession   *model.Session
	inputSessionID *model.SessionID
	inputExpires   bool

	outputSession   *model.Session
	outputSessionID *model.SessionID
	outputError     error
}

// testCase represents one particular test case.
type testCase struct {
	// flags for db
	isGetAd            bool
	isGetAds           bool
	isGetAdsOfUser     bool
	isGetUserWithID    bool
	isNewUser          bool
	isGetUserWithEmail bool
	isEditUser         bool
	isRemoveUser       bool
	isGetUserWithIDImg bool
	isNewAd            bool
	isEditAd           bool
	isRemoveAd         bool

	// flags of actions
	isLoginPage bool

	// flags for sm
	isCreateSession      bool
	isCheckSession       bool
	isSecondCheckSession bool
	isDeleteSession      bool
	isPrepareCheckConnSM bool

	// expect flags
	isPrepareDB bool
	isPrepareSM bool

	// request of client
	request *http.Request

	// data for functions of DB and SM
	sm *mockDataSM
	db *mockDataDB

	// expected response of server
	expectedStatusCode  int
	expectedUser        *model.User
	expectedAd          *model.AdItem
	expectedAds         []*model.AdItem
	expectedUserCreate  *createUserResp
	expectedLogin       *loginUserResp
	expectedCookieValue string
	expectedAdCreate    *createUserResp
}

var testCases = []testCase{
	{
		isGetAd:     true,
		isPrepareDB: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("GET", domain+"/ads/15", nil)
			return r
		}(),
		db: &mockDataDB{
			inputID:     15,
			outputAd:    adsInDB[15],
			outputError: nil,
		},
		expectedAd:         adsInDB[15],
		expectedStatusCode: 200,
	},
	{
		isGetAd:     true,
		isPrepareDB: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("GET", domain+"/ads/15", nil)
			return r
		}(),
		db: &mockDataDB{
			inputID:     15,
			outputAd:    &model.AdItem{ID: -1},
			outputError: errors.New("e"),
		},
		expectedStatusCode: 400,
	},
	{
		isGetAd:     true,
		isPrepareDB: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("GET", domain+"/ads/15", nil)
			return r
		}(),
		db: &mockDataDB{
			inputID:     15,
			outputAd:    &model.AdItem{},
			outputError: errors.New("e"),
		},
		expectedStatusCode: 500,
	},
	{
		isGetAds:    true,
		isPrepareDB: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("GET", domain+"/ads?offset=1&limit=12", nil)
			return r
		}(),
		db: &mockDataDB{
			inputLimit:  12,
			inputOffset: 1,
			outputAds:   []*model.AdItem{adsInDB[15], adsInDB[21]},
			outputError: nil,
		},
		expectedAds:        []*model.AdItem{adsInDB[15], adsInDB[21]},
		expectedStatusCode: 200,
	},
	{
		isGetAds:    true,
		isPrepareDB: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("GET", domain+"/ads?offset=-1&limit=12", nil)
			return r
		}(),
		db: &mockDataDB{
			inputLimit:  12,
			inputOffset: 0,
			outputAds:   []*model.AdItem{adsInDB[15], adsInDB[21]},
			outputError: nil,
		},
		expectedAds:        []*model.AdItem{adsInDB[15], adsInDB[21]},
		expectedStatusCode: 200,
	},
	{
		isGetAds:    true,
		isPrepareDB: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("GET", domain+"/ads?offset=1&limit=de12", nil)
			return r
		}(),
		db: &mockDataDB{
			inputLimit:  15,
			inputOffset: 1,
			outputAds:   []*model.AdItem{adsInDB[15], adsInDB[21]},
			outputError: nil,
		},
		expectedAds:        []*model.AdItem{adsInDB[15], adsInDB[21]},
		expectedStatusCode: 200,
	},
	{
		isGetAds:    true,
		isPrepareDB: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("GET", domain+"/ads?offset=1&limit=12", nil)
			return r
		}(),
		db: &mockDataDB{
			inputLimit:  12,
			inputOffset: 1,
			outputAds:   []*model.AdItem{},
			outputError: errors.New("e"),
		},
		expectedStatusCode: 500,
	},
	{
		isGetAdsOfUser: true,
		isPrepareDB:    true,
		request: func() *http.Request {
			r, _ := http.NewRequest("GET", domain+"/users/54?show_ads=true", nil)
			return r
		}(),
		db: &mockDataDB{
			inputID:     54,
			outputAds:   []*model.AdItem{adsInDB[15], adsInDB[21]},
			outputError: nil,
		},
		expectedAds:        []*model.AdItem{adsInDB[15], adsInDB[21]},
		expectedStatusCode: 200,
	},
	{
		isGetAdsOfUser: true,
		isPrepareDB:    true,
		request: func() *http.Request {
			r, _ := http.NewRequest("GET", domain+"/users/54?show_ads=true", nil)
			return r
		}(),
		db: &mockDataDB{
			inputID:     54,
			outputAds:   make([]*model.AdItem, 0),
			outputError: errors.New("e"),
		},
		expectedStatusCode: 500,
	},
	{
		isGetUserWithID: true,
		isPrepareDB:     true,
		request: func() *http.Request {
			r, _ := http.NewRequest("GET", domain+"/users/18", nil)
			return r
		}(),
		db: &mockDataDB{
			inputID:     18,
			outputUser:  usersInDB[18],
			outputError: nil,
		},
		expectedUser:       usersInDB[18],
		expectedStatusCode: 200,
	},
	{
		isGetUserWithID: true,
		isPrepareDB:     true,
		request: func() *http.Request {
			r, _ := http.NewRequest("GET", domain+"/users/18", nil)
			return r
		}(),
		db: &mockDataDB{
			inputID:     18,
			outputUser:  &model.User{ID: -1},
			outputError: errors.New("e"),
		},
		expectedStatusCode: 400,
	},
	{
		isGetUserWithID: true,
		isPrepareDB:     true,
		request: func() *http.Request {
			r, _ := http.NewRequest("GET", domain+"/users/18", nil)
			return r
		}(),
		db: &mockDataDB{
			inputID:     18,
			outputUser:  &model.User{},
			outputError: errors.New("e"),
		},
		expectedStatusCode: 500,
	},
	{
		isNewUser:   true,
		isPrepareDB: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/new",
				strings.NewReader("email=pet@animal.com&first_name=Ivan&last_name=Ivanov&password=123456"))
			return r
		}(),
		db: &mockDataDB{
			outputID:    17,
			outputError: nil,
		},
		expectedStatusCode: http.StatusCreated,
		expectedUserCreate: &createUserResp{
			ID:  17,
			Ref: "/users/17",
		},
	},
	{
		isNewUser: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/new",
				strings.NewReader("%email=pet@animal.com&first_name=Ivan&last_name=Ivanov&password=123456"))
			return r
		}(),
		expectedStatusCode: http.StatusBadRequest,
	},
	{
		isNewUser:   true,
		isPrepareDB: true,
		request: func() *http.Request {
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			writer.WriteField("email", "pet@animal.com")
			writer.WriteField("first_name", "Ivan")
			writer.WriteField("last_name", "Ivanov")
			writer.WriteField("password", "123456")
			writer.Close()
			r, _ := http.NewRequest("POST", domain+"/users/new", body)
			r.Header.Set("Content-Type", writer.FormDataContentType())
			return r
		}(),
		db: &mockDataDB{
			outputID:    17,
			outputError: nil,
		},
		expectedStatusCode: http.StatusCreated,
		expectedUserCreate: &createUserResp{
			ID:  17,
			Ref: "/users/17",
		},
	},
	{
		isNewUser:   true,
		isPrepareDB: true,
		request: func() *http.Request {
			path := os.Getenv("CI_PROJECT_DIR") + "docs/AuthReq.PNG"
			file, err := os.Open(path)
			if err != nil {
				log.Fatal("Can't open file")
			}
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, err := writer.CreateFormFile("images", filepath.Base(path))
			if err != nil {
				log.Fatal("Can't read file")
			}
			io.Copy(part, file)
			writer.WriteField("email", "pet@animal.com")
			writer.WriteField("first_name", "Ivan")
			writer.WriteField("last_name", "Ivanov")
			writer.WriteField("password", "123456")
			writer.Close()
			file.Close()
			r, _ := http.NewRequest("POST", domain+"/users/new", body)
			r.Header.Set("Content-Type", writer.FormDataContentType())
			return r
		}(),
		db: &mockDataDB{
			outputID:    17,
			outputError: nil,
		},
		expectedStatusCode: http.StatusCreated,
		expectedUserCreate: &createUserResp{
			ID:  17,
			Ref: "/users/17",
		},
	},
	{
		isNewUser: true,
		request: func() *http.Request {
			path := os.Getenv("CI_PROJECT_DIR") + "docs/curlTest.md"
			file, err := os.Open(path)
			if err != nil {
				log.Fatal("Can't open file")
			}
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, err := writer.CreateFormFile("images", filepath.Base(path))
			if err != nil {
				log.Fatal("Can't read file")
			}
			io.Copy(part, file)
			writer.WriteField("email", "pet@animal.com")
			writer.WriteField("first_name", "Ivan")
			writer.WriteField("last_name", "Ivanov")
			writer.WriteField("password", "123456")
			writer.Close()
			file.Close()
			r, _ := http.NewRequest("POST", domain+"/users/new", body)
			r.Header.Set("Content-Type", writer.FormDataContentType())
			return r
		}(),
		expectedStatusCode: http.StatusInternalServerError,
	},
	{
		isNewUser: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/new",
				strings.NewReader("script=bad&email=pet@animal.com&first_name=Ivan&last_name=Ivanov&password=123456"))
			return r
		}(),
		expectedStatusCode: http.StatusBadRequest,
	},
	{
		isNewUser: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/new",
				strings.NewReader("email=pet@animal.com&first_name=&last_name=Ivanov&password=123456"))
			return r
		}(),
		expectedStatusCode: http.StatusBadRequest,
	},
	{
		isNewUser: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/new",
				strings.NewReader("email=pet@animal.om&first_name=Ivan&last_name=Ivanov&password=123456"))
			return r
		}(),
		expectedStatusCode: http.StatusBadRequest,
	},
	{
		isNewUser:   true,
		isPrepareDB: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/new",
				strings.NewReader("email=pet@animal.com&first_name=Ivan&last_name=Ivanov&password=123456"))
			return r
		}(),
		db: &mockDataDB{
			outputID:    -1,
			outputError: errors.New("e"),
		},
		expectedStatusCode: http.StatusBadRequest,
	},
	{
		isNewUser:   true,
		isPrepareDB: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/new",
				strings.NewReader("email=pet@animal.com&first_name=Ivan&last_name=Ivanov&password=123456"))
			return r
		}(),
		db: &mockDataDB{
			outputID:    0,
			outputError: errors.New("e"),
		},
		expectedStatusCode: http.StatusInternalServerError,
	},
	{
		isGetUserWithEmail:   true,
		isCreateSession:      true,
		isPrepareDB:          true,
		isPrepareSM:          true,
		isPrepareCheckConnSM: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/login",
				strings.NewReader("email=pet@animal.com&password=123456"))
			return r
		}(),
		db: &mockDataDB{
			inputEmail:  "pet@animal.com",
			outputUser:  usersInDB[17],
			outputError: nil,
		},
		sm: &mockDataSM{
			inputSession: &model.Session{
				ID:        17,
				Login:     "pet@animal.com",
				UserAgent: "Go-http-client/1.1",
			},
			inputExpires:    true,
			outputSessionID: &model.SessionID{ID: "id"},
		},
		expectedStatusCode:  200,
		expectedCookieValue: "id",
	},
	{
		isGetUserWithEmail:   true,
		isCreateSession:      true,
		isPrepareDB:          true,
		isPrepareCheckConnSM: true,
		isPrepareSM:          true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/login",
				strings.NewReader("email=pet@animal.com&password=123456"))
			r.Header.Set("User-Agent", "Android_app")
			return r
		}(),
		db: &mockDataDB{
			inputEmail:  "pet@animal.com",
			outputUser:  usersInDB[17],
			outputError: nil,
		},
		sm: &mockDataSM{
			inputSession: &model.Session{
				ID:        17,
				Login:     "pet@animal.com",
				UserAgent: "Android_app",
			},
			inputExpires:    false,
			outputSessionID: &model.SessionID{ID: "id"},
		},
		expectedLogin: &loginUserResp{
			ID:        17,
			FirstName: "Ivan",
			LastName:  "Ivanov",
		},
		expectedStatusCode:  200,
		expectedCookieValue: "id",
	},
	{
		isGetUserWithEmail:   true,
		isPrepareCheckConnSM: true,
		isCreateSession:      true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/login",
				strings.NewReader("%email=pet@animal.com&password=123456"))
			return r
		}(),
		expectedStatusCode: 400,
	},
	{
		isGetUserWithEmail:   true,
		isPrepareCheckConnSM: true,
		isCreateSession:      true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/login",
				strings.NewReader("email=pet@animal.cv&password=123456"))
			return r
		}(),
		expectedStatusCode: 400,
	},
	{
		isGetUserWithEmail:   true,
		isPrepareCheckConnSM: true,
		isCreateSession:      true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/login",
				strings.NewReader("script=vvvv&email=pet@animal.com&password=123456"))
			return r
		}(),
		expectedStatusCode: 400,
	},
	{
		isPrepareCheckConnSM: true,
		isGetUserWithEmail:   true,
		isCreateSession:      true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/login",
				strings.NewReader("email=pet@animal.com&password="))
			return r
		}(),
		expectedStatusCode: 400,
	},
	{
		isGetUserWithEmail:   true,
		isCreateSession:      true,
		isPrepareCheckConnSM: true,
		isPrepareDB:          true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/login",
				strings.NewReader("email=pet@animal.com&password=123456"))
			return r
		}(),
		db: &mockDataDB{
			inputEmail:  "pet@animal.com",
			outputUser:  &model.User{ID: -1},
			outputError: errors.New("e"),
		},
		expectedStatusCode: 400,
	},
	{
		isGetUserWithEmail:   true,
		isCreateSession:      true,
		isPrepareCheckConnSM: true,
		isPrepareDB:          true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/login",
				strings.NewReader("email=pet@animal.com&password=123456"))
			return r
		}(),
		db: &mockDataDB{
			inputEmail:  "pet@animal.com",
			outputUser:  &model.User{},
			outputError: errors.New("e"),
		},
		expectedStatusCode: 500,
	},
	{
		isGetUserWithEmail:   true,
		isPrepareCheckConnSM: true,
		isCreateSession:      true,
		isPrepareDB:          true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/login",
				strings.NewReader("email=pet@animal.com&password=12456"))
			return r
		}(),
		db: &mockDataDB{
			inputEmail:  "pet@animal.com",
			outputUser:  usersInDB[17],
			outputError: nil,
		},
		expectedStatusCode: 400,
	},
	{
		isGetUserWithEmail:   true,
		isCreateSession:      true,
		isPrepareCheckConnSM: true,
		isPrepareDB:          true,
		isPrepareSM:          true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/login",
				strings.NewReader("email=pet@animal.com&password=123456"))
			return r
		}(),
		db: &mockDataDB{
			inputEmail:  "pet@animal.com",
			outputUser:  usersInDB[17],
			outputError: nil,
		},
		sm: &mockDataSM{
			inputSession: &model.Session{
				ID:        17,
				Login:     "pet@animal.com",
				UserAgent: "Go-http-client/1.1",
			},
			inputExpires:    true,
			outputSessionID: &model.SessionID{},
			outputError:     errors.New("e"),
		},
		expectedStatusCode: 500,
	},
	{
		isDeleteSession:      true,
		isPrepareCheckConnSM: true,
		isPrepareSM:          true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/logout", nil)
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputError:    nil,
		},
		expectedStatusCode: 200,
	},
	{
		isPrepareCheckConnSM: true,
		isDeleteSession:      true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/logout", nil)
			return r
		}(),
		expectedStatusCode: 200,
	},
	{
		isPrepareCheckConnSM: true,
		isDeleteSession:      true,
		isPrepareSM:          true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/logout", nil)
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputError:    errors.New("e"),
		},
		expectedStatusCode: 200,
	},
	{
		isPrepareCheckConnSM: true,
		isEditUser:           true,
		isSecondCheckSession: true,
		isCheckSession:       true,
		isPrepareDB:          true,
		isPrepareSM:          true,
		isGetUserWithIDImg:   true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/profile",
				strings.NewReader("first_name=Alex&last_name=Ivanov"))
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		db: &mockDataDB{
			inputUser: &model.User{
				ID:        17,
				FirstName: "Alex",
				LastName:  "Ivanov",
			},
			inputIDimg:     17,
			outputID:       17,
			outputUserImg:  usersInDB[17],
			outputError:    nil,
			outputErrorImg: nil,
		},
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputError:    nil,
			outputSession: &model.Session{
				ID:        17,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
		},
		expectedStatusCode: 200,
	},
	{
		isCheckSession:       true,
		isPrepareSM:          true,
		isPrepareCheckConnSM: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/profile",
				strings.NewReader("%first_name=Alex&last_name=Ivanov"))
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputError:    nil,
			outputSession: &model.Session{
				ID:        17,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
		},
		expectedStatusCode: 400,
	},
	{
		isCheckSession:       true,
		isPrepareCheckConnSM: true,
		isPrepareSM:          true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/profile",
				strings.NewReader("script=bbb&first_name=Alex&last_name=Ivanov"))
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputError:    nil,
			outputSession: &model.Session{
				ID:        17,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
		},
		expectedStatusCode: 400,
	},
	{
		isCheckSession:       true,
		isPrepareCheckConnSM: true,
		isPrepareSM:          true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/profile",
				strings.NewReader("first_name=&last_name=Ivanov"))
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputError:    nil,
			outputSession: &model.Session{
				ID:        17,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
		},
		expectedStatusCode: 400,
	},
	{
		isPrepareCheckConnSM: true,
		isCheckSession:       true,
		isPrepareSM:          true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/profile",
				strings.NewReader("first_name=efv45&last_name=Ivanov"))
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputError:    nil,
			outputSession: &model.Session{
				ID:        17,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
		},
		expectedStatusCode: 400,
	},
	{
		isCheckSession:       true,
		isPrepareCheckConnSM: true,
		isPrepareSM:          true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/profile",
				strings.NewReader("first_name=Alex&last_name=Ivanov&tel_number=dwdww"))
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputError:    nil,
			outputSession: &model.Session{
				ID:        17,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
		},
		expectedStatusCode: 400,
	},
	{
		isEditUser:           true,
		isPrepareCheckConnSM: true,
		isSecondCheckSession: true,
		isCheckSession:       true,
		isPrepareDB:          true,
		isPrepareSM:          true,
		isGetUserWithIDImg:   true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/profile",
				strings.NewReader("first_name=Alex&last_name=Ivanov"))
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		db: &mockDataDB{
			inputUser: &model.User{
				ID:        17,
				FirstName: "Alex",
				LastName:  "Ivanov",
			},
			inputIDimg:     17,
			outputErrorImg: nil,
			outputUserImg:  usersInDB[17],
			outputID:       0,
			outputError:    errors.New("e"),
		},
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputError:    nil,
			outputSession: &model.Session{
				ID:        17,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
		},
		expectedStatusCode: 500,
	},
	{
		isPrepareCheckConnSM: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/profile",
				strings.NewReader("first_name=Alex&last_name=Ivanov"))
			return r
		}(),
		expectedStatusCode: 401,
	},
	{
		isPrepareCheckConnSM: true,
		isCheckSession:       true,
		isPrepareSM:          true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/profile",
				strings.NewReader("first_name=Alex&last_name=Ivanov"))
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputError:    errors.New("e"),
			outputSession:  &model.Session{},
		},
		expectedStatusCode: 401,
	},
	{
		isPrepareCheckConnSM: true,
		isEditUser:           true,
		isSecondCheckSession: true,
		isCheckSession:       true,
		isPrepareDB:          true,
		isPrepareSM:          true,
		isGetUserWithIDImg:   true,
		request: func() *http.Request {
			path := os.Getenv("CI_PROJECT_DIR") + "docs/GeneralOverview.png"
			file, err := os.Open(path)
			if err != nil {
				log.Fatal("Can't open file")
			}
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, err := writer.CreateFormFile("images", filepath.Base(path))
			if err != nil {
				log.Fatal("Can't read file")
			}
			io.Copy(part, file)
			writer.WriteField("first_name", "Alex")
			writer.WriteField("last_name", "Ivanov")
			writer.Close()
			file.Close()
			r, _ := http.NewRequest("POST", domain+"/users/profile", body)
			r.Header.Set("Content-Type", writer.FormDataContentType())
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		db: &mockDataDB{
			inputIDimg:     17,
			outputID:       17,
			outputUserImg:  usersInDB[17],
			outputError:    nil,
			outputErrorImg: nil,
		},
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputError:    nil,
			outputSession: &model.Session{
				ID:        17,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
		},
		expectedStatusCode: 200,
	},
	{
		isPrepareCheckConnSM: true,
		isSecondCheckSession: true,
		isCheckSession:       true,
		isPrepareSM:          true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/profile",
				strings.NewReader("first_name=Alex&last_name=Ivanov&avatar_address=blbabalbal"))
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputError:    nil,
			outputSession: &model.Session{
				ID:        17,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
		},
		expectedStatusCode: 400,
	},
	{
		isPrepareCheckConnSM: true,
		isSecondCheckSession: true,
		isCheckSession:       true,
		isPrepareSM:          true,
		isGetUserWithIDImg:   true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/profile",
				strings.NewReader("first_name=Alex&last_name=Ivanov"))
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		db: &mockDataDB{
			inputIDimg:     17,
			outputUserImg:  &model.User{},
			outputErrorImg: errors.New("e"),
		},
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputError:    nil,
			outputSession: &model.Session{
				ID:        17,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
		},
		expectedStatusCode: 500,
	},
	{
		isPrepareCheckConnSM: true,
		isSecondCheckSession: true,
		isCheckSession:       true,
		isPrepareDB:          true,
		isPrepareSM:          true,
		isGetUserWithIDImg:   true,
		request: func() *http.Request {
			path := os.Getenv("CI_PROJECT_DIR") + "docs/curlTest.md"
			file, err := os.Open(path)
			if err != nil {
				log.Fatal("Can't open file")
			}
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, err := writer.CreateFormFile("images", filepath.Base(path))
			if err != nil {
				log.Fatal("Can't read file")
			}
			io.Copy(part, file)
			writer.WriteField("first_name", "Alex")
			writer.WriteField("last_name", "Ivanov")
			writer.Close()
			file.Close()
			r, _ := http.NewRequest("POST", domain+"/users/profile", body)
			r.Header.Set("Content-Type", writer.FormDataContentType())
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		db: &mockDataDB{
			inputIDimg:     17,
			outputUserImg:  usersInDB[17],
			outputError:    nil,
			outputErrorImg: nil,
		},
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputError:    nil,
			outputSession: &model.Session{
				ID:        17,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
		},
		expectedStatusCode: 500,
	},
	{
		isGetUserWithID:      true,
		isPrepareDB:          true,
		isPrepareSM:          true,
		isPrepareCheckConnSM: true,
		isCheckSession:       true,
		isSecondCheckSession: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("GET", domain+"/users/profile", nil)
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		db: &mockDataDB{
			inputID:     17,
			outputUser:  usersInDB[17],
			outputError: nil,
		},
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputSession: &model.Session{
				ID:        17,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
			outputError: nil,
		},
		expectedUser:       usersInDB[17],
		expectedStatusCode: 200,
	},
	{
		isGetUserWithID:      true,
		isPrepareDB:          true,
		isPrepareSM:          true,
		isPrepareCheckConnSM: true,
		isCheckSession:       true,
		isSecondCheckSession: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("GET", domain+"/users/profile", nil)
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		db: &mockDataDB{
			inputID:     17,
			outputUser:  &model.User{},
			outputError: errors.New("e"),
		},
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputSession: &model.Session{
				ID:        17,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
			outputError: nil,
		},
		expectedStatusCode: 500,
	},
	{
		isDeleteSession:      true,
		isRemoveUser:         true,
		isPrepareDB:          true,
		isPrepareCheckConnSM: true,
		isPrepareSM:          true,
		isCheckSession:       true,
		isSecondCheckSession: true,
		isGetUserWithIDImg:   true,
		request: func() *http.Request {
			r, _ := http.NewRequest("DELETE", domain+"/users/profile", nil)
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		db: &mockDataDB{
			inputID:        17,
			outputID:       17,
			outputError:    nil,
			inputIDimg:     17,
			outputUserImg:  usersInDB[17],
			outputErrorImg: nil,
		},
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputSession: &model.Session{
				ID:        17,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
			outputError: nil,
		},
		expectedStatusCode: 200,
	},
	{
		isRemoveUser:         true,
		isPrepareDB:          true,
		isPrepareSM:          true,
		isCheckSession:       true,
		isPrepareCheckConnSM: true,
		isSecondCheckSession: true,
		isGetUserWithIDImg:   true,
		request: func() *http.Request {
			r, _ := http.NewRequest("DELETE", domain+"/users/profile", nil)
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		db: &mockDataDB{
			inputID:        17,
			outputID:       0,
			outputError:    errors.New("e"),
			inputIDimg:     17,
			outputUserImg:  usersInDB[17],
			outputErrorImg: nil,
		},
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputSession: &model.Session{
				ID:        17,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
			outputError: nil,
		},
		expectedStatusCode: 500,
	},
	{
		isPrepareDB:          true,
		isPrepareSM:          true,
		isCheckSession:       true,
		isPrepareCheckConnSM: true,
		isSecondCheckSession: true,
		isGetUserWithID:      true,
		request: func() *http.Request {
			r, _ := http.NewRequest("DELETE", domain+"/users/profile", nil)
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		db: &mockDataDB{
			inputID:     17,
			outputID:    0,
			outputError: errors.New("e"),
		},
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputSession: &model.Session{
				ID:        17,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
			outputError: nil,
		},
		expectedStatusCode: 500,
	},
	{
		isPrepareDB:          true,
		isPrepareSM:          true,
		isNewAd:              true,
		isPrepareCheckConnSM: true,
		isCheckSession:       true,
		isSecondCheckSession: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/ads/new",
				strings.NewReader("title=Building&city=Moscow&description_ad=Awesome"))
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		db: &mockDataDB{
			inputAd:     adsInDB[15],
			outputID:    15,
			outputError: nil,
		},
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputSession: &model.Session{
				ID:        12,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
			outputError: nil,
		},
		expectedStatusCode: 201,
		expectedAdCreate:   &createUserResp{ID: 15, Ref: "/ads/15"},
	},
	{
		isPrepareSM:          true,
		isPrepareCheckConnSM: true,
		isCheckSession:       true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/ads/new",
				strings.NewReader("t%itle=Building&city=Moscow&description_ad=Awesome"))
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputSession: &model.Session{
				ID:        12,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
			outputError: nil,
		},
		expectedStatusCode: 400,
	},
	{
		isPrepareSM:          true,
		isPrepareCheckConnSM: true,
		isCheckSession:       true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/ads/new",
				strings.NewReader("title=Building&city=Moscow&description_ad=Awesome&script=bad"))
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputSession: &model.Session{
				ID:        12,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
			outputError: nil,
		},
		expectedStatusCode: 400,
	},
	{
		isPrepareSM:          true,
		isPrepareCheckConnSM: true,
		isCheckSession:       true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/ads/new",
				strings.NewReader("title=Building&city=Moscow"))
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputSession: &model.Session{
				ID:        12,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
			outputError: nil,
		},
		expectedStatusCode: 400,
	},
	{
		isPrepareDB:          true,
		isPrepareSM:          true,
		isNewAd:              true,
		isPrepareCheckConnSM: true,
		isCheckSession:       true,
		isSecondCheckSession: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/ads/new",
				strings.NewReader("title=Building&city=Moscow&description_ad=Awesome"))
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		db: &mockDataDB{
			inputAd:     adsInDB[15],
			outputID:    0,
			outputError: errors.New("e"),
		},
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputSession: &model.Session{
				ID:        12,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
			outputError: nil,
		},
		expectedStatusCode: 500,
	},
	{
		isPrepareDB:          true,
		isPrepareSM:          true,
		isNewAd:              true,
		isPrepareCheckConnSM: true,
		isCheckSession:       true,
		isSecondCheckSession: true,
		request: func() *http.Request {
			path := os.Getenv("CI_PROJECT_DIR") + "docs/AuthReq.PNG"
			file, err := os.Open(path)
			if err != nil {
				log.Fatal("Can't open file")
			}
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, err := writer.CreateFormFile("images", filepath.Base(path))
			if err != nil {
				log.Fatal("Can't read file")
			}
			io.Copy(part, file)
			writer.WriteField("title", "Building")
			writer.WriteField("city", "Moscow")
			writer.WriteField("description_ad", "Awesome")
			writer.Close()
			file.Close()
			r, _ := http.NewRequest("POST", domain+"/ads/new", body)
			r.Header.Set("Content-Type", writer.FormDataContentType())
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		db: &mockDataDB{
			inputAd:     adsInDB[15],
			outputID:    15,
			outputError: nil,
		},
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputSession: &model.Session{
				ID:        12,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
			outputError: nil,
		},
		expectedStatusCode: 201,
		expectedAdCreate:   &createUserResp{ID: 15, Ref: "/ads/15"},
	},
	{
		isPrepareSM:          true,
		isPrepareCheckConnSM: true,
		isCheckSession:       true,
		request: func() *http.Request {
			path := os.Getenv("CI_PROJECT_DIR") + "docs/curlTest.md"
			file, err := os.Open(path)
			if err != nil {
				log.Fatal("Can't open file")
			}
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, err := writer.CreateFormFile("images", filepath.Base(path))
			if err != nil {
				log.Fatal("Can't read file")
			}
			io.Copy(part, file)
			writer.WriteField("title", "Building")
			writer.WriteField("city", "Moscow")
			writer.WriteField("description_ad", "Awesome")
			writer.Close()
			file.Close()
			r, _ := http.NewRequest("POST", domain+"/ads/new", body)
			r.Header.Set("Content-Type", writer.FormDataContentType())
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputSession: &model.Session{
				ID:        12,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
			outputError: nil,
		},
		expectedStatusCode: 500,
	},
	{
		isPrepareDB:          true,
		isPrepareSM:          true,
		isEditAd:             true,
		isGetAd:              true,
		isPrepareCheckConnSM: true,
		isCheckSession:       true,
		isSecondCheckSession: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/ads/edit/15",
				strings.NewReader("title=Building&city=Moscow&description_ad=Awesome"))
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		db: &mockDataDB{
			inputID:     15,
			inputAd:     adsInDB[17],
			outputID:    15,
			outputAd:    adsInDB[16],
			outputError: nil,
		},
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputSession: &model.Session{
				ID:        12,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
			outputError: nil,
		},
		expectedStatusCode: 200,
	},
	{
		isPrepareSM:          true,
		isPrepareCheckConnSM: true,
		isCheckSession:       true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/ads/edit/15",
				strings.NewReader("ti%tle=Building&city=Moscow&description_ad=Awesome"))
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputSession: &model.Session{
				ID:        12,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
			outputError: nil,
		},
		expectedStatusCode: 400,
	},
	{
		isPrepareSM:          true,
		isPrepareCheckConnSM: true,
		isCheckSession:       true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/ads/edit/15",
				strings.NewReader("script=;DROP TABLE;&title=Building&city=Moscow&description_ad=Awesome"))
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputSession: &model.Session{
				ID:        12,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
			outputError: nil,
		},
		expectedStatusCode: 400,
	},
	{
		isPrepareSM:          true,
		isPrepareCheckConnSM: true,
		isCheckSession:       true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/ads/edit/15",
				strings.NewReader("city=Moscow&description_ad=Awesome"))
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputSession: &model.Session{
				ID:        12,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
			outputError: nil,
		},
		expectedStatusCode: 400,
	},
	{
		isPrepareDB:          true,
		isPrepareSM:          true,
		isGetAd:              true,
		isPrepareCheckConnSM: true,
		isCheckSession:       true,
		isSecondCheckSession: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/ads/edit/15",
				strings.NewReader("title=Building&city=Moscow&description_ad=Awesome"))
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		db: &mockDataDB{
			inputID:     15,
			outputAd:    &model.AdItem{ID: -1},
			outputError: errors.New("e"),
		},
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputSession: &model.Session{
				ID:        12,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
			outputError: nil,
		},
		expectedStatusCode: 400,
	},
	{
		isPrepareDB:          true,
		isPrepareSM:          true,
		isGetAd:              true,
		isPrepareCheckConnSM: true,
		isCheckSession:       true,
		isSecondCheckSession: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/ads/edit/15",
				strings.NewReader("title=Building&city=Moscow&description_ad=Awesome"))
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		db: &mockDataDB{
			inputID:     15,
			outputAd:    &model.AdItem{},
			outputError: errors.New("e"),
		},
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputSession: &model.Session{
				ID:        12,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
			outputError: nil,
		},
		expectedStatusCode: 500,
	},
	{
		isPrepareDB:          true,
		isPrepareSM:          true,
		isGetAd:              true,
		isPrepareCheckConnSM: true,
		isCheckSession:       true,
		isSecondCheckSession: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/ads/edit/15",
				strings.NewReader("title=Building&city=Moscow&description_ad=Awesome"))
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		db: &mockDataDB{
			inputID:     15,
			outputAd:    adsInDB[21],
			outputError: nil,
		},
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputSession: &model.Session{
				ID:        12,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
			outputError: nil,
		},
		expectedStatusCode: 403,
	},
	{
		isPrepareDB:          true,
		isPrepareSM:          true,
		isEditAd:             true,
		isGetAd:              true,
		isPrepareCheckConnSM: true,
		isCheckSession:       true,
		isSecondCheckSession: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/ads/edit/15",
				strings.NewReader("title=Building&city=Moscow&description_ad=Awesome"))
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		db: &mockDataDB{
			inputID:        15,
			inputAd:        adsInDB[17],
			outputID:       15,
			outputAd:       adsInDB[16],
			outputError:    nil,
			outputErrorSec: errors.New("e"),
		},
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputSession: &model.Session{
				ID:        12,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
			outputError: nil,
		},
		expectedStatusCode: 500,
	},
	{
		isPrepareDB:          true,
		isPrepareSM:          true,
		isEditAd:             true,
		isGetAd:              true,
		isPrepareCheckConnSM: true,
		isCheckSession:       true,
		isSecondCheckSession: true,
		request: func() *http.Request {
			path := os.Getenv("CI_PROJECT_DIR") + "docs/AuthReq.PNG"
			file, err := os.Open(path)
			if err != nil {
				log.Fatal("Can't open file")
			}
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, err := writer.CreateFormFile("images", filepath.Base(path))
			if err != nil {
				log.Fatal("Can't read file")
			}
			io.Copy(part, file)
			writer.WriteField("title", "Building")
			writer.WriteField("city", "Moscow")
			writer.WriteField("description_ad", "Awesome")
			writer.Close()
			file.Close()
			r, _ := http.NewRequest("POST", domain+"/ads/edit/15", body)
			r.Header.Set("Content-Type", writer.FormDataContentType())
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		db: &mockDataDB{
			inputID:     15,
			inputAd:     adsInDB[17],
			outputID:    15,
			outputAd:    adsInDB[16],
			outputError: nil,
		},
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputSession: &model.Session{
				ID:        12,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
			outputError: nil,
		},
		expectedStatusCode: 200,
	},
	{
		isPrepareDB:          true,
		isPrepareSM:          true,
		isGetAd:              true,
		isPrepareCheckConnSM: true,
		isCheckSession:       true,
		isSecondCheckSession: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/ads/edit/15",
				strings.NewReader("title=Building&city=Moscow&description_ad=Awesome&ad_images=/eee"))
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		db: &mockDataDB{
			inputID:     15,
			outputAd:    adsInDB[16],
			outputError: nil,
		},
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputSession: &model.Session{
				ID:        12,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
			outputError: nil,
		},
		expectedStatusCode: 400,
	},
	{
		isPrepareDB:          true,
		isPrepareSM:          true,
		isGetAd:              true,
		isPrepareCheckConnSM: true,
		isCheckSession:       true,
		isSecondCheckSession: true,
		request: func() *http.Request {
			path := os.Getenv("CI_PROJECT_DIR") + "docs/curlTest.md"
			file, err := os.Open(path)
			if err != nil {
				log.Fatal("Can't open file")
			}
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, err := writer.CreateFormFile("images", filepath.Base(path))
			if err != nil {
				log.Fatal("Can't read file")
			}
			io.Copy(part, file)
			writer.WriteField("title", "Building")
			writer.WriteField("city", "Moscow")
			writer.WriteField("description_ad", "Awesome")
			writer.Close()
			file.Close()
			r, _ := http.NewRequest("POST", domain+"/ads/edit/15", body)
			r.Header.Set("Content-Type", writer.FormDataContentType())
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		db: &mockDataDB{
			inputID:     15,
			inputAd:     adsInDB[17],
			outputID:    15,
			outputAd:    adsInDB[16],
			outputError: nil,
		},
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputSession: &model.Session{
				ID:        12,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
			outputError: nil,
		},
		expectedStatusCode: 500,
	},
	{
		isPrepareDB:          true,
		isPrepareSM:          true,
		isGetAd:              true,
		isPrepareCheckConnSM: true,
		isCheckSession:       true,
		isSecondCheckSession: true,
		isRemoveAd:           true,
		request: func() *http.Request {
			r, _ := http.NewRequest("DELETE", domain+"/ads/delete/15", nil)
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		db: &mockDataDB{
			inputID:     15,
			inputAd:     adsInDB[17],
			outputID:    15,
			outputAd:    adsInDB[16],
			outputError: nil,
		},
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputSession: &model.Session{
				ID:        12,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
			outputError: nil,
		},
		expectedStatusCode: 200,
	},
	{
		isPrepareDB:          true,
		isPrepareSM:          true,
		isGetAd:              true,
		isPrepareCheckConnSM: true,
		isCheckSession:       true,
		isSecondCheckSession: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("DELETE", domain+"/ads/delete/15", nil)
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		db: &mockDataDB{
			inputID:     15,
			outputAd:    &model.AdItem{ID: -1},
			outputError: nil,
		},
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputSession: &model.Session{
				ID:        12,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
			outputError: nil,
		},
		expectedStatusCode: 400,
	},
	{
		isPrepareDB:          true,
		isPrepareSM:          true,
		isGetAd:              true,
		isPrepareCheckConnSM: true,
		isCheckSession:       true,
		isSecondCheckSession: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("DELETE", domain+"/ads/delete/15", nil)
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		db: &mockDataDB{
			inputID:     15,
			outputAd:    &model.AdItem{},
			outputError: errors.New("e"),
		},
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputSession: &model.Session{
				ID:        12,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
			outputError: nil,
		},
		expectedStatusCode: 500,
	},
	{
		isPrepareDB:          true,
		isPrepareSM:          true,
		isGetAd:              true,
		isPrepareCheckConnSM: true,
		isCheckSession:       true,
		isSecondCheckSession: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("DELETE", domain+"/ads/delete/15", nil)
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		db: &mockDataDB{
			inputID:     15,
			outputAd:    adsInDB[21],
			outputError: nil,
		},
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputSession: &model.Session{
				ID:        12,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
			outputError: nil,
		},
		expectedStatusCode: 403,
	},
	{
		isPrepareDB:          true,
		isPrepareSM:          true,
		isGetAd:              true,
		isPrepareCheckConnSM: true,
		isCheckSession:       true,
		isSecondCheckSession: true,
		isRemoveAd:           true,
		request: func() *http.Request {
			r, _ := http.NewRequest("DELETE", domain+"/ads/delete/15", nil)
			r.Header.Set("Cookie", "session_id=123abc")
			return r
		}(),
		db: &mockDataDB{
			inputID:        15,
			inputAd:        adsInDB[17],
			outputID:       15,
			outputAd:       adsInDB[16],
			outputError:    nil,
			outputErrorSec: errors.New("e"),
		},
		sm: &mockDataSM{
			inputSessionID: &model.SessionID{ID: "123abc"},
			outputSession: &model.Session{
				ID:        12,
				UserAgent: "Go-http-client/1.1",
				Login:     "pet@animal.com",
			},
			outputError: nil,
		},
		expectedStatusCode: 500,
	},
	{
		request: func() *http.Request {
			img := image.NewRGBA(image.Rect(0, 0, 100, 50))
			img.Set(2, 3, color.RGBA{255, 0, 0, 255})
			os.Mkdir("images", 0777)
			f, _ := os.OpenFile(os.Getenv("CI_PROJECT_DIR")+"pkg/api/images/image.png", os.O_WRONLY|os.O_CREATE, 0777)
			png.Encode(f, img)
			r, _ := http.NewRequest("GET", domain+"/images/image.png", nil)
			f.Close()
			return r
		}(),
		expectedStatusCode: 200,
	},
	{
		request: func() *http.Request {
			r, _ := http.NewRequest("GET", domain+"/images/image123.png", nil)
			return r
		}(),
		expectedStatusCode: 400,
	},
}

func TestInterfaceOfAPI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log.SetOutput(ioutil.Discard)

	for _, tCase := range testCases {
		func() {
			mockDB := mock_model.NewMockDB(ctrl)

			// need GetAd
			if tCase.isGetAd && tCase.isPrepareDB {
				mockDB.EXPECT().GetAd(tCase.db.inputID).
					Return(tCase.db.outputAd, tCase.db.outputError)
			}

			// need GetAds
			if tCase.isGetAds && tCase.isPrepareDB {
				mockDB.EXPECT().GetAds(&model.SearchParams{
					Limit:  tCase.db.inputLimit,
					Offset: tCase.db.inputOffset,
				}).Return(tCase.db.outputAds, tCase.db.outputError)
			}

			// need GetUserWithID
			if tCase.isGetUserWithID && tCase.isPrepareDB {
				mockDB.EXPECT().GetUserWithID(tCase.db.inputID).
					Return(tCase.db.outputUser, tCase.db.outputError)
			}

			// need GetAdsOfUser
			if tCase.isGetAdsOfUser && tCase.isPrepareDB {
				mockDB.EXPECT().GetAdsOfUser(tCase.db.inputID).
					Return(tCase.db.outputAds, tCase.db.outputError)
			}

			// need NewUser
			if tCase.isNewUser && tCase.isPrepareDB {
				mockDB.EXPECT().NewUser(gomock.Any()).
					Return(tCase.db.outputID, tCase.db.outputError)
			}

			// need GetUserWithEmail
			if tCase.isGetUserWithEmail && tCase.isPrepareDB {
				mockDB.EXPECT().GetUserWithEmail(tCase.db.inputEmail).
					Return(tCase.db.outputUser, tCase.db.outputError)
			}

			// need EditUser
			if tCase.isEditUser && tCase.isPrepareDB {
				if tCase.db.inputUser == nil {
					mockDB.EXPECT().EditUser(gomock.Not(nil)).
						Return(tCase.db.outputID, tCase.db.outputError)
				} else {
					mockDB.EXPECT().EditUser(tCase.db.inputUser).
						Return(tCase.db.outputID, tCase.db.outputError)
				}
			}

			// need DeleteUser
			if tCase.isRemoveUser && tCase.isPrepareDB {
				mockDB.EXPECT().RemoveUser(tCase.db.inputID).
					Return(tCase.db.outputID, tCase.db.outputError)
			}

			// need for image processing
			if tCase.isGetUserWithIDImg {
				mockDB.EXPECT().GetUserWithID(tCase.db.inputIDimg).
					Return(tCase.db.outputUserImg, tCase.db.outputErrorImg)
			}

			// need NewAd
			if tCase.isNewAd && tCase.isPrepareDB {
				if strings.Contains(tCase.request.Header.Get("Content-Type"), "multipart") {
					mockDB.EXPECT().NewAd(gomock.Any()).
						Return(tCase.db.outputID, tCase.db.outputError)
				} else {
					mockDB.EXPECT().NewAd(tCase.db.inputAd).
						Return(tCase.db.outputID, tCase.db.outputError)
				}
			}

			// need EditAd
			if tCase.isEditAd && tCase.isPrepareDB {
				if tCase.db.outputErrorSec != nil {
					mockDB.EXPECT().EditAd(tCase.db.inputAd).
						Return(tCase.db.outputID, tCase.db.outputErrorSec)
				} else if strings.Contains(tCase.request.Header.Get("Content-Type"), "multipart") {
					mockDB.EXPECT().EditAd(gomock.Any()).
						Return(tCase.db.outputID, tCase.db.outputError)
				} else {
					mockDB.EXPECT().EditAd(tCase.db.inputAd).
						Return(tCase.db.outputID, tCase.db.outputError)
				}
			}

			// need RemoveAd
			if tCase.isRemoveAd && tCase.isPrepareDB {
				if tCase.db.outputErrorSec != nil {
					mockDB.EXPECT().RemoveAd(tCase.db.inputID).
						Return(tCase.db.outputID, tCase.db.outputErrorSec)
				} else {
					mockDB.EXPECT().RemoveAd(tCase.db.inputID).
						Return(tCase.db.outputID, tCase.db.outputError)
				}
			}

			mockSM := mock_model.NewMockSM(ctrl)

			// need CreateSession
			if tCase.isCreateSession && tCase.isPrepareSM {
				mockSM.EXPECT().CreateSession(tCase.sm.inputSession, tCase.sm.inputExpires).
					Return(tCase.sm.outputSessionID, tCase.sm.outputError)
			}

			// need CheckSession
			if tCase.isCheckSession && tCase.isPrepareSM {
				mockSM.EXPECT().CheckSession(tCase.sm.inputSessionID).
					Return(tCase.sm.outputSession, tCase.sm.outputError)
			}

			// need CheckSession second
			if tCase.isSecondCheckSession && tCase.isPrepareSM {
				mockSM.EXPECT().CheckSession(tCase.sm.inputSessionID).
					Return(tCase.sm.outputSession, tCase.sm.outputError)
			}

			// need DeleteSession
			if tCase.isDeleteSession && tCase.isPrepareSM {
				mockSM.EXPECT().DeleteSession(tCase.sm.inputSessionID).
					Return(tCase.sm.outputError)
			}

			if tCase.isPrepareCheckConnSM {
				mockSM.EXPECT().IsConnected().Return(true)
			}

			tModel := model.New(mockDB, mockSM)

			srv, ch := api.StartServer(api.Config{
				Address:      "localhost:49123",
				ReadTimeout:  "25s",
				WriteTimeout: "25s",
				IdleTimeout:  "25s",
			}, tModel)

			time.Sleep(time.Millisecond * 50) // time to start the server
			defer func() {
				srv.Shutdown(nil)
				<-ch
			}()

			// send request to server
			client := http.DefaultClient
			if tCase.request.Header.Get("Content-Type") == "" {
				tCase.request.Header.Set("Content-Type",
					"application/x-www-form-urlencoded; charset=utf-8")
			}
			client.Timeout = time.Second * 25
			result, err := client.Do(tCase.request)

			if err != nil {
				t.Fatal("Expected no error while request\nGot: ", err.Error())
			}
			defer result.Body.Close()

			if result.StatusCode != tCase.expectedStatusCode {
				buf := make([]byte, 1024)
				result.Body.Read(buf)
				fmt.Println(string(buf))
				result.Body.Close()
				t.Errorf("Expected equal status codes:\nExpected:%v\nReceived:%v",
					tCase.expectedStatusCode, result.StatusCode)
			}

			var data []byte
			// if response is one ad
			if result.StatusCode == http.StatusOK && tCase.expectedAd != nil {
				data, _ = json.Marshal(tCase.expectedAd)
			}

			// if response is several ads
			if result.StatusCode == http.StatusOK && tCase.expectedAds != nil {
				data, _ = json.Marshal(tCase.expectedAds)
			}

			// if response is user
			if result.StatusCode == http.StatusOK && tCase.expectedUser != nil {
				data, _ = json.Marshal(tCase.expectedUser)
			}

			// if response is create user
			if result.StatusCode == http.StatusCreated && tCase.expectedUserCreate != nil {
				data, _ = json.Marshal(tCase.expectedUserCreate)
			}

			// if response is login from android app
			if result.StatusCode == http.StatusOK && tCase.expectedLogin != nil {
				data, _ = json.Marshal(tCase.expectedLogin)
			}

			// if response is login
			if result.StatusCode == http.StatusOK && tCase.expectedCookieValue != "" {
				cookies := result.Cookies()
				if len(cookies) == 0 {
					t.Error("Expected set-cookie")
				} else if cookies[0].Value != tCase.expectedCookieValue {
					t.Errorf("Expected equal:\nExpected:%s\nReceived:%s",
						tCase.expectedCookieValue, cookies[0].Value)
				}
			}

			// if response is ad create
			if result.StatusCode == http.StatusOK && tCase.expectedAdCreate != nil {
				data, _ = json.Marshal(tCase.expectedAdCreate)
			}

			if data != nil {
				body, err := ioutil.ReadAll(result.Body)
				if err != nil {
					t.Error("Expected no error while reading body")
				}

				if string(body) != string(data) {
					t.Errorf("Expected equal:\nExpected:%s\nReceived:%s",
						string(data), string(body))
				}
			}
		}()
	}
}

func TestMiddleware(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mock_model.NewMockDB(ctrl)
	sm := mock_model.NewMockSM(ctrl)

	sm.EXPECT().IsConnected().Return(false)
	sm.EXPECT().TryReconnect().Return(nil)

	m := model.New(db, sm)

	srv, ch := api.StartServer(api.Config{
		Address:      "localhost:49123",
		ReadTimeout:  "25s",
		WriteTimeout: "25s",
		IdleTimeout:  "25s",
	}, m)

	time.Sleep(time.Millisecond * 50) // time to start the server

	res, _ := http.Post(domain+"/users/logout", "application/x-www-form-urlencoded; charset=utf-8", nil)
	if res.StatusCode != 200 {
		t.Error("Expected status 200")
	}

	srv.Shutdown(nil)
	<-ch

	sm.EXPECT().IsConnected().Return(false)
	sm.EXPECT().TryReconnect().Return(errors.New("e"))

	m = model.New(db, sm)

	srv, ch = api.StartServer(api.Config{
		Address:      "localhost:49123",
		ReadTimeout:  "25s",
		WriteTimeout: "25s",
		IdleTimeout:  "25s",
	}, m)

	time.Sleep(time.Millisecond * 50) // time to start the server

	res, _ = http.Post(domain+"/users/logout", "application/x-www-form-urlencoded; charset=utf-8", nil)
	if res.StatusCode != 500 {
		t.Error("Expected status 500")
	}

	srv.Shutdown(nil)
	<-ch
}
