package starling

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var customersTestCases = []struct {
	name string
	mock string
}{
	{
		name: "sample customer",
		mock: `{
			"_links": {
				"accounts": {
					"href": "api/v1/accounts",
					"templated": false
				},
				"addresses": {
					"href": "api/v1/addresses",
					"templated": false
				},
				"contacts": {
					"href": "api/v1/contacts",
					"templated": false
				},
				"self": {
					"href": "api/v1/customers",
					"templated": false
				},
				"transactions": {
					"href": "api/v1/transactions?from={fromDate}&to={toDate}",
					"templated": false
				}
			},
			"customerUid": "6d2aa528-b9d1-4083-ae7c-53d460cd8d88",
			"firstName": "Vincent",
			"lastName": "Adultman",
			"dateOfBirth": "1960-01-01",
			"email": "6d2aa528-b9d1-4083-ae7c-53d460cd8d88@starlingbank.com",
			"phone": "+447886725871",
			"accountHolderType": "INDIVIDUAL"
		}`,
	},
	{
		name: "sample customer without HAL wrapper",
		mock: `{
			"customerUid": "6d2aa528-b9d1-4083-ae7c-53d460cd8d88",
			"firstName": "Vincent",
			"lastName": "Adultman",
			"dateOfBirth": "1960-01-01",
			"email": "6d2aa528-b9d1-4083-ae7c-53d460cd8d88@starlingbank.com",
			"phone": "+447886725871",
			"accountHolderType": "INDIVIDUAL"
		}`,
	},
}

func TestCustomer(t *testing.T) {
	for _, tc := range customersTestCases {
		t.Run(tc.name, func(st *testing.T) {
			testCustomer(st, tc.name, tc.mock)
		})
	}
}

func testCustomer(t *testing.T, name, mock string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/customers", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		fmt.Fprint(w, mock)
	})

	got, _, err := client.Customer(context.Background())
	checkNoError(t, err)

	want := &Customer{}
	json.Unmarshal([]byte(mock), want)

	if !reflect.DeepEqual(got, want) {
		t.Error("should return an identity matching the mock response", cross)
	}
}

func TestCustomerForbidden(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/customers", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusForbidden)
	})

	got, resp, err := client.Customer(context.Background())
	checkHasError(t, err)

	if resp.StatusCode != http.StatusForbidden {
		t.Error("should return HTTP 403 status")
	}

	if got != nil {
		t.Error("should not return a customer")
	}
}
