package starling

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var cardTestCases = []struct {
	name string
	mock string
}{
	{
		name: "sample card",
		mock: `{
			"_links": {
				"transactions": {
					"href": "/api/v1/transactions/mastercard?from={fromDate}&to={toDate}",
					"templated": true
				}
			},
			"id": "8e9c955c-b209-4887-af32-a9e4999e985e",
			"nameOnCard": "Vincent Adultman",
			"type": "ContactlessDebitMastercard",
			"enabled": true,
			"cancelled": false,
			"activationRequested": true,
			"activated": true,
			"dispatchDate": "2018-03-13",
			"lastFourDigits": "0142"
		}`,
	},
	{
		name: "sample card without HAL wrapper",
		mock: `{
			"id": "8e9c955c-b209-4887-af32-a9e4999e985e",
			"nameOnCard": "Vincent Adultman",
			"type": "ContactlessDebitMastercard",
			"enabled": true,
			"cancelled": false,
			"activationRequested": true,
			"activated": true,
			"dispatchDate": "2018-03-13",
			"lastFourDigits": "0142"
		}`,
	},
}

func TestCard(t *testing.T) {
	for _, tc := range cardTestCases {
		t.Run(tc.name, func(st *testing.T) {
			testCard(st, tc.name, tc.mock)
		})
	}
}

func testCard(t *testing.T, name, mock string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/cards", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		fmt.Fprint(w, mock)
	})

	got, _, err := client.Card(context.Background())
	checkNoError(t, err)

	want := &Card{}
	json.Unmarshal([]byte(mock), want)

	if !reflect.DeepEqual(got, want) {
		t.Error("should return a card matching the mock response", cross)
	}
}
