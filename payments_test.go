package starling

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var paymentsCasesLocal = []struct {
	name    string
	payment LocalPayment
	mock    string
}{
	{
		name: "valid payment",
		payment: LocalPayment{
			Reference:             "sample payment",
			DestinationAccountUID: "99970be2-2bc7-49d3-8d23-ebef9f746ecf",
			Payment: PaymentAmount{
				Currency: "GBP",
				Amount:   10.24,
			},
		},
		mock: `{
			"payment": {
			  "currency": "GBP",
			  "amount": 10.24
			},
			"destinationAccountUid": "99970be2-2bc7-49d3-8d23-ebef9f746ecf",
			"reference": "sample payment"
		 }`,
	},
}

func TestLocalPayment(t *testing.T) {
	for _, tc := range paymentsCasesLocal {
		t.Run(tc.name, func(st *testing.T) {
			testLocalPayment(st, tc.name, tc.payment, tc.mock)
		})
	}
}

func testLocalPayment(t *testing.T, name string, payment LocalPayment, mock string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/payments/local", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodPost)

		var reqPayment = LocalPayment{}
		err := json.NewDecoder(r.Body).Decode(&reqPayment)
		if err != nil {
			t.Fatal("should send a request that the API can parse", err)
		}

		if !reflect.DeepEqual(payment, reqPayment) {
			t.Error("should send a local payment that matches the mock", cross)
		}

		w.WriteHeader(http.StatusAccepted)
	})

	resp, err := client.MakeLocalPayment(context.Background(), payment)
	checkNoError(t, err)
	checkStatus(t, resp, http.StatusAccepted)

}

var paymentsCasesScheduled = []struct {
	name string
	mock string
}{
	{
		name: "list of scheduled payments",
		mock: `{
			"_links": {
				 "nextPage": {
					  "href": "NOT_YET_IMPLEMENTED"
				 }
			},
			"_embedded": {
				 "paymentOrders": [
					  {
							"_links": {
								 "receivingContactAccount": {
									  "href": "api/v1/contacts/157e8e67-c642-427a-a62c-b978fb6a6f55/accounts/2f543dbd-a8dc-443f-8962-521bbb45b5b6",
									  "templated": false
								 }
							},
							"paymentOrderId": "1e22a383-0dd6-4845-a5fd-17c55920381d",
							"currency": "GBP",
							"amount": 16.55,
							"reference": "External Payment",
							"receivingContactAccountId": "2f543dbd-a8dc-443f-8962-521bbb45b5b6",
							"recipientName": null,
							"immediate": true,
							"startDate": "2018-03-09",
							"nextDate": "2018-03-09",
							"endDate": "2018-03-09",
							"paymentType": "STANDING_ORDER"
					  },
					  {
							"_links": {
								 "receivingContactAccount": {
									  "href": "api/v1/contacts/819e5a8f-54b5-4638-b961-9492ffd0d142/accounts/99970be2-2bc7-49d3-8d23-ebef9f746ecf",
									  "templated": false
								 }
							},
							"paymentOrderId": "f8e714f1-f5a3-4bd8-a6f5-28e44a6b1416",
							"currency": "GBP",
							"amount": 10.24,
							"reference": "Dinner",
							"receivingContactAccountId": "99970be2-2bc7-49d3-8d23-ebef9f746ecf",
							"recipientName": null,
							"immediate": true,
							"startDate": "2018-04-29",
							"nextDate": "2018-04-29",
							"endDate": "2018-04-29",
							"paymentType": "STANDING_ORDER"
					  }
				 ]
			}
	  }`,
	},
	{
		name: "blank response",
		mock: ``,
	},
}

func TestScheduledPayment(t *testing.T) {
	for _, tc := range paymentsCasesScheduled {
		t.Run(tc.name, func(st *testing.T) {
			testScheduledPayment(st, tc.name, tc.mock)
		})
	}
}

func testScheduledPayment(t *testing.T, name string, mock string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/payments/scheduled", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, mock)
	})

	got, resp, err := client.ScheduledPayments(context.Background())
	checkNoError(t, err)
	checkStatus(t, resp, http.StatusOK)

	hPO := new(halPaymentOrders)
	json.Unmarshal([]byte(mock), hPO)
	var want []ScheduledPayment
	if hPO.Embedded != nil {
		want = hPO.Embedded.PaymentOrders
	}

	if !reflect.DeepEqual(got, want) {
		t.Error("should return a list of scheduled payments matching the mock")
	}

	if len(want) > 0 {
		if got[0].UID == "" {
			t.Error("first scheduled payment should have a non zero UID")
		}
	}
}
