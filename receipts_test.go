package starling

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

var rcptTestCases = []struct {
	name    string
	receipt Receipt
	status  int
	mock    string
}{
	{
		name: "valid receipt",
		receipt: Receipt{
			UID:                "9af397e9-63e8-4a72-b3f6-4f0b068b2ed0",
			EventUID:           "feada116-963e-4246-8fac-ed5e499b03d5",
			MetadataSource:     "tests",
			ReceiptIdentifier:  "0987654321",
			MerchantIdentifier: "2b03a13a-bbfc-4479-8d4d-abb6a9119d27",
			MerchantAddress:    "1 Merchant Way, Go Tests, UK",
			TotalAmount:        1.23,
			TotalTax:           0.45,
			TaxReference:       "1234567890",
			AuthCode:           "639682",
			CardLast4:          "9012",
			ProviderName:       "demo",
			Items: []ReceiptItem{
				{
					UID:         "d74b17d5-f114-42d2-8adb-a6ddedc29298",
					Description: "Large coffee",
					Quantity:    2,
					Amount:      12.34,
					Tax:         1.23,
					URL:         "https://crossenvsync",
				},
			},
			Notes: []ReceiptNote{
				{
					UID:         "152dcc19-c86b-4f03-ba46-82aadbdcc957",
					Description: "Large coffee",
					URL:         "https://crossenvsync",
				},
			},
		},
		status: http.StatusAccepted,
		mock:   ``,
	},
}

func TestCreateReceipt(t *testing.T) {
	for _, tc := range rcptTestCases {
		t.Run(tc.name, func(st *testing.T) {
			testCreateReceipt(st, tc.name, tc.receipt, tc.mock, tc.status)
		})
	}
}

func testCreateReceipt(t *testing.T, name string, rcpt Receipt, mock string, status int) {
	client, mux, _, teardown := setup()
	defer teardown()

	txnUID := "ab4b76ef-8e98-4181-8be9-164d73e16d99"

	mux.HandleFunc("/api/v1/transactions/mastercard/"+txnUID+"/receipt", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodPost)

		var reqRcpt = Receipt{}
		err := json.NewDecoder(r.Body).Decode(&reqRcpt)
		if err != nil {
			t.Fatal("should send a request that the API can parse", cross, err)
		}

		if !reflect.DeepEqual(rcpt, reqRcpt) {
			t.Error("should send a contact account that matches the mock", cross)
			t.Error(rcpt)
			t.Error(reqRcpt)
		}

		w.WriteHeader(http.StatusAccepted)
	})

	resp, err := client.CreateReceipt(context.Background(), txnUID, rcpt)
	if status <= 299 {
		checkNoError(t, err)
	} else {
		checkHasError(t, err)
	}

	checkStatus(t, resp, status)
}
