package starling

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"reflect"
	"testing"
	"time"
)

var txnTestCases = []struct {
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
	{
		name:      "without date range",
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
	{
		name:      "with multiple transactions",
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
							 "href": "api/v1/transactions/mastercard/6f03a23a-bbfc-4479-8d4d-abb6a9119d27",
							 "templated": false
							}
						},
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
	for _, tc := range txnTestCases {
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
		checkMethod(t, r, http.MethodGet)

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

	got, _, err := client.Transactions(context.Background(), dr)
	if err != nil {
		t.Fatal("\t\tshould be able to make the request", cross, err)
	} else {
		t.Log("\t\tshould be able to make the request", tick)
	}

	hal := &halTransactions{}
	json.Unmarshal([]byte(mock), hal)
	want := hal.Embedded

	if got == nil {
		t.Fatal("\t\tshould not return 'nil'", cross)
	} else {
		t.Log("\t\tshould not return 'nil'", tick)
	}

	if !reflect.DeepEqual(*got, want.Transactions) {
		t.Error("\t\tshould return a list matching the mock response", cross)
	} else {
		t.Log("\t\tshould return a transaction list matching the mock response", tick)
	}

	if len(*got) == 0 {
		t.Errorf("\t\tshould have at least one transaction %s %d", cross, len(*got))
	} else {
		t.Log("\t\tshould have at least one transaction", tick)
	}

}

var transactionCases = []struct {
	name string
	uid  string
	mock string
}{
	{
		name: "direct debit transaction",
		uid:  "474642e6-c4f5-43af-9b93-fe5ddbfcb857",
		mock: `{
			"_links": {
				 "detail": {
					  "href": "api/v1/transactions/direct-debit/474642e6-c4f5-43af-9b93-fe5ddbfcb857",
					  "templated": false
				 }
			},
			"id": "474642e6-c4f5-43af-9b93-fe5ddbfcb857",
			"currency": "GBP",
			"amount": -42.13,
			"direction": "OUTBOUND",
			"created": "2018-04-16T23:30:00.000Z",
			"narrative": "Society of Antiquaries",
			"source": "DIRECT_DEBIT"
	  }`,
	},
	{
		name: "direct debit transaction (without hal)",
		uid:  "474642e6-c4f5-43af-9b93-fe5ddbfcb857",
		mock: `{
			"id": "474642e6-c4f5-43af-9b93-fe5ddbfcb857",
			"currency": "GBP",
			"amount": -42.13,
			"direction": "OUTBOUND",
			"created": "2018-04-16T23:30:00.000Z",
			"narrative": "Society of Antiquaries",
			"source": "DIRECT_DEBIT"
	  }`,
	},
}

func TestTransaction(t *testing.T) {
	for _, tc := range transactionCases {
		t.Run(tc.name, func(st *testing.T) {
			testTransaction(st, tc.name, tc.uid, tc.mock)
		})
	}
}

func testTransaction(t *testing.T, name, uid, mock string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/transactions/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)

		reqUID := path.Base(r.URL.Path)
		if reqUID != uid {
			t.Error("should send a request with the correct UID", cross, reqUID)
		}

		fmt.Fprint(w, mock)
	})

	got, _, err := client.Transaction(context.Background(), uid)
	checkNoError(t, err)

	want := &Transaction{}
	json.Unmarshal([]byte(mock), want)

	if !reflect.DeepEqual(got, want) {
		t.Error("should return a single transaction matching the mock response", cross)
	}

	if got.UID != want.UID {
		t.Error("should have the correct UID", cross, got.UID)
	}
}

var txnsTestCasesDD = []struct {
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
			"_links": {
				 "nextPage": {
					  "href": "NOT_YET_IMPLEMENTED"
				 }
			},
			"_embedded": {
				 "transactions": [
					  {
							"id": "474642e6-c4f5-43af-9b93-fe5ddbfcb857",
							"currency": "GBP",
							"amount": -42.13,
							"direction": "OUTBOUND",
							"created": "2018-04-16T23:30:00.000Z",
							"narrative": "Society of Antiquaries",
							"source": "DIRECT_DEBIT",
							"mandate": {
								 "href": "/api/v1/direct-debit/mandates/fa7998f6-07ce-42a9-ba5b-ce45ea8aff89",
								 "templated": false
							},
							"merchant": {
								 "href": "/api/v1/merchants/b6c146f7-666e-4868-beed-21344b7e6e47",
								 "templated": false
							},
							"merchantLocation": {
								 "href": "/api/v1/merchants/b6c146f7-666e-4868-beed-21344b7e6e47/locations/7dda8396-7c7a-46d3-b5af-61a187bf00f9",
								 "templated": false
							},
							"mandateId": "fa7998f6-07ce-42a9-ba5b-ce45ea8aff89",
							"type": "FIRST_PAYMENT_OF_DIRECT_DEBIT",
							"merchantId": "b6c146f7-666e-4868-beed-21344b7e6e47",
							"merchantLocationId": "7dda8396-7c7a-46d3-b5af-61a187bf00f9",
							"spendingCategory": "GENERAL"
					  }
				 ]
			}
	  }`,
	},
	{
		name:      "without date range",
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
							"id": "474642e6-c4f5-43af-9b93-fe5ddbfcb857",
							"currency": "GBP",
							"amount": -42.13,
							"direction": "OUTBOUND",
							"created": "2018-04-16T23:30:00.000Z",
							"narrative": "Society of Antiquaries",
							"source": "DIRECT_DEBIT",
							"mandate": {
								 "href": "/api/v1/direct-debit/mandates/fa7998f6-07ce-42a9-ba5b-ce45ea8aff89",
								 "templated": false
							},
							"merchant": {
								 "href": "/api/v1/merchants/b6c146f7-666e-4868-beed-21344b7e6e47",
								 "templated": false
							},
							"merchantLocation": {
								 "href": "/api/v1/merchants/b6c146f7-666e-4868-beed-21344b7e6e47/locations/7dda8396-7c7a-46d3-b5af-61a187bf00f9",
								 "templated": false
							},
							"mandateId": "fa7998f6-07ce-42a9-ba5b-ce45ea8aff89",
							"type": "FIRST_PAYMENT_OF_DIRECT_DEBIT",
							"merchantId": "b6c146f7-666e-4868-beed-21344b7e6e47",
							"merchantLocationId": "7dda8396-7c7a-46d3-b5af-61a187bf00f9",
							"spendingCategory": "GENERAL"
					  }
				 ]
			}
	  }`,
	},
}

func TestDDTransactions(t *testing.T) {
	for _, tc := range txnsTestCasesDD {
		t.Run(tc.name, func(t *testing.T) {
			testDDTransactions(t, tc.name, tc.mock, tc.dateRange)
		})
	}
}

func testDDTransactions(t *testing.T, name, mock string, dr *DateRange) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/transactions/direct-debit", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)

		params := r.URL.Query()

		if dr != nil {
			if got, want := params.Get("from"), dr.From.Format("2006-01-02"); got != want {
				t.Errorf("should include 'from=%s' query string parameter %s 'from=%s'", want, cross, got)
			}

			if got, want := params.Get("to"), dr.To.Format("2006-01-02"); got != want {
				t.Errorf("should include 'to=%s' query string parameter %s 'to=%s'", want, cross, got)
			}
		} else {
			if got, want := params.Get("from"), ""; got != want {
				t.Errorf("should not include 'from' query string parameter %s 'from=%s'", cross, got)
			}

			if got, want := params.Get("to"), ""; got != want {
				t.Errorf("should not include 'to' query string parameter %s 'to=%s'", cross, got)
			}
		}

		fmt.Fprint(w, mock)
	})

	got, _, err := client.DDTransactions(context.Background(), dr)
	if err != nil {
		t.Fatal("should be able to make the request", cross, err)
	}

	hal := &halDDTransactions{}
	json.Unmarshal([]byte(mock), hal)
	want := hal.Embedded

	if got == nil {
		t.Fatal("should not return 'nil'", cross)
	}

	if !reflect.DeepEqual(*got, want.Transactions) {
		t.Error("should return a list matching the mock response", cross)
	}

	if len(*got) == 0 {
		t.Errorf("should have at least one transaction %s %d", cross, len(*got))
	}

	first := (*got)[0]

	if first.MandateUID == "" {
		t.Error("should have a MandateID specified", cross)
	}

	if first.Type == "" {
		t.Error("should have a Type specified", cross)
	}

	if first.MerchantUID == "" {
		t.Error("should have a MerchantUID specified", cross)
	}

	if first.MerchantLocationUID == "" {
		t.Error("should have a MerchantLocationUID specified", cross)
	}

	if first.SpendingCategory == "" {
		t.Error("should have a SpendingCategory specified", cross)
	}

}

var txnTestCasesDD = []struct {
	name string
	uid  string
	mock string
}{
	{
		name: "single direct-debit transaction",
		uid:  "474642e6-c4f5-43af-9b93-fe5ddbfcb857",
		mock: `{
			"id": "474642e6-c4f5-43af-9b93-fe5ddbfcb857",
			"currency": "GBP",
			"amount": -42.13,
			"direction": "OUTBOUND",
			"created": "2018-04-16T23:30:00.000Z",
			"narrative": "Society of Antiquaries",
			"source": "DIRECT_DEBIT",
			"mandate": {
				 "href": "/api/v1/direct-debit/mandates/fa7998f6-07ce-42a9-ba5b-ce45ea8aff89",
				 "templated": false
			},
			"merchant": {
				 "href": "/api/v1/merchants/b6c146f7-666e-4868-beed-21344b7e6e47",
				 "templated": false
			},
			"merchantLocation": {
				 "href": "/api/v1/merchants/b6c146f7-666e-4868-beed-21344b7e6e47/locations/7dda8396-7c7a-46d3-b5af-61a187bf00f9",
				 "templated": false
			},
			"mandateId": "fa7998f6-07ce-42a9-ba5b-ce45ea8aff89",
			"type": "FIRST_PAYMENT_OF_DIRECT_DEBIT",
			"merchantId": "b6c146f7-666e-4868-beed-21344b7e6e47",
			"merchantLocationId": "7dda8396-7c7a-46d3-b5af-61a187bf00f9",
			"spendingCategory": "GENERAL"
	  }`,
	},
}

func TestDDTransaction(t *testing.T) {
	for _, tc := range txnTestCasesDD {
		t.Run(tc.name, func(t *testing.T) {
			testDDTransaction(t, tc.name, tc.uid, tc.mock)
		})
	}
}

func testDDTransaction(t *testing.T, name, uid, mock string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/transactions/direct-debit/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)

		reqUID := path.Base(r.URL.Path)
		if reqUID != uid {
			t.Error("should send a request with the correct UID", cross, reqUID)
		}

		fmt.Fprint(w, mock)
	})

	want := &DDTransaction{}
	json.Unmarshal([]byte(mock), want)

	got, _, err := client.DDTransaction(context.Background(), uid)
	if err != nil {
		t.Fatal("should be able to make the request", cross, err)
	}

	if got == nil {
		t.Fatal("should not return 'nil'", cross)
	}

	if !reflect.DeepEqual(got, want) {
		t.Error("should return a transaction matching the mock response", cross)
	}

	if got.MandateUID == "" {
		t.Error("should have a MandateID specified", cross)
	}

	if got.Type == "" {
		t.Error("should have a Type specified", cross)
	}

	if got.MerchantUID == "" {
		t.Error("should have a MerchantUID specified", cross)
	}

	if got.MerchantLocationUID == "" {
		t.Error("should have a MerchantLocationUID specified", cross)
	}

	if got.SpendingCategory == "" {
		t.Error("should have a SpendingCategory specified", cross)
	}

}
