package api_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"bmstu.codes/developers34/SBWeb/internal/api"
	"bmstu.codes/developers34/SBWeb/internal/api/mock_model"

	"github.com/golang/mock/gomock"

	"bmstu.codes/developers34/SBWeb/internal/model"
)

const (
	domain = "http://localhost:54000"
)

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
	inputEmail  string
	inputUser   *model.User
	inputAd     *model.User

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
	// flags
	isGetAd         bool
	isGetAds        bool
	isGetAdsOfUser  bool
	isGetUserWithID bool
	isNewUser       bool

	// expect flags
	isPrepareDB bool
	isPrepareSM bool

	// request of client
	request *http.Request

	// data for functions of DB and SM
	sm *mockDataSM
	db *mockDataDB

	// expected response of server
	expectedStatusCode int
	expectedUser       *model.User
	expectedAd         *model.AdItem
	expectedAds        []*model.AdItem
	expectedUserCreate *createUserResp
	expectedLogin      *loginUserResp
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
			r, _ := http.NewRequest("POST", domain+"/users/18", nil)
			return r
		}(),
	},
}

func TestThis(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log.SetOutput(ioutil.Discard)

	for _, tCase := range testCases {
		func() {
			mockDB := mock_model.NewMockDB(ctrl)
			mockSM := mock_model.NewMockSM(ctrl)

			// need GetAd
			if tCase.isGetAd && tCase.isPrepareDB {
				mockDB.EXPECT().GetAd(tCase.db.inputID).Return(tCase.db.outputAd, tCase.db.outputError)
			}

			// need GetAds
			if tCase.isGetAds && tCase.isPrepareDB {
				mockDB.EXPECT().GetAds(tCase.db.inputLimit, tCase.db.inputOffset).Return(tCase.db.outputAds, tCase.db.outputError)
			}

			// need GetUserWithID
			if tCase.isGetUserWithID && tCase.isPrepareDB {
				mockDB.EXPECT().GetUserWithID(tCase.db.inputID).Return(tCase.db.outputUser, tCase.db.outputError)
			}

			// need GetAdsOfUser
			if tCase.isGetAdsOfUser && tCase.isPrepareDB {
				mockDB.EXPECT().GetAdsOfUser(tCase.db.inputID).Return(tCase.db.outputAds, tCase.db.outputError)
			}

			// need NewUser
			if tCase.isNewUser && tCase.isPrepareDB {
				mockDB.EXPECT().NewUser(tCase.db.inputUser).Return(tCase.db.outputID, tCase.db.outputError)
			}

			tModel := model.New(mockDB, mockSM)

			srv, ch := api.StartServer(api.Config{
				Address:      "localhost:54000",
				ReadTimeout:  "10s",
				WriteTimeout: "10s",
				IdleTimeout:  "10s",
			}, tModel)

			time.Sleep(time.Millisecond * 10) // time to start the server
			defer func() {
				srv.Shutdown(nil)
				<-ch
			}()

			client := http.DefaultClient
			tCase.request.Header.Set("Content-Type",
				"application/x-www-form-urlencoded; charset=utf-8")
			result, err := client.Do(tCase.request)

			if err != nil {
				t.Fatal("Expected no error while request\nGot: ", err.Error())
			}
			defer result.Body.Close()

			if result.StatusCode != tCase.expectedStatusCode {
				t.Errorf("Expected equal status codes:\nExpected:%v\nReceived:%v",
					tCase.expectedStatusCode, result.StatusCode)
			}

			// if response is one ad
			if result.StatusCode == http.StatusOK && tCase.expectedAd != nil {
				adDataExpected, _ := json.Marshal(tCase.expectedAd)

				body, err := ioutil.ReadAll(result.Body)
				if err != nil {
					t.Error("Expected no error while reading body")
				}

				if string(body) != string(adDataExpected) {
					t.Errorf("Expected equal ads:\nExpected:%s\nReceived:%s",
						string(adDataExpected), string(body))
				}
			}

			// if response is several ads
			if result.StatusCode == http.StatusOK && tCase.expectedAds != nil {
				adsDataExpected, _ := json.Marshal(tCase.expectedAds)

				body, err := ioutil.ReadAll(result.Body)
				if err != nil {
					t.Error("Expected no error while reading body")
				}

				if string(body) != string(adsDataExpected) {
					t.Errorf("Expected equal ads:\nExpected:%s\nReceived:%s",
						string(adsDataExpected), string(body))
				}
			}

			// if response is user
			if result.StatusCode == http.StatusOK && tCase.expectedUser != nil {
				userDataExpected, _ := json.Marshal(tCase.expectedUser)

				body, err := ioutil.ReadAll(result.Body)
				if err != nil {
					t.Error("Expected no error while reading body")
				}

				if string(body) != string(userDataExpected) {
					t.Errorf("Expected equal ads:\nExpected:%s\nReceived:%s",
						string(userDataExpected), string(body))
				}
			}

			// if response is create user
			if result.StatusCode == http.StatusOK && tCase.expectedUserCreate != nil {
				userDataExpected, _ := json.Marshal(tCase.expectedUserCreate)

				body, err := ioutil.ReadAll(result.Body)
				if err != nil {
					t.Error("Expected no error while reading body")
				}

				if string(body) != string(userDataExpected) {
					t.Errorf("Expected equal ads:\nExpected:%s\nReceived:%s",
						string(userDataExpected), string(body))
				}
			}
		}()
	}
}
