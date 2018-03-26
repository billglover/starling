package starling

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

const (
	tick  = "\u2713"
	cross = "\u2717"
)

var testCases = []struct {
	name      string
	mock      string
	dateRange *DateRange
}{
	{
		name: "with date range",
		dateRange: &DateRange{
			From: time.Date(2017, time.January, 01, 0, 0, 0, 0, time.Local),
			To:   time.Date(2017, time.January, 03, 0, 0, 0, 0, time.Local),
		},
		mock: `{
			"transactions": [
			  {
				"id": "6f03a23a-bbfc-4479-8d4d-abb6a9119d27",
				"currency": "GBP",
				"amount": -23.45,
				"direction": "OUTBOUND",
				"created": "2017-07-05T18:27:02.335Z",
				"narrative": "Borough Barista",
				"source": "MASTER_CARD",
				"balance": 254.12
			  }
			]
		  }`,
	},
	{
		name:      "without date range",
		dateRange: nil,
		mock: `{
			"transactions": [
			  {
				"id": "6f03a23a-bbfc-4479-8d4d-abb6a9119d27",
				"currency": "GBP",
				"amount": -23.45,
				"direction": "OUTBOUND",
				"created": "2017-07-05T18:27:02.335Z",
				"narrative": "Borough Barista",
				"source": "MASTER_CARD",
				"balance": 254.12
			  }
			]
		  }`,
	},
	{
		name:      "with multiple transactions",
		dateRange: nil,
		mock: `{
			"transactions": [
			  {
				"id": "6f03a23a-bbfc-4479-8d4d-abb6a9119d27",
				"currency": "GBP",
				"amount": -23.45,
				"direction": "OUTBOUND",
				"created": "2017-07-05T18:27:02.335Z",
				"narrative": "Borough Barista",
				"source": "MASTER_CARD",
				"balance": 254.12
			  },
			  {
				"id": "6f03a23a-bbfc-4479-8d4d-abb6a9119d27",
				"currency": "GBP",
				"amount": -23.45,
				"direction": "OUTBOUND",
				"created": "2017-07-05T18:27:02.335Z",
				"narrative": "Borough Barista",
				"source": "MASTER_CARD",
				"balance": 254.12
			  }
			]
		  }`,
	},
	{
		name:      "with HAL wrapper",
		dateRange: nil,
		mock: `{
			"_links": {
			  "nextPage": {
				 "href": "NOT_YET_IMPLEMENTED"
			  }
			},
			"_embedded": {
			  "transactions": [
				 {
					"_links": {
					  "detail": {
						 "href": "api/v1/transactions/mastercard/0e70192c-e602-40ac-b306-c21630e6874e",
						 "templated": false
					  }
					},
					"id": "0e70192c-e602-40ac-b306-c21630e6874e",
					"currency": "GBP",
					"amount": -13.99,
					"direction": "OUTBOUND",
					"created": "2018-03-25T11:55:26.865Z",
					"narrative": "Mastercard",
					"source": "MASTER_CARD",
					"balance": 13081.32
				 }
			  ]
			}
		 }`,
	},
}

func TestGetTransactions(t *testing.T) {

	t.Log("Given the need to test fetching transactions:")

	// Run each of the test cases a subtest.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testGetTransactions(t, tc.name, tc.mock, tc.dateRange)
		})
	}
}

func testGetTransactions(t *testing.T, name, mock string, dr *DateRange) {
	t.Logf("\tWhen making a call to GetTransactions() %s:", name)

	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/transactions", func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, "GET"; got != want {
			t.Errorf("\t\tshould send a %s request to the API %s %s", want, cross, got)
		} else {
			t.Logf("\t\tshould send a %s request to the API %s", want, tick)
		}

		params := r.URL.Query()

		if dr != nil {

			// If we are expecting a date range to be passed to the API, validate that it was
			// passed correctly as part of the query string.
			if got, want := params.Get("from"), dr.From.Format("2006-01-02"); got != want {
				t.Errorf("\t\tshould include 'from=%s' query string parameter %s 'from=%s'", want, cross, got)
			} else {
				t.Logf("\t\tshould include 'from=%s' query string parameter %s", want, tick)
			}

			if got, want := params.Get("to"), dr.To.Format("2006-01-02"); got != want {
				t.Errorf("\t\tshould include 'to=%s' query string parameter %s 'to=%s'", want, cross, got)
			} else {
				t.Logf("\t\tshould include 'to=%s' query string parameter %s", want, tick)
			}
		} else {

			// If we weren't expecting a date range to be passed to the API, validate that the
			// API was called without the 'from' and 'to' query parameters.
			if got, want := params.Get("from"), ""; got != want {
				t.Errorf("\t\tshould not include 'from' query string parameter %s 'from=%s'", cross, got)
			} else {
				t.Logf("\t\tshould not include 'from' query string parameter %s", tick)
			}

			if got, want := params.Get("to"), ""; got != want {
				t.Errorf("\t\tshould not include 'to' query string parameter %s 'to=%s'", cross, got)
			} else {
				t.Logf("\t\tshould not include 'to' query string parameter %s", tick)
			}
		}

		// Return the mock response to the client.
		fmt.Fprint(w, mock)
	})

	got, _, err := client.GetTransactions(context.Background(), dr)
	if err != nil {
		t.Fatal("\t\tshould be able to make the request", cross, err)
	} else {
		t.Log("\t\tshould be able to make the request", tick)
	}

	want := &Transactions{}
	json.Unmarshal([]byte(mock), want)

	if !reflect.DeepEqual(got, want) {
		t.Error("\t\tshould return a list matching the mock response", cross)
	} else {
		t.Log("\t\tshould return a transaction list matching the mock response", tick)
	}

	if len(got.Transactions) == 0 {
		t.Errorf("\t\tshould have at least one transaction %s %d", cross, len(got.Transactions))
	} else {
		t.Log("\t\tshould have at least one transaction", tick)
	}

}
