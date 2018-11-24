package api_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"bmstu.codes/developers34/SBWeb/internal/api"
	mock_model "bmstu.codes/developers34/SBWeb/internal/api/mock_model"
	"bmstu.codes/developers34/SBWeb/internal/model"
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

type mockReadInfo struct {
	expectedAd      *model.AdItem
	expectedAds     []*model.AdItem
	expectedUser    *model.User
	expectedError   error
	expectedID      int64
	expectedTocken  string
	expectedErrorSM error
	expectedSession *model.Session
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

type testCase struct {
	providedID       int64
	providedOffset   int
	providedLimit    int
	providedParams   string
	providedUser     *model.User
	providedEmail    string
	providedPassword string
	providedAgent    string

	isOneAd                bool
	isSeveralAds           bool
	isOneUser              bool
	isSeveralAdsOfUser     bool
	isCreateUser           bool
	isNeedSetExpect        bool
	isNeedSetExpectSM      bool
	isLogin                bool
	isAnroid               bool
	isLogout               bool
	isUpdateUser           bool
	isNeedSetAnotherExpect bool
	isDeleteUser           bool

	expectedReturnAd         model.AdItem
	expectedReturnAds        []model.AdItem
	expectedReturnUser       model.User
	expectedReturnUserCreate createUserResp
	expectedReturnUserLogin  loginUserResp
	expectedReturnError      apiError
	expectedStatusCode       int

	mockReadInfo mockReadInfo
	addr         string
}

func TestStartWithBadConfig(t *testing.T) {
	srv, ch := api.StartServer(api.Config{
		Address:      "localhost:54000",
		ReadTimeout:  "10eregeurgi",
		WriteTimeout: "10s",
		IdleTimeout:  "10s",
	}, nil)
	defer srv.Shutdown(nil)
	err := <-ch
	if err == http.ErrServerClosed {
		t.Error("Expected other error")
	}
}

var testCases = []testCase{
	testCase{
		isOneAd:            true,
		providedID:         15,
		expectedReturnAd:   *adsInDB[15],
		expectedStatusCode: http.StatusOK,
		mockReadInfo: mockReadInfo{
			expectedAd:    adsInDB[15],
			expectedError: nil,
		},
		addr: "/ads/15",
	},
	testCase{
		isOneAd:            true,
		providedID:         14,
		expectedStatusCode: http.StatusBadRequest,
		mockReadInfo: mockReadInfo{
			expectedAd: &model.AdItem{
				ID: -1,
			},
			expectedError: sql.ErrNoRows,
		},
		addr: "/ads/14",
	},
	testCase{
		isOneAd:            true,
		providedID:         14,
		expectedStatusCode: http.StatusInternalServerError,
		mockReadInfo: mockReadInfo{
			expectedAd:    &model.AdItem{},
			expectedError: errors.New("some err"),
		},
		addr: "/ads/14",
	},
	testCase{
		isSeveralAds:       true,
		providedOffset:     2,
		providedLimit:      2,
		expectedStatusCode: http.StatusInternalServerError,
		mockReadInfo: mockReadInfo{
			expectedAds:   make([]*model.AdItem, 0),
			expectedError: errors.New("some err"),
		},
		addr: "/ads?offset=2&limit=2",
	},
	testCase{
		expectedReturnAds:  []model.AdItem{*adsInDB[21], *adsInDB[54]},
		isSeveralAds:       true,
		providedOffset:     2,
		providedLimit:      2,
		expectedStatusCode: http.StatusOK,
		mockReadInfo: mockReadInfo{
			expectedAds:   []*model.AdItem{adsInDB[21], adsInDB[54]},
			expectedError: nil,
		},
		addr: "/ads?offset=2&limit=2",
	},
	testCase{
		expectedReturnAds:  []model.AdItem{*adsInDB[21], *adsInDB[54]},
		isSeveralAds:       true,
		providedOffset:     0,
		providedLimit:      15,
		expectedStatusCode: http.StatusOK,
		mockReadInfo: mockReadInfo{
			expectedAds:   []*model.AdItem{adsInDB[21], adsInDB[54]},
			expectedError: nil,
		},
		addr: "/ads",
	},
	testCase{
		providedID:         12,
		expectedReturnAds:  []model.AdItem{*adsInDB[21], *adsInDB[54]},
		isSeveralAdsOfUser: true,
		expectedStatusCode: http.StatusOK,
		mockReadInfo: mockReadInfo{
			expectedAds:   []*model.AdItem{adsInDB[21], adsInDB[54]},
			expectedError: nil,
		},
		addr: "/users/12?show_ads=true",
	},
	testCase{
		providedID:         12,
		isSeveralAdsOfUser: true,
		expectedStatusCode: http.StatusInternalServerError,
		mockReadInfo: mockReadInfo{
			expectedAds:   make([]*model.AdItem, 0),
			expectedError: errors.New("some err"),
		},
		addr: "/users/12?show_ads=true",
	},
	testCase{
		providedID:         12,
		isOneUser:          true,
		expectedStatusCode: http.StatusInternalServerError,
		mockReadInfo: mockReadInfo{
			expectedUser:  &model.User{},
			expectedError: errors.New("some err"),
		},
		addr: "/users/12",
	},
	testCase{
		providedID:         12,
		isOneUser:          true,
		expectedStatusCode: http.StatusBadRequest,
		mockReadInfo: mockReadInfo{
			expectedUser:  &model.User{ID: -1},
			expectedError: sql.ErrNoRows,
		},
		addr: "/users/12",
	},
	testCase{
		providedID:         12,
		isOneUser:          true,
		expectedReturnUser: *usersInDB[12],
		expectedStatusCode: http.StatusOK,
		mockReadInfo: mockReadInfo{
			expectedUser:  usersInDB[12],
			expectedError: nil,
		},
		addr: "/users/12",
	},
	testCase{
		isNeedSetExpect:          true,
		isCreateUser:             true,
		providedParams:           "first_name=Ivan&last_name=Ivanov&email=Ivan@ivanov.com&password=123456",
		expectedReturnUserCreate: createUserResp{ID: 12, Ref: "/users/12"},
		expectedStatusCode:       http.StatusCreated,
		mockReadInfo: mockReadInfo{
			expectedID:    12,
			expectedError: nil,
		},
		addr: "/users/new",
	},
	testCase{
		isNeedSetExpect:    false,
		isCreateUser:       true,
		providedParams:     "%first_name=Ivan&last_name=Ivanov&email=Ivan@ivanov.com&password=123456",
		expectedStatusCode: http.StatusBadRequest,
		addr:               "/users/new",
	},
	testCase{
		isNeedSetExpect:    false,
		isCreateUser:       true,
		providedParams:     "script='DROP TABLE';first_name=Ivan&last_name=Ivanov&email=Ivan@ivanov.com&password=123456",
		expectedStatusCode: http.StatusInternalServerError,
		addr:               "/users/new",
	},
	testCase{
		isNeedSetExpect:    false,
		isCreateUser:       true,
		providedParams:     "last_name=Ivanov&email=Ivan@ivanov.com&password=123456",
		expectedStatusCode: http.StatusBadRequest,
		addr:               "/users/new",
	},
	testCase{
		isNeedSetExpect:    false,
		isCreateUser:       true,
		providedParams:     "first_name=Ivan&last_name=Iva124nov&email=Ivan@ivanov.com&password=123456",
		expectedStatusCode: http.StatusBadRequest,
		addr:               "/users/new",
	},
	testCase{
		isNeedSetExpect:    true,
		isCreateUser:       true,
		providedParams:     "first_name=Ivan&last_name=Ivanov&email=Ivan@ivanov.com&password=123456",
		expectedStatusCode: http.StatusBadRequest,
		mockReadInfo: mockReadInfo{
			expectedID:    -1,
			expectedError: errors.New("some err"),
		},
		addr: "/users/new",
	},
	testCase{
		isNeedSetExpect:    true,
		isCreateUser:       true,
		providedParams:     "first_name=Ivan&last_name=Ivanov&email=Ivan@ivanov.com&password=123456",
		expectedStatusCode: http.StatusInternalServerError,
		mockReadInfo: mockReadInfo{
			expectedID:    0,
			expectedError: errors.New("some err"),
		},
		addr: "/users/new",
	},
	testCase{
		isNeedSetExpect:    true,
		isNeedSetExpectSM:  true,
		isLogin:            true,
		isAnroid:           false,
		providedUser:       usersInDB[12],
		providedParams:     "email=Ivan@ivanov.com&password=123456",
		providedAgent:      "Go-http-client/1.1",
		expectedStatusCode: http.StatusOK,
		mockReadInfo: mockReadInfo{
			expectedTocken:  "",
			expectedUser:    usersInDB[12],
			expectedError:   nil,
			expectedErrorSM: nil,
		},
		addr: "/users/login",
	},
	testCase{
		isNeedSetExpect:    true,
		isNeedSetExpectSM:  true,
		isLogin:            true,
		isAnroid:           true,
		providedUser:       usersInDB[12],
		providedParams:     "email=Ivan@ivanov.com&password=123456",
		providedAgent:      "Android_app",
		expectedStatusCode: http.StatusOK,
		expectedReturnUserLogin: loginUserResp{
			ID:        12,
			FirstName: "Ivan",
			LastName:  "Ivanov",
		},
		mockReadInfo: mockReadInfo{
			expectedTocken:  "123",
			expectedUser:    usersInDB[12],
			expectedError:   nil,
			expectedErrorSM: nil,
		},
		addr: "/users/login",
	},
	testCase{
		isLogin:            true,
		providedParams:     "%email=Ivan@ivanov.com&password=123456",
		expectedStatusCode: http.StatusBadRequest,
		addr:               "/users/login",
	},
	testCase{
		isLogin:            true,
		providedParams:     "script=eer&email=Ivan@ivanov.com&password=123456",
		expectedStatusCode: http.StatusInternalServerError,
		addr:               "/users/login",
	},
	testCase{
		isLogin:            true,
		providedParams:     "email=Ivan@ivanov.com&password=",
		expectedStatusCode: http.StatusBadRequest,
		addr:               "/users/login",
	},
	testCase{
		isLogin:            true,
		providedParams:     "email=Ivan@ivanov.c&password=123456",
		expectedStatusCode: http.StatusBadRequest,
		addr:               "/users/login",
	},
	testCase{
		isLogin:            true,
		isNeedSetExpect:    true,
		providedParams:     "email=Ivan@ivanov.com&password=123456",
		expectedStatusCode: http.StatusBadRequest,
		providedUser:       usersInDB[12],
		addr:               "/users/login",
		mockReadInfo: mockReadInfo{
			expectedUser:  &model.User{ID: -1},
			expectedError: errors.New("some error"),
		},
	},
	testCase{
		isLogin:            true,
		isNeedSetExpect:    true,
		providedParams:     "email=Ivan@ivanov.com&password=123456",
		expectedStatusCode: http.StatusInternalServerError,
		providedUser:       usersInDB[12],
		addr:               "/users/login",
		mockReadInfo: mockReadInfo{
			expectedUser:  &model.User{},
			expectedError: errors.New("some error"),
		},
	},
	testCase{
		isLogin:            true,
		isNeedSetExpect:    true,
		providedParams:     "email=Ivan@ivanov.com&password=1234567",
		expectedStatusCode: http.StatusBadRequest,
		providedUser:       usersInDB[12],
		addr:               "/users/login",
		mockReadInfo: mockReadInfo{
			expectedUser:  usersInDB[12],
			expectedError: nil,
		},
	},
	testCase{
		isLogin:            true,
		isNeedSetExpect:    true,
		isNeedSetExpectSM:  true,
		providedParams:     "email=Ivan@ivanov.com&password=123456",
		expectedStatusCode: http.StatusInternalServerError,
		providedAgent:      "Go-http-client/1.1",
		providedUser:       usersInDB[12],
		addr:               "/users/login",
		mockReadInfo: mockReadInfo{
			expectedUser:    usersInDB[12],
			expectedError:   nil,
			expectedTocken:  "",
			expectedErrorSM: errors.New("some err"),
		},
	},
	testCase{
		isNeedSetExpectSM:  true,
		isLogout:           true,
		providedParams:     "123456789",
		expectedStatusCode: http.StatusOK,
		providedAgent:      "Go-http-client/1.1",
		addr:               "/users/logout",
		mockReadInfo: mockReadInfo{
			expectedTocken:  "123456789",
			expectedErrorSM: nil,
		},
	},
	testCase{
		isLogout:           true,
		providedParams:     "",
		expectedStatusCode: http.StatusOK,
		providedAgent:      "Go-http-client/1.1",
		addr:               "/users/logout",
	},
	testCase{
		isNeedSetExpectSM:  true,
		isLogout:           true,
		providedParams:     "1512",
		expectedStatusCode: http.StatusOK,
		providedAgent:      "Go-http-client/1.1",
		addr:               "/users/logout",
		mockReadInfo: mockReadInfo{
			expectedTocken:  "1512",
			expectedErrorSM: errors.New("some err"),
		},
	},
	testCase{
		isNeedSetExpectSM:      true,
		isNeedSetExpect:        true,
		isUpdateUser:           true,
		isNeedSetAnotherExpect: true,
		providedParams:         "first_name=Ivan&last_name=Ivanov&email=pet@animal.com",
		providedUser:           usersInDB[17],
		expectedStatusCode:     http.StatusOK,
		providedAgent:          "Go-http-client/1.1",
		addr:                   "/users/profile",
		mockReadInfo: mockReadInfo{
			expectedID:      17,
			expectedTocken:  "1512",
			expectedErrorSM: nil,
			expectedError:   nil,
			expectedSession: &model.Session{
				ID:        17,
				Login:     usersInDB[17].Email,
				UserAgent: "Go-http-client/1.1",
			},
		},
	},
	testCase{
		isNeedSetExpectSM:      true,
		isUpdateUser:           true,
		isNeedSetAnotherExpect: false,
		providedParams:         "%first_name=Ivan&last_name=Ivanov&email=pet@animal.com",
		providedUser:           usersInDB[17],
		expectedStatusCode:     http.StatusBadRequest,
		providedAgent:          "Go-http-client/1.1",
		addr:                   "/users/profile",
		mockReadInfo: mockReadInfo{
			expectedID:      17,
			expectedTocken:  "1512",
			expectedErrorSM: nil,
			expectedError:   nil,
			expectedSession: &model.Session{
				ID:        17,
				Login:     usersInDB[17].Email,
				UserAgent: "Go-http-client/1.1",
			},
		},
	},
	testCase{
		isNeedSetExpectSM:      true,
		isUpdateUser:           true,
		isNeedSetAnotherExpect: false,
		providedParams:         "first_name=&last_name=Ivanov&email=pet@animal.com",
		providedUser:           usersInDB[17],
		expectedStatusCode:     http.StatusBadRequest,
		providedAgent:          "Go-http-client/1.1",
		addr:                   "/users/profile",
		mockReadInfo: mockReadInfo{
			expectedID:      17,
			expectedTocken:  "1512",
			expectedErrorSM: nil,
			expectedError:   nil,
			expectedSession: &model.Session{
				ID:        17,
				Login:     usersInDB[17].Email,
				UserAgent: "Go-http-client/1.1",
			},
		},
	},
	testCase{
		isNeedSetExpectSM:      true,
		isUpdateUser:           true,
		isNeedSetAnotherExpect: false,
		providedParams:         "first_name=15fwe5&last_name=Ivanov&email=pet@animal.com",
		providedUser:           usersInDB[17],
		expectedStatusCode:     http.StatusBadRequest,
		providedAgent:          "Go-http-client/1.1",
		addr:                   "/users/profile",
		mockReadInfo: mockReadInfo{
			expectedID:      17,
			expectedTocken:  "1512",
			expectedErrorSM: nil,
			expectedError:   nil,
			expectedSession: &model.Session{
				ID:        17,
				Login:     usersInDB[17].Email,
				UserAgent: "Go-http-client/1.1",
			},
		},
	},
	testCase{
		isNeedSetExpectSM:      true,
		isUpdateUser:           true,
		isNeedSetAnotherExpect: false,
		providedParams:         "first_name=Ivan&script=ffff&last_name=Ivanov&email=pet@animal.com",
		providedUser:           usersInDB[17],
		expectedStatusCode:     http.StatusInternalServerError,
		providedAgent:          "Go-http-client/1.1",
		addr:                   "/users/profile",
		mockReadInfo: mockReadInfo{
			expectedID:      17,
			expectedTocken:  "1512",
			expectedErrorSM: nil,
			expectedError:   nil,
			expectedSession: &model.Session{
				ID:        17,
				Login:     usersInDB[17].Email,
				UserAgent: "Go-http-client/1.1",
			},
		},
	},
	testCase{
		isNeedSetExpectSM:      true,
		isUpdateUser:           true,
		isNeedSetAnotherExpect: true,
		isNeedSetExpect:        true,
		providedParams:         "first_name=Ivan&last_name=Ivanov&email=pet@animal.com",
		providedUser:           usersInDB[17],
		expectedStatusCode:     http.StatusInternalServerError,
		providedAgent:          "Go-http-client/1.1",
		addr:                   "/users/profile",
		mockReadInfo: mockReadInfo{
			expectedID:      17,
			expectedTocken:  "1512",
			expectedErrorSM: nil,
			expectedError:   errors.New("err"),
			expectedSession: &model.Session{
				ID:        17,
				Login:     usersInDB[17].Email,
				UserAgent: "Go-http-client/1.1",
			},
		},
	},
	testCase{
		isNeedSetExpectSM:      true,
		isUpdateUser:           true,
		isOneUser:              true,
		isNeedSetAnotherExpect: true,
		isNeedSetExpect:        true,
		providedID:             17,
		providedUser:           usersInDB[17],
		expectedStatusCode:     http.StatusOK,
		providedAgent:          "Go-http-client/1.1",
		addr:                   "/users/profile",
		expectedReturnUser:     *usersInDB[17],
		mockReadInfo: mockReadInfo{
			expectedID:      17,
			expectedUser:    usersInDB[17],
			expectedTocken:  "1512",
			expectedErrorSM: nil,
			expectedError:   nil,
			expectedSession: &model.Session{
				ID:        17,
				Login:     usersInDB[17].Email,
				UserAgent: "Go-http-client/1.1",
			},
		},
	},
	testCase{
		isNeedSetExpectSM:      true,
		isUpdateUser:           true,
		isOneUser:              true,
		isNeedSetAnotherExpect: true,
		isNeedSetExpect:        true,
		providedID:             17,
		providedUser:           usersInDB[17],
		expectedStatusCode:     http.StatusInternalServerError,
		providedAgent:          "Go-http-client/1.1",
		addr:                   "/users/profile",
		expectedReturnUser:     *usersInDB[17],
		mockReadInfo: mockReadInfo{
			expectedID:      17,
			expectedUser:    usersInDB[17],
			expectedTocken:  "1512",
			expectedErrorSM: nil,
			expectedError:   errors.New("err"),
			expectedSession: &model.Session{
				ID:        17,
				Login:     usersInDB[17].Email,
				UserAgent: "Go-http-client/1.1",
			},
		},
	},
	testCase{
		isNeedSetExpectSM:      true,
		isUpdateUser:           true,
		isNeedSetAnotherExpect: true,
		isNeedSetExpect:        true,
		isDeleteUser:           true,
		providedID:             17,
		providedUser:           usersInDB[17],
		expectedStatusCode:     http.StatusOK,
		providedAgent:          "Go-http-client/1.1",
		addr:                   "/users/profile",
		expectedReturnUser:     *usersInDB[17],
		mockReadInfo: mockReadInfo{
			expectedID:      17,
			expectedUser:    usersInDB[17],
			expectedTocken:  "1512",
			expectedErrorSM: nil,
			expectedError:   nil,
			expectedSession: &model.Session{
				ID:        17,
				Login:     usersInDB[17].Email,
				UserAgent: "Go-http-client/1.1",
			},
		},
	},
	testCase{
		isNeedSetExpectSM:      true,
		isUpdateUser:           true,
		isNeedSetAnotherExpect: true,
		isNeedSetExpect:        true,
		isDeleteUser:           true,
		providedID:             17,
		providedUser:           usersInDB[17],
		expectedStatusCode:     http.StatusInternalServerError,
		providedAgent:          "Go-http-client/1.1",
		addr:                   "/users/profile",
		expectedReturnUser:     *usersInDB[17],
		mockReadInfo: mockReadInfo{
			expectedID:      17,
			expectedUser:    usersInDB[17],
			expectedTocken:  "1512",
			expectedErrorSM: nil,
			expectedError:   errors.New("err"),
			expectedSession: &model.Session{
				ID:        17,
				Login:     usersInDB[17].Email,
				UserAgent: "Go-http-client/1.1",
			},
		},
	},
}

func TestReadInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, testItem := range testCases {
		func(testItem testCase) {

			mockDB := mock_model.NewMockDB(ctrl)

			if testItem.isOneAd {
				mockDB.EXPECT().GetAd(testItem.providedID).
					Return(testItem.mockReadInfo.expectedAd, testItem.mockReadInfo.expectedError)
			} else if testItem.isSeveralAds {
				mockDB.EXPECT().GetAds(testItem.providedLimit, testItem.providedOffset).
					Return(testItem.mockReadInfo.expectedAds, testItem.mockReadInfo.expectedError)
			} else if testItem.isSeveralAdsOfUser {
				mockDB.EXPECT().GetAdsOfUser(testItem.providedID).
					Return(testItem.mockReadInfo.expectedAds, testItem.mockReadInfo.expectedError)
			} else if testItem.isOneUser {
				mockDB.EXPECT().GetUserWithID(testItem.providedID).
					Return(testItem.mockReadInfo.expectedUser, testItem.mockReadInfo.expectedError)
			} else if testItem.isNeedSetExpect && testItem.isCreateUser {
				mockDB.EXPECT().NewUser(gomock.Any()).
					Return(testItem.mockReadInfo.expectedID, testItem.mockReadInfo.expectedError)
			} else if testItem.isNeedSetExpect && testItem.isLogin {
				mockDB.EXPECT().GetUserWithEmail(testItem.providedUser.Email).
					Return(testItem.mockReadInfo.expectedUser, testItem.mockReadInfo.expectedError)
			} else if testItem.isUpdateUser && testItem.isNeedSetExpect && !testItem.isDeleteUser {
				user := testItem.providedUser
				mockDB.EXPECT().EditUser(user).
					Return(testItem.mockReadInfo.expectedID, testItem.mockReadInfo.expectedError)
			} else if testItem.isDeleteUser {
				mockDB.EXPECT().RemoveUser(testItem.mockReadInfo.expectedID).
					Return(testItem.mockReadInfo.expectedID, testItem.mockReadInfo.expectedError)
			}

			mockSM := mock_model.NewMockSM(ctrl)
			if testItem.isLogin && testItem.isNeedSetExpectSM {
				mockSM.EXPECT().CreateSession(&model.Session{
					ID:        testItem.providedUser.ID,
					Login:     testItem.providedUser.Email,
					UserAgent: testItem.providedAgent,
				}, !testItem.isAnroid).Return(
					&model.SessionID{ID: testItem.mockReadInfo.expectedTocken},
					testItem.mockReadInfo.expectedErrorSM)
			} else if testItem.isLogout && testItem.isNeedSetExpectSM &&
				testItem.providedParams != "" {
				mockSM.EXPECT().DeleteSession(
					&model.SessionID{ID: testItem.mockReadInfo.expectedTocken}).
					Return(testItem.mockReadInfo.expectedErrorSM)
			} else if testItem.isUpdateUser && testItem.isNeedSetExpectSM {
				mockSM.EXPECT().CheckSession(&model.SessionID{ID: testItem.mockReadInfo.expectedTocken}).
					Return(testItem.mockReadInfo.expectedSession, testItem.mockReadInfo.expectedErrorSM)
				if testItem.isNeedSetAnotherExpect {
					mockSM.EXPECT().CheckSession(&model.SessionID{ID: testItem.mockReadInfo.expectedTocken}).
						Return(testItem.mockReadInfo.expectedSession, testItem.mockReadInfo.expectedErrorSM)
				}
			}

			testModel := model.New(mockDB, mockSM)

			srv, ch := api.StartServer(api.Config{
				Address:      "localhost:54000",
				ReadTimeout:  "10s",
				WriteTimeout: "10s",
				IdleTimeout:  "10s",
			}, testModel)

			time.Sleep(time.Millisecond * 10) // time to start the server

			defer func() {
				srv.Shutdown(nil)
				<-ch
			}()

			client := http.DefaultClient

			var result *http.Response
			var err error
			if testItem.isOneAd || testItem.isOneUser ||
				testItem.isSeveralAds || testItem.isSeveralAdsOfUser {
				req, _ := http.NewRequest(
					"GET",
					"http://localhost:54000"+testItem.addr, nil)
				req.Header.Set("Content-Type",
					"application/x-www-form-urlencoded; charset=utf-8")

				if testItem.isUpdateUser {
					req.Header.Set("Cookie", "session_id="+testItem.mockReadInfo.expectedTocken)
				}

				result, err = client.Do(req)

			} else if testItem.isCreateUser || testItem.isLogin {
				req, _ := http.NewRequest(
					"POST",
					"http://localhost:54000"+testItem.addr,
					strings.NewReader(testItem.providedParams))
				req.Header.Set("Content-Type",
					"application/x-www-form-urlencoded; charset=utf-8")
				if testItem.isAnroid {
					req.Header.Set("User-Agent", "Android_app")
				}

				result, err = client.Do(req)
			} else if testItem.isLogout {
				req, _ := http.NewRequest(
					"POST",
					"http://localhost:54000"+testItem.addr,
					strings.NewReader(testItem.providedParams))
				req.Header.Set("Content-Type",
					"application/x-www-form-urlencoded; charset=utf-8")
				if testItem.providedParams != "" {
					req.Header.Set("Cookie", "session_id="+testItem.providedParams)
				}

				result, err = client.Do(req)
			} else if testItem.isUpdateUser {
				req, _ := http.NewRequest(
					"POST",
					"http://localhost:54000"+testItem.addr,
					strings.NewReader(testItem.providedParams))

				if testItem.isDeleteUser {
					req, _ = http.NewRequest(
						"DELETE",
						"http://localhost:54000"+testItem.addr,
						nil)
				}
				req.Header.Set("Content-Type",
					"application/x-www-form-urlencoded; charset=utf-8")

				req.Header.Set("Cookie", "session_id="+testItem.mockReadInfo.expectedTocken)
				result, err = client.Do(req)
			}

			if err != nil {
				t.Fatal("Expected no error while request\nGot: ", err.Error())
			}

			defer result.Body.Close()

			if result.StatusCode != testItem.expectedStatusCode {
				t.Errorf("Expected equal status codes:\nExpected:%v\nReceived:%v",
					testItem.expectedStatusCode, result.StatusCode)
			}

			if result.StatusCode == http.StatusOK && testItem.isOneAd {
				adDataExpected, _ := json.Marshal(testItem.expectedReturnAd)

				body, err := ioutil.ReadAll(result.Body)
				if err != nil {
					t.Error("Expected no error while reading body")
				}

				if string(body) != string(adDataExpected) {
					t.Errorf("Expected equal ads:\nExpected:%s\nReceived:%s",
						string(adDataExpected), string(body))
				}
			} else if result.StatusCode == http.StatusOK &&
				(testItem.isSeveralAds || testItem.isSeveralAdsOfUser) {
				adsDataExpected, _ := json.Marshal(testItem.expectedReturnAds)

				body, err := ioutil.ReadAll(result.Body)
				if err != nil {
					t.Error("Expected no error while reading body")
				}

				if string(body) != string(adsDataExpected) {
					t.Errorf("Expected equal ads:\nExpected:%s\nReceived:%s",
						string(adsDataExpected), string(body))
				}
			} else if result.StatusCode == http.StatusOK && testItem.isOneUser {
				userDataExpected, _ := json.Marshal(testItem.expectedReturnUser)

				body, err := ioutil.ReadAll(result.Body)
				if err != nil {
					t.Error("Expected no error while reading body")
				}

				if string(body) != string(userDataExpected) {
					t.Errorf("Expected equal users:\nExpected:%s\nReceived:%s",
						string(userDataExpected), string(body))
				}
			} else if result.StatusCode == http.StatusOK && testItem.isCreateUser {
				userDataExpected, _ := json.Marshal(testItem.expectedReturnUserCreate)

				body, err := ioutil.ReadAll(result.Body)
				if err != nil {
					t.Error("Expected no error while reading body")
				}

				if string(body) != string(userDataExpected) {
					t.Errorf("Expected equal users:\nExpected:%s\nReceived:%s",
						string(userDataExpected), string(body))
				}
			} else if result.StatusCode == http.StatusOK && testItem.isLogin {
				if testItem.isAnroid {
					userDataExpected, _ := json.Marshal(testItem.expectedReturnUserLogin)

					body, err := ioutil.ReadAll(result.Body)
					if err != nil {
						t.Error("Expected no error while reading body")
					}

					if string(body) != string(userDataExpected) {
						t.Errorf("Expected equal data:\nExpected:%s\nReceived:%s",
							string(userDataExpected), string(body))
					}
				}

				cookie := result.Cookies()
				if len(cookie) == 0 {
					t.Error("Expected set cookie on login")
				}

				if cookie[0].Value != testItem.mockReadInfo.expectedTocken {
					t.Errorf("Expected equal tockens:\nExpected:%s\nReceived:%s",
						testItem.mockReadInfo.expectedTocken, cookie[0].Value)
				}
			} else if result.StatusCode == http.StatusOK && testItem.isLogout {
				if testItem.providedParams != "" {
					cookie := result.Cookies()
					if len(cookie) == 0 {
						t.Error("Expected set cookie on login")
					}

					if cookie[0].Value != testItem.mockReadInfo.expectedTocken {
						t.Errorf("Expected equal tockens:\nExpected:%s\nReceived:%s",
							testItem.mockReadInfo.expectedTocken, cookie[0].Value)
					}
				}
			}
		}(testItem)
	}
}
