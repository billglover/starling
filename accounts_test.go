package starling

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"reflect"
	"testing"
)

var accountsTC = []struct {
	name  string
	count int
	mock  string
}{
	{
		name:  "single account",
		count: 1,
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
		name:  "two personal accounts",
		count: 2,
		mock: `{
					"accounts": [
				 		{
							"accountUid": "24492cc9-77dd-4155-87a2-ec2580daf139",
							"defaultCategory": "8d8c0f3b-f685-49ed-835e-db2ff8cef703",
							"currency": "GBP",
							"createdAt": "2017-05-24T07:43:46.664Z"
						},
						{
							"accountUid": "654BB6AB-3C10-49C2-9D4E-D49968772BB0",
							"defaultCategory": "09e7e421-1afc-483a-98be-0b9da90f9a57",
							"currency": "GBP",
							"createdAt": "2017-05-24T07:43:46.664Z"
						}
					]
				}`,
	},
}

func TestAccounts(t *testing.T) {
	for _, tc := range accountsTC {
		t.Run(tc.name, func(st *testing.T) {
			testAccounts(st, tc.name, tc.count, tc.mock)
		})
	}
}

func testAccounts(t *testing.T, name string, count int, mock string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		fmt.Fprint(w, mock)
	})

	got, _, err := client.Accounts(context.Background())
	checkNoError(t, err)

	if len(got) != count {
		t.Error("should return the correct number of accounts", got)
	}

	want := &accounts{}
	json.Unmarshal([]byte(mock), want)

	if !reflect.DeepEqual(got, want.Accounts) {
		t.Error("should return an account matching the mock response")
		t.Error(got)
		t.Error(want)
	}
}

func TestAccountsForbidden(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusForbidden)
	})

	got, resp, err := client.Accounts(context.Background())
	checkHasError(t, err)

	if resp.StatusCode != http.StatusForbidden {
		t.Error("should return HTTP 403 status")
	}

	if got != nil {
		t.Error("should not return an account")
	}
}

var accountIDsTC = []struct {
	name string
	uid  string
	mock string
}{
	{
		name: "personal account",
		uid:  "2c7a379d-c0d8-4541-8520-ca41cc26d56a",
		mock: `{
			"accountIdentifier": "12345678",
			"bankIdentifier": "608371",
			"iban": "GB50SRLG60837112345678",
			"bic": "SRLGGB2L"
		 }`,
	},
}

func TestAccountIDs(t *testing.T) {
	for _, tc := range accountIDsTC {
		t.Run(tc.name, func(st *testing.T) {
			testAccountIDs(st, tc.name, tc.uid, tc.mock)
		})
	}
}

func testAccountIDs(t *testing.T, name, uid, mock string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v2/accounts/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)

		reqActUID := path.Base(path.Dir(r.URL.Path))
		if reqActUID != uid {
			t.Error("should send a request with the correct account UID")
		}

		resource := path.Base(r.URL.Path)
		if resource != "identifiers" {
			t.Error("should send a request for the identifiers resource")
		}

		fmt.Fprint(w, mock)
	})

	got, _, err := client.AccountID(context.Background(), uid)
	checkNoError(t, err)

	want := new(AccountID)
	json.Unmarshal([]byte(mock), want)

	if !reflect.DeepEqual(got, want) {
		t.Error("should return an AccountID matching the mock response")
	}
}

func TestAccountIDsForbidden(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v2/accounts/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusForbidden)
	})

	got, resp, err := client.AccountID(context.Background(), "2c7a379d-c0d8-4541-8520-ca41cc26d56a")
	checkHasError(t, err)

	if resp.StatusCode != http.StatusForbidden {
		t.Error("should return HTTP 403 status")
	}

	if got != nil {
		t.Error("should not return an account")
	}
}

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

func TestAccount(t *testing.T) {
	for _, tc := range accountTestCases {
		t.Run(tc.name, func(st *testing.T) {
			testAccount(st, tc.name, tc.mock)
		})
	}
}

func testAccount(t *testing.T, name, mock string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
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

	mux.HandleFunc("/api/v2/accounts", func(w http.ResponseWriter, r *http.Request) {
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
