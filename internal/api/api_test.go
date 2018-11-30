package api_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"bmstu.codes/developers34/SBWeb/internal/api"
	"bmstu.codes/developers34/SBWeb/internal/api/mock_model"

	"github.com/golang/mock/gomock"

	"bmstu.codes/developers34/SBWeb/internal/model"
)

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
	15: &model.AdItem{
		ID:    15,
		Title: "Building",
		City:  "Moscow",
		User: model.User{
			ID:        12,
			FirstName: "Ivan",
			LastName:  "Ivanov",
			Email:     "Ivan@ivanov.com",
		},
		Description: "Awesome",
	},
	21: &model.AdItem{
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
	54: &model.AdItem{
		ID:    54,
		Title: "House",
		City:  "New York",
		User: model.User{
			ID:        17,
			FirstName: "Petya",
			LastName:  "Succeded",
			Email:     "pet@animal.com",
		},
		Description: "very good house",
	},
}

var usersInDB = map[int64]*model.User{
	12: &model.User{
		ID:        12,
		FirstName: "Ivan",
		LastName:  "Ivanov",
		Email:     "Ivan@ivanov.com",
		Password:  "$2a$10$kCw55AyZuzo3u6GicdIjg.RKscZZ0IxvZMSKIiuY.fQ7R2F9OPbba",
	},
	17: &model.User{
		ID:        17,
		FirstName: "Ivan",
		LastName:  "Ivanov",
		Email:     "pet@animal.com",
		Password:  "$2a$10$kCw55AyZuzo3u6GicdIjg.RKscZZ0IxvZMSKIiuY.fQ7R2F9OPbba",
	},
	18: &model.User{
		ID:        18,
		FirstName: "John",
		LastName:  "Johnov",
		Email:     "john@johnov.com",
	},
}

type apiError struct {
	Description string `json:"description"`
	Message     string `json:"message"`
	ErrorCode   string `json:"error"`
}

type createUserResp struct {
	ID  int64
	Ref string
}

type loginUserResp struct {
	ID        int64  `json:"id,"`
	FirstName string `json:"first_name,"`
	LastName  string `json:"last_name,"`
}

type mockDataDB struct {
	inputID     int64
	inputLimit  int
	inputOffset int
	inputQuery  string
	inputEmail  string
	inputUser   *model.User
	inputAd     *model.User

	inputIDimg int64

	outputErrorImg error
	outputUserImg  *model.User

	outputUser       *model.User
	outputUserCreate *createUserResp
	outputAd         *model.AdItem
	outputAds        []*model.AdItem
	outputError      error
	outputID         int64
}

type mockDataSM struct {
	inputSession   *model.Session
	inputSessionID *model.SessionID
	inputExpires   bool

	outputSession   *model.Session
	outputSessionID *model.SessionID
	outputError     error
}

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

	// flags of actions
	isLoginPage bool

	// flags for sm
	isCreateSession      bool
	isCheckSession       bool
	isSecondCheckSession bool
	isDeleteSession      bool

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
}

var testCases = []testCase{
	testCase{
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
	testCase{
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
	testCase{
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
	testCase{
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
	testCase{
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
	testCase{
		isGetAds: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("GET", domain+"/ads?offset=1&limit=12&query=\003", nil)
			return r
		}(),
		expectedStatusCode: 400,
	},
	testCase{
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
	testCase{
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
	testCase{
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
	testCase{
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
	testCase{
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
	testCase{
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
	testCase{
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
	testCase{
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
	testCase{
		isNewUser: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/new",
				strings.NewReader("%email=pet@animal.com&first_name=Ivan&last_name=Ivanov&password=123456"))
			return r
		}(),
		expectedStatusCode: http.StatusBadRequest,
	},
	testCase{
		isNewUser: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/new",
				strings.NewReader("script=bad&email=pet@animal.com&first_name=Ivan&last_name=Ivanov&password=123456"))
			return r
		}(),
		expectedStatusCode: http.StatusBadRequest,
	},
	testCase{
		isNewUser: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/new",
				strings.NewReader("email=pet@animal.com&first_name=&last_name=Ivanov&password=123456"))
			return r
		}(),
		expectedStatusCode: http.StatusBadRequest,
	},
	testCase{
		isNewUser: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/new",
				strings.NewReader("email=pet@animal.om&first_name=Ivan&last_name=Ivanov&password=123456"))
			return r
		}(),
		expectedStatusCode: http.StatusBadRequest,
	},
	testCase{
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
	testCase{
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
	testCase{
		isGetUserWithEmail: true,
		isCreateSession:    true,
		isPrepareDB:        true,
		isPrepareSM:        true,
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
	testCase{
		isGetUserWithEmail: true,
		isCreateSession:    true,
		isPrepareDB:        true,
		isPrepareSM:        true,
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
	testCase{
		isGetUserWithEmail: true,
		isCreateSession:    true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/login",
				strings.NewReader("%email=pet@animal.com&password=123456"))
			return r
		}(),
		expectedStatusCode: 400,
	},
	testCase{
		isGetUserWithEmail: true,
		isCreateSession:    true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/login",
				strings.NewReader("email=pet@animal.cv&password=123456"))
			return r
		}(),
		expectedStatusCode: 400,
	},
	testCase{
		isGetUserWithEmail: true,
		isCreateSession:    true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/login",
				strings.NewReader("script=vvvv&email=pet@animal.com&password=123456"))
			return r
		}(),
		expectedStatusCode: 400,
	},
	testCase{
		isGetUserWithEmail: true,
		isCreateSession:    true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/login",
				strings.NewReader("email=pet@animal.com&password="))
			return r
		}(),
		expectedStatusCode: 400,
	},
	testCase{
		isGetUserWithEmail: true,
		isCreateSession:    true,
		isPrepareDB:        true,
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
	testCase{
		isGetUserWithEmail: true,
		isCreateSession:    true,
		isPrepareDB:        true,
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
	testCase{
		isGetUserWithEmail: true,
		isCreateSession:    true,
		isPrepareDB:        true,
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
	testCase{
		isGetUserWithEmail: true,
		isCreateSession:    true,
		isPrepareDB:        true,
		isPrepareSM:        true,
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
	testCase{
		isDeleteSession: true,
		isPrepareSM:     true,
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
	testCase{
		isDeleteSession: true,
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/logout", nil)
			return r
		}(),
		expectedStatusCode: 200,
	},
	testCase{
		isDeleteSession: true,
		isPrepareSM:     true,
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
	testCase{
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
	testCase{
		isCheckSession: true,
		isPrepareSM:    true,
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
	testCase{
		isCheckSession: true,
		isPrepareSM:    true,
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
	testCase{
		isCheckSession: true,
		isPrepareSM:    true,
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
	testCase{
		isCheckSession: true,
		isPrepareSM:    true,
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
	testCase{
		isCheckSession: true,
		isPrepareSM:    true,
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
	testCase{
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
	testCase{
		request: func() *http.Request {
			r, _ := http.NewRequest("POST", domain+"/users/profile",
				strings.NewReader("first_name=Alex&last_name=Ivanov"))
			return r
		}(),
		expectedStatusCode: 401,
	},
	testCase{
		isCheckSession: true,
		isPrepareSM:    true,
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
	testCase{
		isGetUserWithID:      true,
		isPrepareDB:          true,
		isPrepareSM:          true,
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
	testCase{
		isGetUserWithID:      true,
		isPrepareDB:          true,
		isPrepareSM:          true,
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
	testCase{
		isDeleteSession:      true,
		isRemoveUser:         true,
		isPrepareDB:          true,
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
	testCase{
		isRemoveUser:         true,
		isPrepareDB:          true,
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
}

func TestThis(t *testing.T) {
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
				mockDB.EXPECT().EditUser(tCase.db.inputUser).
					Return(tCase.db.outputID, tCase.db.outputError)
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

			tModel := model.New(mockDB, mockSM)

			srv, ch := api.StartServer(api.Config{
				Address:      "localhost:49123",
				ReadTimeout:  "25s",
				WriteTimeout: "25s",
				IdleTimeout:  "25s",
			}, tModel)

			time.Sleep(time.Millisecond * 10) // time to start the server
			defer func() {
				srv.Shutdown(nil)
				<-ch
			}()

			client := http.DefaultClient
			tCase.request.Header.Set("Content-Type",
				"application/x-www-form-urlencoded; charset=utf-8")
			client.Timeout = time.Second * 25
			result, err := client.Do(tCase.request)

			if err != nil {
				t.Fatal("Expected no error while request\nGot: ", err.Error())
			}
			defer result.Body.Close()

			if result.StatusCode != tCase.expectedStatusCode {
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
