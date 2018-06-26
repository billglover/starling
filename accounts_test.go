package starling

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var accountsTC = []struct {
	name string
	mock string
}{
	{
		name: "single account",
		mock: `{
			"accounts": [
				{
					"accountUid": "24492cc9-77dd-4155-87a2-ec2580daf139",
					"defaultCategory": "8d8c0f3b-f685-49ed-835e-db2ff8cef703",
					"currency": "GBP",
					"createdAt": "2017-05-24T07:43:46.664Z"
				}
			]
	  }`,
	},
	{
		name: "two accounts",
		mock: `{
			"accounts": [
				{
					"accountUid": "24492cc9-77dd-4155-87a2-ec2580daf139",
					"defaultCategory": "8d8c0f3b-f685-49ed-835e-db2ff8cef703",
					"currency": "GBP",
					"createdAt": "2017-05-24T07:43:46.664Z"
				},
				{
					"accountUid": "24492cc9-77dd-4155-87a2-ec2580daf139",
					"defaultCategory": "8d8c0f3b-f685-49ed-835e-db2ff8cef703",
					"currency": "GBP",
					"createdAt": "2017-05-24T07:43:46.664Z"
				}
			]
	  }`,
	},
}

func TestAccount(t *testing.T) {
	for _, tc := range accountsTC {
		t.Run(tc.name, func(st *testing.T) {
			testAccount(st, tc.name, tc.mock)
		})
	}
}

func testAccount(t *testing.T, name, mock string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/accounts", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		fmt.Fprint(w, mock)
	})

	got, _, err := client.Account(context.Background())
	checkNoError(t, err)

	want := &Account{}
	json.Unmarshal([]byte(mock), want)

	if !reflect.DeepEqual(got, want) {
		t.Error("should return an account matching the mock response", cross)
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

func TestAccountBalance(t *testing.T) {
	for _, tc := range balanceTestCases {
		t.Run(tc.name, func(st *testing.T) {
			testAccountBalance(st, tc.name, tc.mock)
		})
	}
}

func testAccountBalance(t *testing.T, name, mock string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/accounts/balance", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		fmt.Fprint(w, mock)
	})

	got, _, err := client.AccountBalance(context.Background())
	checkNoError(t, err)

	want := &Balance{}
	json.Unmarshal([]byte(mock), want)

	if !reflect.DeepEqual(got, want) {
		t.Error("should return an account balance matching the mock response", cross)
	}
}

func TestAccountForbidden(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/accounts", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusForbidden)
	})

	got, resp, err := client.Account(context.Background())
	checkHasError(t, err)

	if resp.StatusCode != http.StatusForbidden {
		t.Error("should return HTTP 403 status")
	}

	if got != nil {
		t.Error("should not return an account")
	}
}

func TestBalanceForbidden(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/accounts/balance", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusForbidden)
	})

	got, resp, err := client.AccountBalance(context.Background())
	checkHasError(t, err)

	if resp.StatusCode != http.StatusForbidden {
		t.Error("should return HTTP 403 status")
	}

	if got != nil {
		t.Error("should not return an account")
	}
}
