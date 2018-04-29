package starling

import (
	"context"
	"encoding/json"
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
