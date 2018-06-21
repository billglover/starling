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
	for _, tc := range txnTestCases {
		t.Run(tc.name, func(t *testing.T) {
			testGetTransactions(t, tc.name, tc.mock, tc.dateRange)
		})
	}
}

func testGetTransactions(t *testing.T, name, mock string, dr *DateRange) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/transactions", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)

		params := r.URL.Query()

		if dr != nil {

			// If we are expecting a date range to be passed to the API, validate that it was
			// passed correctly as part of the query string.
			if got, want := params.Get("from"), dr.From.Format("2006-01-02"); got != want {
				t.Errorf("should include 'from=%s' query string parameter %s 'from=%s'", want, cross, got)
			}

			if got, want := params.Get("to"), dr.To.Format("2006-01-02"); got != want {
				t.Errorf("should include 'to=%s' query string parameter %s 'to=%s'", want, cross, got)
			}
		} else {

			// If we weren't expecting a date range to be passed to the API, validate that the
			// API was called without the 'from' and 'to' query parameters.
			if got, want := params.Get("from"), ""; got != want {
				t.Errorf("should not include 'from' query string parameter %s 'from=%s'", cross, got)
			}

			if got, want := params.Get("to"), ""; got != want {
				t.Errorf("should not include 'to' query string parameter %s 'to=%s'", cross, got)
			}
		}

		// Return the mock response to the client.
		fmt.Fprint(w, mock)
	})

	got, _, err := client.Transactions(context.Background(), dr)
	if err != nil {
		t.Fatal("should be able to make the request", cross, err)
	}

	hal := &halTransactions{}
	json.Unmarshal([]byte(mock), hal)
	want := hal.Embedded

	if got == nil {
		t.Fatal("should not return 'nil'", cross)
	}

	if !reflect.DeepEqual(got, want.Transactions) {
		t.Error("should return a list matching the mock response", cross)
	}

	if len(got) == 0 {
		t.Errorf("should have at least one transaction %s %d", cross, len(got))
	}

}

func TestTxnsForbidden(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/transactions/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusForbidden)
	})

	got, resp, err := client.Transactions(context.Background(), nil)
	checkHasError(t, err)

	if resp.StatusCode != http.StatusForbidden {
		t.Error("should return HTTP 403 status")
	}

	if got != nil {
		t.Error("should not return a slice of transactions")
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

func TestTxnForbidden(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/transactions/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusForbidden)
	})

	got, resp, err := client.Transaction(context.Background(), "nil")
	checkHasError(t, err)

	if resp.StatusCode != http.StatusForbidden {
		t.Error("should return HTTP 403 status")
	}

	if got != nil {
		t.Error("should not return a transaction")
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

	if !reflect.DeepEqual(got, want.Transactions) {
		t.Error("should return a list matching the mock response", cross)
	}

	if len(got) == 0 {
		t.Errorf("should have at least one transaction %s %d", cross, len(got))
	}

	first := (got)[0]

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

func TestTxnsDDForbidden(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/transactions/direct-debit", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusForbidden)
	})

	got, resp, err := client.DDTransactions(context.Background(), nil)
	checkHasError(t, err)

	if resp.StatusCode != http.StatusForbidden {
		t.Error("should return HTTP 403 status")
	}

	if got != nil {
		t.Error("should not return a slice of transactions")
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

func TestDDTxnForbidden(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/transactions/direct-debit/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusForbidden)
	})

	got, resp, err := client.DDTransaction(context.Background(), "nil")
	checkHasError(t, err)

	if resp.StatusCode != http.StatusForbidden {
		t.Error("should return HTTP 403 status")
	}

	if got != nil {
		t.Error("should not return a transaction")
	}
}

var setDDSpendingCategoryCases = []struct {
	name     string
	uid      string
	category string
	status   int
	mock     string
}{
	{
		name:     "set transaction category to charity",
		uid:      "474642e6-c4f5-43af-9b93-fe5ddbfcb857",
		category: "CHARITY",
		status:   http.StatusAccepted,
		mock:     ``,
	},
	{
		name:     "set invalid transaction category",
		uid:      "474642e6-c4f5-43af-9b93-fe5ddbfcb857",
		category: "INVALID",
		status:   http.StatusBadRequest,
		mock: `[
			"Can not deserialize value of type com.starlingbank.connectors.customer.SpendingCategory from String \"DEMOX\": value not one of declared Enum instance names: [GIFTS, FAMILY, ENTERTAINMENT, TRANSPORT, GROCERIES, PAYMENTS, PETS, LIFESTYLE, CHARITY, BILLS_AND_SERVICES, SAVING, HOLIDAYS, HOME, GENERAL, NONE, EXPENSES, INCOME, SHOPPING, EATING_OUT]"
	  ]`,
	},
}

func TestSetDDSpendingCategory(t *testing.T) {
	for _, tc := range setDDSpendingCategoryCases {
		t.Run(tc.name, func(st *testing.T) {
			testSetDDSpendingCategory(st, tc.name, tc.uid, tc.category, tc.status, tc.mock)
		})
	}
}

func testSetDDSpendingCategory(t *testing.T, name, uid, cat string, status int, mock string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/transactions/direct-debit/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodPut)

		reqUID := path.Base(r.URL.Path)
		resource := path.Base(path.Dir(r.URL.Path))

		if reqUID != uid {
			t.Error("should send a request with the correct UID", cross, reqUID)
		}

		if resource != "direct-debit" {
			t.Error("should request direct-debit", cross, resource)
		}

		var reqCat = SpendingCategory{}
		err := json.NewDecoder(r.Body).Decode(&reqCat)
		if err != nil {
			t.Fatal("should send a request that the API can parse", cross, err)
		}

		if reqCat.SpendingCategory != cat {
			t.Error("should request the correct spending category", cross, reqCat)
		}

		w.WriteHeader(status)
		fmt.Fprintln(w, mock)
	})

	resp, err := client.SetDDSpendingCategory(context.Background(), uid, cat)
	if status >= 400 {
		checkHasError(t, err)
	} else {
		checkNoError(t, err)
	}

	if resp.StatusCode != status {
		t.Error("should return the correct status", cross, resp.Status)
	}
}

func TestDDSpendingCategoryForbidden(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/transactions/direct-debit/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodPut)
		w.WriteHeader(http.StatusForbidden)
	})

	resp, err := client.SetDDSpendingCategory(context.Background(), "474642e6-c4f5-43af-9b93-fe5ddbfcb857", "CHARITY")
	checkHasError(t, err)

	if resp.StatusCode != http.StatusForbidden {
		t.Error("should return HTTP 403 status")
	}
}

var txnsTestCasesFPSIn = []struct {
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
							"id": "4f39ce4a-b760-42d8-811d-792e366486ef",
							"currency": "GBP",
							"amount": 33.14,
							"direction": "INBOUND",
							"created": "2018-03-28T13:48:56.000Z",
							"narrative": "someRef",
							"source": "FASTER_PAYMENTS_IN"
					  },
					  {
							"id": "e517d335-2fb8-486b-91b3-2762ae7d929a",
							"currency": "GBP",
							"amount": 19.44,
							"direction": "INBOUND",
							"created": "2018-03-28T13:48:51.000Z",
							"narrative": "someRef",
							"source": "FASTER_PAYMENTS_IN"
					  },
					  {
							"id": "94d24e13-61c1-47fe-adb2-a903a9bf6982",
							"currency": "GBP",
							"amount": 200,
							"direction": "INBOUND",
							"created": "2018-03-28T13:48:17.000Z",
							"narrative": "someRef",
							"source": "FASTER_PAYMENTS_IN"
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
							"id": "4f39ce4a-b760-42d8-811d-792e366486ef",
							"currency": "GBP",
							"amount": 33.14,
							"direction": "INBOUND",
							"created": "2018-03-28T13:48:56.000Z",
							"narrative": "someRef",
							"source": "FASTER_PAYMENTS_IN"
					  },
					  {
							"id": "e517d335-2fb8-486b-91b3-2762ae7d929a",
							"currency": "GBP",
							"amount": 19.44,
							"direction": "INBOUND",
							"created": "2018-03-28T13:48:51.000Z",
							"narrative": "someRef",
							"source": "FASTER_PAYMENTS_IN"
					  },
					  {
							"id": "94d24e13-61c1-47fe-adb2-a903a9bf6982",
							"currency": "GBP",
							"amount": 200,
							"direction": "INBOUND",
							"created": "2018-03-28T13:48:17.000Z",
							"narrative": "someRef",
							"source": "FASTER_PAYMENTS_IN"
					  }
				 ]
			}
	  }`,
	},
}

func TestFPSInTransactions(t *testing.T) {
	for _, tc := range txnsTestCasesFPSIn {
		t.Run(tc.name, func(t *testing.T) {
			testFPSInTransactions(t, tc.name, tc.mock, tc.dateRange)
		})
	}
}

func testFPSInTransactions(t *testing.T, name, mock string, dr *DateRange) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/transactions/fps/in", func(w http.ResponseWriter, r *http.Request) {
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

	got, _, err := client.FPSTransactionsIn(context.Background(), dr)
	if err != nil {
		t.Fatal("should be able to make the request", cross, err)
	}

	hal := &halTransactions{}
	json.Unmarshal([]byte(mock), hal)
	want := hal.Embedded

	if got == nil {
		t.Fatal("should not return 'nil'", cross)
	}

	if !reflect.DeepEqual(got, want.Transactions) {
		t.Error("should return a list matching the mock response", cross)
	}

	if len(got) == 0 {
		t.Errorf("should have at least one transaction %s %d", cross, len(got))
	}

	first := (got)[0]

	if first.UID == "" {
		t.Error("should have a UID specified", cross)
	}

	if first.Currency == "" {
		t.Error("should have a Currency specified", cross)
	}

	if first.Direction == "" {
		t.Error("should have a Direction specified", cross)
	}

	if first.Created == "" {
		t.Error("should have a Created date specified", cross)
	}

	if first.Narrative == "" {
		t.Error("should have a Narrative specified", cross)
	}

	if first.Source == "" {
		t.Error("should have a Source specified", cross)
	}

}

func TestFPSTxnsInForbidden(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/transactions/fps/in", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusForbidden)
	})

	got, resp, err := client.FPSTransactionsIn(context.Background(), nil)
	checkHasError(t, err)

	if resp.StatusCode != http.StatusForbidden {
		t.Error("should return HTTP 403 status")
	}

	if got != nil {
		t.Error("should not return a slice of transactions")
	}
}

var txnTestCasesFPSIn = []struct {
	name string
	uid  string
	mock string
}{
	{
		name: "single direct-debit transaction",
		uid:  "4f39ce4a-b760-42d8-811d-792e366486ef",
		mock: `{
			"id": "4f39ce4a-b760-42d8-811d-792e366486ef",
			"currency": "GBP",
			"amount": 33.14,
			"direction": "INBOUND",
			"created": "2018-03-28T13:48:56.000Z",
			"narrative": "someRef",
			"source": "FASTER_PAYMENTS_IN"
	  }`,
	},
}

func TestFPSInTransaction(t *testing.T) {
	for _, tc := range txnTestCasesDD {
		t.Run(tc.name, func(t *testing.T) {
			testFPSInTransaction(t, tc.name, tc.uid, tc.mock)
		})
	}
}

func testFPSInTransaction(t *testing.T, name, uid, mock string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/transactions/fps/in/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)

		reqUID := path.Base(r.URL.Path)
		if reqUID != uid {
			t.Error("should send a request with the correct UID", cross, reqUID)
		}

		dir := path.Dir(r.URL.Path)
		if dir != "/api/v1/transactions/fps/in" {
			t.Error("should send a request to the correct endpoint", cross, dir)
		}

		fmt.Fprint(w, mock)
	})

	want := &Transaction{}
	json.Unmarshal([]byte(mock), want)

	got, _, err := client.FPSTransactionIn(context.Background(), uid)
	if err != nil {
		t.Fatal("should be able to make the request", cross, err)
	}

	if got == nil {
		t.Fatal("should not return 'nil'", cross)
	}

	if !reflect.DeepEqual(got, want) {
		t.Error("should return a transaction matching the mock response", cross)
	}

	if got.UID == "" {
		t.Error("should have a UID specified", cross)
	}

	if got.Currency == "" {
		t.Error("should have a Currency specified", cross)
	}

	if got.Direction == "" {
		t.Error("should have a Direction specified", cross)
	}

	if got.Created == "" {
		t.Error("should have a Created date specified", cross)
	}

	if got.Narrative == "" {
		t.Error("should have a Narrative specified", cross)
	}

	if got.Source == "" {
		t.Error("should have a Source specified", cross)
	}
}

func TestFPSTxnInForbidden(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/transactions/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusForbidden)
	})

	got, resp, err := client.FPSTransactionIn(context.Background(), "nil")
	checkHasError(t, err)

	if resp.StatusCode != http.StatusForbidden {
		t.Error("should return HTTP 403 status")
	}

	if got != nil {
		t.Error("should not return a transaction")
	}
}

var txnsTestCasesFPSOut = []struct {
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
							"id": "7d3e646a-a485-41af-bd3a-d46bbb3aca8f",
							"currency": "GBP",
							"amount": -12.46,
							"direction": "OUTBOUND",
							"created": "2018-03-28T13:48:49.702Z",
							"narrative": "External Payment",
							"source": "FASTER_PAYMENTS_OUT"
					  },
					  {
							"id": "93ad883d-0883-48b9-82c1-dbe3ff57d5c8",
							"currency": "GBP",
							"amount": -14.85,
							"direction": "OUTBOUND",
							"created": "2018-03-28T13:48:48.832Z",
							"narrative": "External Payment",
							"source": "FASTER_PAYMENTS_OUT"
					  },
					  {
							"id": "c1d3b7ff-dc46-4391-82e3-6ccef72be971",
							"currency": "GBP",
							"amount": -31.17,
							"direction": "OUTBOUND",
							"created": "2018-03-28T13:48:40.047Z",
							"narrative": "External Payment",
							"source": "FASTER_PAYMENTS_OUT"
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
							"id": "7d3e646a-a485-41af-bd3a-d46bbb3aca8f",
							"currency": "GBP",
							"amount": -12.46,
							"direction": "OUTBOUND",
							"created": "2018-03-28T13:48:49.702Z",
							"narrative": "External Payment",
							"source": "FASTER_PAYMENTS_OUT"
					  },
					  {
							"id": "93ad883d-0883-48b9-82c1-dbe3ff57d5c8",
							"currency": "GBP",
							"amount": -14.85,
							"direction": "OUTBOUND",
							"created": "2018-03-28T13:48:48.832Z",
							"narrative": "External Payment",
							"source": "FASTER_PAYMENTS_OUT"
					  },
					  {
							"id": "c1d3b7ff-dc46-4391-82e3-6ccef72be971",
							"currency": "GBP",
							"amount": -31.17,
							"direction": "OUTBOUND",
							"created": "2018-03-28T13:48:40.047Z",
							"narrative": "External Payment",
							"source": "FASTER_PAYMENTS_OUT"
					  }
				 ]
			}
	  }`,
	},
}

func TestFPSOutTransactions(t *testing.T) {
	for _, tc := range txnsTestCasesFPSIn {
		t.Run(tc.name, func(t *testing.T) {
			testFPSOutTransactions(t, tc.name, tc.mock, tc.dateRange)
		})
	}
}

func testFPSOutTransactions(t *testing.T, name, mock string, dr *DateRange) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/transactions/fps/out", func(w http.ResponseWriter, r *http.Request) {
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

	got, _, err := client.FPSTransactionsOut(context.Background(), dr)
	if err != nil {
		t.Fatal("should be able to make the request", cross, err)
	}

	hal := &halTransactions{}
	json.Unmarshal([]byte(mock), hal)
	want := hal.Embedded

	if got == nil {
		t.Fatal("should not return 'nil'", cross)
	}

	if !reflect.DeepEqual(got, want.Transactions) {
		t.Error("should return a list matching the mock response", cross)
	}

	if len(got) == 0 {
		t.Errorf("should have at least one transaction %s %d", cross, len(got))
	}

	first := (got)[0]

	if first.UID == "" {
		t.Error("should have a UID specified", cross)
	}

	if first.Currency == "" {
		t.Error("should have a Currency specified", cross)
	}

	if first.Direction == "" {
		t.Error("should have a Direction specified", cross)
	}

	if first.Created == "" {
		t.Error("should have a Created date specified", cross)
	}

	if first.Narrative == "" {
		t.Error("should have a Narrative specified", cross)
	}

	if first.Source == "" {
		t.Error("should have a Source specified", cross)
	}
}

func TestFPSTxnsOutForbidden(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/transactions/fps/out", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusForbidden)
	})

	got, resp, err := client.FPSTransactionsOut(context.Background(), nil)
	checkHasError(t, err)

	if resp.StatusCode != http.StatusForbidden {
		t.Error("should return HTTP 403 status")
	}

	if got != nil {
		t.Error("should not return a slice of transactions")
	}
}

var txnTestCasesFPSOut = []struct {
	name string
	uid  string
	mock string
}{
	{
		name: "single outbound faster payments transaction",
		uid:  "4f39ce4a-b760-42d8-811d-792e366486ef",
		mock: `{
			"id": "7d3e646a-a485-41af-bd3a-d46bbb3aca8f",
			"currency": "GBP",
			"amount": -12.46,
			"direction": "OUTBOUND",
			"created": "2018-03-28T13:48:49.702Z",
			"narrative": "External Payment",
			"source": "FASTER_PAYMENTS_OUT"
	  }`,
	},
}

func TestFPSOutTransaction(t *testing.T) {
	for _, tc := range txnTestCasesDD {
		t.Run(tc.name, func(t *testing.T) {
			testFPSOutTransaction(t, tc.name, tc.uid, tc.mock)
		})
	}
}

func testFPSOutTransaction(t *testing.T, name, uid, mock string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/transactions/fps/out/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)

		reqUID := path.Base(r.URL.Path)
		if reqUID != uid {
			t.Error("should send a request with the correct UID", cross, reqUID)
		}

		dir := path.Dir(r.URL.Path)
		if dir != "/api/v1/transactions/fps/out" {
			t.Error("should send a request to the correct endpoint", cross, dir)
		}

		fmt.Fprint(w, mock)
	})

	want := &Transaction{}
	json.Unmarshal([]byte(mock), want)

	got, _, err := client.FPSTransactionOut(context.Background(), uid)
	if err != nil {
		t.Fatal("should be able to make the request", cross, err)
	}

	if got == nil {
		t.Fatal("should not return 'nil'", cross)
	}

	if !reflect.DeepEqual(got, want) {
		t.Error("should return a transaction matching the mock response", cross)
	}

	if got.UID == "" {
		t.Error("should have a UID specified", cross)
	}

	if got.Currency == "" {
		t.Error("should have a Currency specified", cross)
	}

	if got.Direction == "" {
		t.Error("should have a Direction specified", cross)
	}

	if got.Created == "" {
		t.Error("should have a Created date specified", cross)
	}

	if got.Narrative == "" {
		t.Error("should have a Narrative specified", cross)
	}

	if got.Source == "" {
		t.Error("should have a Source specified", cross)
	}
}

func TestFPSTxnOutForbidden(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/transactions/fps/out/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusForbidden)
	})

	got, resp, err := client.FPSTransactionOut(context.Background(), "nil")
	checkHasError(t, err)

	if resp.StatusCode != http.StatusForbidden {
		t.Error("should return HTTP 403 status")
	}

	if got != nil {
		t.Error("should not return a transaction")
	}
}

var txnsTestCasesMastercard = []struct {
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
								 "merchant": {
									  "href": "/api/v1/merchants/e5b21fd4-fb62-4c40-8f56-4890a16bab51",
									  "templated": false
								 },
								 "merchantLocation": {
									  "href": "/api/v1/merchants/e5b21fd4-fb62-4c40-8f56-4890a16bab51/locations/96fc27d9-164d-4fb5-a60f-799e9b2b2294",
									  "templated": false
								 }
							},
							"id": "1b19cf9e-0499-4a99-83b5-b442db58d176",
							"currency": "GBP",
							"amount": -15.51,
							"direction": "OUTBOUND",
							"created": "2018-03-28T13:48:59.434Z",
							"narrative": "Sofitel",
							"source": "MASTER_CARD",
							"mastercardTransactionMethod": "CHIP_AND_PIN",
							"status": "SETTLED",
							"sourceAmount": -15.51,
							"sourceCurrency": "GBP",
							"merchantId": "e5b21fd4-fb62-4c40-8f56-4890a16bab51",
							"merchantLocationId": "96fc27d9-164d-4fb5-a60f-799e9b2b2294",
							"spendingCategory": "HOLIDAYS",
							"country": "GBR",
							"posTimestamp": "13:48:58",
							"authorisationCode": "597448",
							"eventUid": "5a9ca682-8b8e-40bc-887d-43e00617e2b9",
							"cardLast4": "0142"
					  },
					  {
							"_links": {
								 "merchant": {
									  "href": "/api/v1/merchants/c052f76f-e919-427d-85fc-f46a75a3ff26",
									  "templated": false
								 },
								 "merchantLocation": {
									  "href": "/api/v1/merchants/c052f76f-e919-427d-85fc-f46a75a3ff26/locations/371c62bc-dcfc-4799-8b23-b070626772f7",
									  "templated": false
								 }
							},
							"id": "ab4b76ef-8e98-4181-8be9-164d73e16d99",
							"currency": "GBP",
							"amount": -35.24,
							"direction": "OUTBOUND",
							"created": "2018-03-28T13:48:58.677Z",
							"narrative": "Yorkshire Bank III",
							"source": "MASTER_CARD",
							"mastercardTransactionMethod": "CHIP_AND_PIN",
							"status": "SETTLED",
							"sourceAmount": -35.24,
							"sourceCurrency": "GBP",
							"merchantId": "c052f76f-e919-427d-85fc-f46a75a3ff26",
							"merchantLocationId": "371c62bc-dcfc-4799-8b23-b070626772f7",
							"spendingCategory": "HOLIDAYS",
							"country": "GBR",
							"posTimestamp": "13:48:58",
							"authorisationCode": "639682",
							"eventUid": "feada116-963e-4246-8fac-ed5e499b03d5",
							"cardLast4": "0142"
					  },
					  {
							"_links": {
								 "merchant": {
									  "href": "/api/v1/merchants/acc81c26-071b-4e22-904d-5635e4bd6089",
									  "templated": false
								 },
								 "merchantLocation": {
									  "href": "/api/v1/merchants/acc81c26-071b-4e22-904d-5635e4bd6089/locations/5c1abdb1-8a8f-4a9f-a415-d8a5af8df0f6",
									  "templated": false
								 }
							},
							"id": "a7593bd6-b003-4076-9c12-198fdc53fa25",
							"currency": "GBP",
							"amount": -12.3,
							"direction": "OUTBOUND",
							"created": "2018-03-28T13:48:58.675Z",
							"narrative": "Sofitel",
							"source": "MASTER_CARD",
							"mastercardTransactionMethod": "CHIP_AND_PIN",
							"status": "SETTLED",
							"sourceAmount": -12.3,
							"sourceCurrency": "GBP",
							"merchantId": "acc81c26-071b-4e22-904d-5635e4bd6089",
							"merchantLocationId": "5c1abdb1-8a8f-4a9f-a415-d8a5af8df0f6",
							"spendingCategory": "HOLIDAYS",
							"country": "GBR",
							"posTimestamp": "13:48:58",
							"authorisationCode": "675742",
							"eventUid": "1b448a68-1bc5-46b3-9b52-aab689511516",
							"cardLast4": "0142"
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
								 "merchant": {
									  "href": "/api/v1/merchants/e5b21fd4-fb62-4c40-8f56-4890a16bab51",
									  "templated": false
								 },
								 "merchantLocation": {
									  "href": "/api/v1/merchants/e5b21fd4-fb62-4c40-8f56-4890a16bab51/locations/96fc27d9-164d-4fb5-a60f-799e9b2b2294",
									  "templated": false
								 }
							},
							"id": "1b19cf9e-0499-4a99-83b5-b442db58d176",
							"currency": "GBP",
							"amount": -15.51,
							"direction": "OUTBOUND",
							"created": "2018-03-28T13:48:59.434Z",
							"narrative": "Sofitel",
							"source": "MASTER_CARD",
							"mastercardTransactionMethod": "CHIP_AND_PIN",
							"status": "SETTLED",
							"sourceAmount": -15.51,
							"sourceCurrency": "GBP",
							"merchantId": "e5b21fd4-fb62-4c40-8f56-4890a16bab51",
							"merchantLocationId": "96fc27d9-164d-4fb5-a60f-799e9b2b2294",
							"spendingCategory": "HOLIDAYS",
							"country": "GBR",
							"posTimestamp": "13:48:58",
							"authorisationCode": "597448",
							"eventUid": "5a9ca682-8b8e-40bc-887d-43e00617e2b9",
							"cardLast4": "0142"
					  },
					  {
							"_links": {
								 "merchant": {
									  "href": "/api/v1/merchants/c052f76f-e919-427d-85fc-f46a75a3ff26",
									  "templated": false
								 },
								 "merchantLocation": {
									  "href": "/api/v1/merchants/c052f76f-e919-427d-85fc-f46a75a3ff26/locations/371c62bc-dcfc-4799-8b23-b070626772f7",
									  "templated": false
								 }
							},
							"id": "ab4b76ef-8e98-4181-8be9-164d73e16d99",
							"currency": "GBP",
							"amount": -35.24,
							"direction": "OUTBOUND",
							"created": "2018-03-28T13:48:58.677Z",
							"narrative": "Yorkshire Bank III",
							"source": "MASTER_CARD",
							"mastercardTransactionMethod": "CHIP_AND_PIN",
							"status": "SETTLED",
							"sourceAmount": -35.24,
							"sourceCurrency": "GBP",
							"merchantId": "c052f76f-e919-427d-85fc-f46a75a3ff26",
							"merchantLocationId": "371c62bc-dcfc-4799-8b23-b070626772f7",
							"spendingCategory": "HOLIDAYS",
							"country": "GBR",
							"posTimestamp": "13:48:58",
							"authorisationCode": "639682",
							"eventUid": "feada116-963e-4246-8fac-ed5e499b03d5",
							"cardLast4": "0142"
					  },
					  {
							"_links": {
								 "merchant": {
									  "href": "/api/v1/merchants/acc81c26-071b-4e22-904d-5635e4bd6089",
									  "templated": false
								 },
								 "merchantLocation": {
									  "href": "/api/v1/merchants/acc81c26-071b-4e22-904d-5635e4bd6089/locations/5c1abdb1-8a8f-4a9f-a415-d8a5af8df0f6",
									  "templated": false
								 }
							},
							"id": "a7593bd6-b003-4076-9c12-198fdc53fa25",
							"currency": "GBP",
							"amount": -12.3,
							"direction": "OUTBOUND",
							"created": "2018-03-28T13:48:58.675Z",
							"narrative": "Sofitel",
							"source": "MASTER_CARD",
							"mastercardTransactionMethod": "CHIP_AND_PIN",
							"status": "SETTLED",
							"sourceAmount": -12.3,
							"sourceCurrency": "GBP",
							"merchantId": "acc81c26-071b-4e22-904d-5635e4bd6089",
							"merchantLocationId": "5c1abdb1-8a8f-4a9f-a415-d8a5af8df0f6",
							"spendingCategory": "HOLIDAYS",
							"country": "GBR",
							"posTimestamp": "13:48:58",
							"authorisationCode": "675742",
							"eventUid": "1b448a68-1bc5-46b3-9b52-aab689511516",
							"cardLast4": "0142"
					  }
				 ]
			}
	  }`,
	},
}

func TestMastercardTransactions(t *testing.T) {
	for _, tc := range txnsTestCasesMastercard {
		t.Run(tc.name, func(t *testing.T) {
			testMastercardTransactions(t, tc.name, tc.mock, tc.dateRange)
		})
	}
}

func testMastercardTransactions(t *testing.T, name, mock string, dr *DateRange) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/transactions/mastercard", func(w http.ResponseWriter, r *http.Request) {
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

		resource := path.Base(r.URL.Path)

		if resource != "mastercard" {
			t.Error("should request mastercard", cross, resource)
		}

		fmt.Fprint(w, mock)
	})

	got, _, err := client.MastercardTransactions(context.Background(), dr)
	if err != nil {
		t.Fatal("should be able to make the request", cross, err)
	}

	hal := &halMastercardTransactions{}
	json.Unmarshal([]byte(mock), hal)
	want := hal.Embedded

	if got == nil {
		t.Fatal("should not return 'nil'", cross)
	}

	if !reflect.DeepEqual(got, want.Transactions) {
		t.Error("should return a list matching the mock response", cross)
	}

	if len(got) == 0 {
		t.Errorf("should have at least one transaction %s %d", cross, len(got))
	}

	first := (got)[0]

	if first.UID == "" {
		t.Error("should have a UID specified", cross)
	}

	if first.Currency == "" {
		t.Error("should have a Currency specified", cross)
	}

	if first.Direction == "" {
		t.Error("should have a Direction specified", cross)
	}

	if first.Created == "" {
		t.Error("should have a Created date specified", cross)
	}

	if first.Narrative == "" {
		t.Error("should have a Narrative specified", cross)
	}

	if first.Source == "" {
		t.Error("should have a Source specified", cross)
	}
}

func TestTxnsCardForbidden(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/transactions/mastercard", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusForbidden)
	})

	got, resp, err := client.MastercardTransactions(context.Background(), nil)
	checkHasError(t, err)

	if resp.StatusCode != http.StatusForbidden {
		t.Error("should return HTTP 403 status")
	}

	if got != nil {
		t.Error("should not return a slice of transactions")
	}
}

var txnTestCasesMastercard = []struct {
	name string
	uid  string
	mock string
}{
	{
		name: "single mastercard transaction",
		uid:  "1b19cf9e-0499-4a99-83b5-b442db58d176",
		mock: `{
			"_links": {
				 "merchant": {
					  "href": "/api/v1/merchants/e5b21fd4-fb62-4c40-8f56-4890a16bab51",
					  "templated": false
				 },
				 "merchantLocation": {
					  "href": "/api/v1/merchants/e5b21fd4-fb62-4c40-8f56-4890a16bab51/locations/96fc27d9-164d-4fb5-a60f-799e9b2b2294",
					  "templated": false
				 }
			},
			"id": "1b19cf9e-0499-4a99-83b5-b442db58d176",
			"currency": "GBP",
			"amount": -15.51,
			"direction": "OUTBOUND",
			"created": "2018-03-28T13:48:59.434Z",
			"narrative": "Sofitel",
			"source": "MASTER_CARD",
			"mastercardTransactionMethod": "CHIP_AND_PIN",
			"status": "SETTLED",
			"sourceAmount": -15.51,
			"sourceCurrency": "GBP",
			"merchantId": "e5b21fd4-fb62-4c40-8f56-4890a16bab51",
			"merchantLocationId": "96fc27d9-164d-4fb5-a60f-799e9b2b2294",
			"spendingCategory": "HOLIDAYS",
			"country": "GBR",
			"posTimestamp": "13:48:58",
			"authorisationCode": "597448",
			"eventUid": "5a9ca682-8b8e-40bc-887d-43e00617e2b9",
			"cardLast4": "0142"
	  }`,
	},
}

func TestMastercardTransaction(t *testing.T) {
	for _, tc := range txnTestCasesMastercard {
		t.Run(tc.name, func(t *testing.T) {
			testMastercardTransaction(t, tc.name, tc.uid, tc.mock)
		})
	}
}

func testMastercardTransaction(t *testing.T, name, uid, mock string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/transactions/mastercard/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)

		reqUID := path.Base(r.URL.Path)
		if reqUID != uid {
			t.Error("should send a request with the correct UID", cross, reqUID)
		}

		dir := path.Dir(r.URL.Path)
		if dir != "/api/v1/transactions/mastercard" {
			t.Error("should send a request to the correct endpoint", cross, dir)
		}

		fmt.Fprint(w, mock)
	})

	want := &MastercardTransaction{}
	json.Unmarshal([]byte(mock), want)

	got, _, err := client.MastercardTransaction(context.Background(), uid)
	if err != nil {
		t.Fatal("should be able to make the request", cross, err)
	}

	if got == nil {
		t.Fatal("should not return 'nil'", cross)
	}

	if !reflect.DeepEqual(got, want) {
		t.Error("should return a transaction matching the mock response", cross)
	}

	if got.UID == "" {
		t.Error("should have a UID specified", cross)
	}

	if got.Currency == "" {
		t.Error("should have a Currency specified", cross)
	}

	if got.Direction == "" {
		t.Error("should have a Direction specified", cross)
	}

	if got.Created == "" {
		t.Error("should have a Created date specified", cross)
	}

	if got.Narrative == "" {
		t.Error("should have a Narrative specified", cross)
	}

	if got.Source == "" {
		t.Error("should have a Source specified", cross)
	}

	if got.CardLast4 == "" {
		t.Error("should have the last four digits of the card number specified", cross)
	}

	if got.POSTimestamp == "" {
		t.Error("should have the POS timestamp specified", cross)
	}
}

func TestCardTxnForbidden(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/transactions/mastercard/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusForbidden)
	})

	got, resp, err := client.MastercardTransaction(context.Background(), "nil")
	checkHasError(t, err)

	if resp.StatusCode != http.StatusForbidden {
		t.Error("should return HTTP 403 status")
	}

	if got != nil {
		t.Error("should not return a transaction")
	}
}

var setMastercardSpendingCategoryCases = []struct {
	name     string
	uid      string
	category string
	status   int
	mock     string
}{
	{
		name:     "set transaction category to charity",
		uid:      "474642e6-c4f5-43af-9b93-fe5ddbfcb857",
		category: "CHARITY",
		status:   http.StatusAccepted,
		mock:     ``,
	},
	{
		name:     "set invalid transaction category",
		uid:      "474642e6-c4f5-43af-9b93-fe5ddbfcb857",
		category: "INVALID",
		status:   http.StatusBadRequest,
		mock: `[
			"Can not deserialize value of type com.starlingbank.connectors.customer.SpendingCategory from String \"DEMOX\": value not one of declared Enum instance names: [GIFTS, FAMILY, ENTERTAINMENT, TRANSPORT, GROCERIES, PAYMENTS, PETS, LIFESTYLE, CHARITY, BILLS_AND_SERVICES, SAVING, HOLIDAYS, HOME, GENERAL, NONE, EXPENSES, INCOME, SHOPPING, EATING_OUT]"
	  ]`,
	},
}

func TestSetMastercardSpendingCategory(t *testing.T) {
	for _, tc := range setDDSpendingCategoryCases {
		t.Run(tc.name, func(st *testing.T) {
			testSetMastercardSpendingCategory(st, tc.name, tc.uid, tc.category, tc.status, tc.mock)
		})
	}
}

func testSetMastercardSpendingCategory(t *testing.T, name, uid, cat string, status int, mock string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/transactions/mastercard/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodPut)

		reqUID := path.Base(r.URL.Path)
		resource := path.Base(path.Dir(r.URL.Path))

		if reqUID != uid {
			t.Error("should send a request with the correct UID", cross, reqUID)
		}

		if resource != "mastercard" {
			t.Error("should request mastercard", cross, resource)
		}

		var reqCat = SpendingCategory{}
		err := json.NewDecoder(r.Body).Decode(&reqCat)
		if err != nil {
			t.Fatal("should send a request that the API can parse", cross, err)
		}

		if reqCat.SpendingCategory != cat {
			t.Error("should request the correct spending category", cross, reqCat)
		}

		w.WriteHeader(status)
		fmt.Fprintln(w, mock)
	})

	resp, err := client.SetMastercardSpendingCategory(context.Background(), uid, cat)
	if status >= 400 {
		checkHasError(t, err)
	} else {
		checkNoError(t, err)
	}

	if resp.StatusCode != status {
		t.Error("should return the correct status", cross, resp.Status)
	}
}

func TestCardSpendingCategoryForbidden(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/transactions/mastercard/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodPut)
		w.WriteHeader(http.StatusForbidden)
	})

	resp, err := client.SetMastercardSpendingCategory(context.Background(), "474642e6-c4f5-43af-9b93-fe5ddbfcb857", "CHARITY")
	checkHasError(t, err)

	if resp.StatusCode != http.StatusForbidden {
		t.Error("should return HTTP 403 status")
	}
}
