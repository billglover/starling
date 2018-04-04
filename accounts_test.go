package starling

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var accountTestCases = []struct {
	name string
	mock string
}{
	{
		name: "sample account",
		mock: `{
			"_links": {
				"card": {
					"href": "api/v1/cards",
					"templated": false
				},
				"customer": {
					"href": "api/v1/customers",
					"templated": false
				},
				"mandates": {
					"href": "/api/v1/direct-debit/mandates",
					"templated": false
				},
				"payees": {
					"href": "api/v1/contacts",
					"templated": false
				},
				"transactions": {
					"href": "api/v1/transactions?from={fromDate}&to={toDate}",
					"templated": true
				}
			},
			"id": "6f5a3548-f25d-4dfe-9f8e-3078fe8bfa2a",
			"name": "6d2aa528-b9d1-4083-ae7c-53d460cd8d88 GBP",
			"number": "04829435",
			"accountNumber": "04829435",
			"sortCode": "608371",
			"currency": "GBP",
			"iban": "GB28SRLG60837104829435",
			"bic": "SRLGGB2L",
			"createdAt": "2017-03-09T17:58:15.848Z"
		}`,
	},
	{
		name: "sample account without HAL wrapper",
		mock: `{
			"id": "6f5a3548-f25d-4dfe-9f8e-3078fe8bfa2a",
			"name": "6d2aa528-b9d1-4083-ae7c-53d460cd8d88 GBP",
			"number": "04829435",
			"accountNumber": "04829435",
			"sortCode": "608371",
			"currency": "GBP",
			"iban": "GB28SRLG60837104829435",
			"bic": "SRLGGB2L",
			"createdAt": "2017-03-09T17:58:15.848Z"
		}`,
	},
}

func TestGetAccount(t *testing.T) {

	t.Log("Given the need to test fetching account details:")

	for _, tc := range accountTestCases {
		t.Run(tc.name, func(st *testing.T) {
			testGetAccount(st, tc.name, tc.mock)
		})
	}
}

func testGetAccount(t *testing.T, name, mock string) {
	t.Logf("\tWhen making a call to GetAccount() with %s:", name)

	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/accounts", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		fmt.Fprint(w, mock)
	})

	got, _, err := client.GetAccount(context.Background())
	checkNoError(t, err)

	want := &Account{}
	json.Unmarshal([]byte(mock), want)

	if !reflect.DeepEqual(got, want) {
		t.Error("\t\tshould return an account matching the mock response", cross)
	} else {
		t.Log("\t\tshould return an account matching the mock response", tick)
	}
}

var balanceTestCases = []struct {
	name string
	mock string
}{
	{
		name: "positive balance",
		mock: `{
			"clearedBalance": 15260.82,
			"effectiveBalance": 15260.82,
			"pendingTransactions": 0,
			"availableToSpend": 15260.82,
			"currency": "GBP",
			"amount": 15260.82
		}`,
	},
	{
		name: "negative balance",
		mock: `{
			"clearedBalance": -15260.82,
			"effectiveBalance": -15260.82,
			"pendingTransactions": 0,
			"availableToSpend": 0,
			"currency": "GBP",
			"amount": -15260.82
		}`,
	},
	{
		name: "very large balance",
		mock: `{
			"clearedBalance": -15260.82,
			"effectiveBalance": -15260.82,
			"pendingTransactions": 0,
			"availableToSpend": 0,
			"currency": "GBP",
			"amount": 1.797693134862315708145274237317043567981e+308
		}`,
	},
}

func TestGetAccountBalance(t *testing.T) {

	t.Log("Given the need to test fetching account balance:")

	for _, tc := range balanceTestCases {
		t.Run(tc.name, func(st *testing.T) {
			testGetAccountBalance(st, tc.name, tc.mock)
		})
	}
}

func testGetAccountBalance(t *testing.T, name, mock string) {
	t.Logf("\tWhen making a call to GetAccountBalance() with a %s:", name)

	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/accounts/balance", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		fmt.Fprint(w, mock)
	})

	got, _, err := client.GetAccountBalance(context.Background())
	checkNoError(t, err)

	want := &Balance{}
	json.Unmarshal([]byte(mock), want)

	if !reflect.DeepEqual(got, want) {
		t.Error("\t\tshould return an account balance matching the mock response", cross)
	} else {
		t.Log("\t\tshould return an account balancne matching the mock response", tick)
	}
}