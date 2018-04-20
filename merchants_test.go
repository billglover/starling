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

var merchantCases = []struct {
	name string
	uid  string
	mock string
}{
	{
		name: "sample merchant",
		uid:  "c052f76f-e919-427d-85fc-f46a75a3ff26",
		mock: `{
			"merchantUid": "c052f76f-e919-427d-85fc-f46a75a3ff26",
			"name": "Mastercard",
			"website": "http://mastercard.co.uk",
			"twitterUsername": "@mastercard"
	  }`,
	},
}

func TestMerchant(t *testing.T) {
	for _, tc := range merchantCases {
		t.Run(tc.name, func(st *testing.T) {
			testMerchant(st, tc.name, tc.uid, tc.mock)
		})
	}
}

func testMerchant(t *testing.T, name, uid, mock string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/merchants/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)

		reqUID := path.Base(r.URL.Path)
		if reqUID != uid {
			t.Error("should send a request with the correct UID", cross, reqUID)
		}

		fmt.Fprint(w, mock)
	})

	got, _, err := client.Merchant(context.Background(), uid)
	checkNoError(t, err)

	mer := &Merchant{}
	json.Unmarshal([]byte(mock), mer)

	if !reflect.DeepEqual(got, mer) {
		t.Error("should return a merchant matching the mock response", cross)
	}
}

var merchantLocationCases = []struct {
	name string
	mUID string
	lUID string
	mock string
}{
	{
		name: "sample merchant location",
		mUID: "c052f76f-e919-427d-85fc-f46a75a3ff26",
		lUID: "371c62bc-dcfc-4799-8b23-b070626772f7",
		mock: `{
			"merchantUid": "c052f76f-e919-427d-85fc-f46a75a3ff26",
			"merchant": {
				 "href": "/api/v1/merchants/c052f76f-e919-427d-85fc-f46a75a3ff26",
				 "templated": false
			},
			"merchantLocationUid": "371c62bc-dcfc-4799-8b23-b070626772f7",
			"merchantName": "Mastercard",
			"locationName": "Mastercard UK",
			"googlePlaceId": "ChIJJ9ZEdgG4h0gRqJZP_Z5tgCc",
			"mastercardMerchantCategoryCode": 3619
	  }`,
	},
}

func TestMerchantLocation(t *testing.T) {
	for _, tc := range merchantLocationCases {
		t.Run(tc.name, func(st *testing.T) {
			testMerchantLocation(st, tc.name, tc.mUID, tc.lUID, tc.mock)
		})
	}
}

func testMerchantLocation(t *testing.T, name, mUID, lUID, mock string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/merchants/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)

		reqMerUID := path.Base(path.Dir(path.Dir(r.URL.Path)))
		if reqMerUID != mUID {
			t.Error("should send a request with the correct merchant UID", cross, reqMerUID)
		}

		resource := path.Base(path.Dir(r.URL.Path))
		if resource != "locations" {
			t.Error("should send a request for the locations resource", cross, resource)
		}

		reqLocUID := path.Base(r.URL.Path)
		if reqLocUID != lUID {
			t.Error("should send a request with the correct location UID", cross, reqLocUID)
		}

		fmt.Fprint(w, mock)
	})

	got, _, err := client.MerchantLocation(context.Background(), mUID, lUID)
	checkNoError(t, err)

	merLoc := &MerchantLocation{}
	json.Unmarshal([]byte(mock), merLoc)

	if !reflect.DeepEqual(got, merLoc) {
		t.Error("should return a merchant matching the mock response", cross)
	}
}
