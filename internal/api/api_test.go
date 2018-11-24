package api_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
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
			Email:     "Ivan@ivanov.iva",
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
			Email:     "john@johnov.iva",
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
			Email:     "pet@animal.iva",
		},
		Description: "very good house",
	},
}

var usersInDB = map[int64]*model.User{
	12: &model.User{
		ID:        12,
		FirstName: "Ivan",
		LastName:  "Ivanov",
		Email:     "Ivan@ivanov.iva",
	},
	17: &model.User{
		ID:        17,
		FirstName: "Petya",
		LastName:  "Succeded",
		Email:     "pet@animal.iva",
	},
	18: &model.User{
		ID:        18,
		FirstName: "John",
		LastName:  "Johnov",
		Email:     "john@johnov.iva",
	},
}

type apiError struct {
	Description string `json:"description"`
	Message     string `json:"message"`
	ErrorCode   string `json:"error"`
}

type mockReadInfo struct {
	expectedAd    *model.AdItem
	expectedAds   []*model.AdItem
	expectedUser  *model.User
	expectedError error
}

type testCase struct {
	providedID     int64
	providedOffset int
	providedLimit  int

	isOneAd            bool
	isSeveralAds       bool
	isOneUser          bool
	isSeveralAdsOfUser bool

	expectedReturnAd    model.AdItem
	expectedReturnAds   []model.AdItem
	expectedReturnUser  model.User
	expectedReturnError apiError
	expectedStatusCode  int

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

func TestReadInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []testCase{
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
	}

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
			}

			mockSM := mock_model.NewMockSM(ctrl)

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

			client := http.Client{
				Timeout: time.Second * 5,
			}

			result, err := client.Get("http://localhost:54000" + testItem.addr)
			if err != nil {
				t.Fatal("Expected no error while request\nGot: ", err.Error())
			}

			defer result.Body.Close()

			if result.StatusCode != testItem.expectedStatusCode {
				t.Errorf("Expected equal status codes:\nExpected:%v\nReceived:%v", testItem.expectedStatusCode, result.StatusCode)
			}

			if result.StatusCode == http.StatusOK && testItem.isOneAd {
				adDataExpected, _ := json.Marshal(testItem.expectedReturnAd)

				body, err := ioutil.ReadAll(result.Body)
				if err != nil {
					t.Error("Expected no error while reading body")
				}

				if string(body) != string(adDataExpected) {
					t.Errorf("Expected equal ads:\nExpected:%s\nReceived:%s", string(adDataExpected), string(body))
				}
			} else if result.StatusCode == http.StatusOK && (testItem.isSeveralAds || testItem.isSeveralAdsOfUser) {
				adsDataExpected, _ := json.Marshal(testItem.expectedReturnAds)

				body, err := ioutil.ReadAll(result.Body)
				if err != nil {
					t.Error("Expected no error while reading body")
				}

				if string(body) != string(adsDataExpected) {
					t.Errorf("Expected equal ads:\nExpected:%s\nReceived:%s", string(adsDataExpected), string(body))
				}
			} else if result.StatusCode == http.StatusOK && testItem.isOneUser {
				userDataExpected, _ := json.Marshal(testItem.expectedReturnUser)

				body, err := ioutil.ReadAll(result.Body)
				if err != nil {
					t.Error("Expected no error while reading body")
				}

				if string(body) != string(userDataExpected) {
					t.Errorf("Expected equal users:\nExpected:%s\nReceived:%s", string(userDataExpected), string(body))
				}
			}
		}(testItem)
	}
}
