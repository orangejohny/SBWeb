package apitests

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"

	"bmstu.codes/developers34/SBWeb/internal/api"
	"bmstu.codes/developers34/SBWeb/internal/model"
)

func TestReadOneAd(t *testing.T) {
	testSM := smTestCase{
		sessionID: model.SessionID{ID: "dwedqf1234ewdw"},
		err:       nil,
	}

	type testCase struct {
		dbTestCase
		smTestCase
	}

	testCases := []testCase{
		testCase{
			dbTestCase: dbTestCase{
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
				id:  15,
			},
			smTestCase: testSM,
		},
	}

	for _, testItem := range testCases {
		go api.StartServer(api.Config{
			Address:      "localhost:54000",
			ReadTimeout:  "10s",
			WriteTimeout: "10s",
			IdleTimeout:  "10s",
		}, model.New(&testItem.dbTestCase, &testItem.smTestCase))

		client := http.Client{
			Timeout: time.Second * 5,
		}

		result, err := client.Get("http://localhost:54000/ads/" + strconv.FormatInt(testItem.dbTestCase.id, 10))
		if err != nil {
			t.Error("Expected no error while request\nGot: ", err.Error())
		}

		adData, _ := json.Marshal(testItem.dbTestCase.ad)

		defer result.Body.Close()
		body, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Error("Expected no error while reading body")
		}

		if string(body) != string(adData) {
			t.Errorf("Expected equal ads:\nExpected:%s\nReceived:%s", string(adData), string(body))
		}
	}
}
