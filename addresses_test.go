package starling

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var addressesTestCases = []struct {
	name string
	mock string
}{
	{
		name: "single address",
		mock: `{
			"current": {
				"line1": "1A Admiralty Arch",
				"line2": "The Mall",
				"line3": "City of Westminster",
				"postTown": "London",
				"countryCode": "GB",
				"postCode": "SW1A 2WH"
			},
			"previous": []
		}`,
	},
	{
		name: "single previous address",
		mock: `{
			"current": {
				"line1": "1A Admiralty Arch",
				"line2": "The Mall",
				"line3": "City of Westminster",
				"postTown": "London",
				"countryCode": "GB",
				"postCode": "SW1A 2WH"
			},
			"previous": [{
				"line1": "1A Admiralty Arch",
				"line2": "The Mall",
				"line3": "City of Westminster",
				"postTown": "London",
				"countryCode": "GB",
				"postCode": "SW1A 2WH"
			}]
		}`,
	},
	{
		name: "multiple previous addresses",
		mock: `{
			"current": {
				"line1": "1A Admiralty Arch",
				"line2": "The Mall",
				"line3": "City of Westminster",
				"postTown": "London",
				"countryCode": "GB",
				"postCode": "SW1A 2WH"
			},
			"previous": [{
				"line1": "1A Admiralty Arch",
				"line2": "The Mall",
				"line3": "City of Westminster",
				"postTown": "London",
				"countryCode": "GB",
				"postCode": "SW1A 2WH"
			},
			{
				"line1": "1A Admiralty Arch",
				"line2": "The Mall",
				"line3": "City of Westminster",
				"postTown": "London",
				"countryCode": "GB",
				"postCode": "SW1A 2WH"
			}]
		}`,
	},
}

func TestAddressHistory(t *testing.T) {
	for _, tc := range addressesTestCases {
		t.Run(tc.name, func(st *testing.T) {
			testAddressHistory(st, tc.name, tc.mock)
		})
	}
}

func testAddressHistory(t *testing.T, name, mock string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v2/addresses", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		fmt.Fprint(w, mock)
	})

	got, _, err := client.AddressHistory(context.Background())
	checkNoError(t, err)

	want := &AddressHistory{}
	json.Unmarshal([]byte(mock), want)

	if !reflect.DeepEqual(got, want) {
		t.Error("should return addresses matching the mock response", cross)
	}
}

func TestAddressHistoryForbidden(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v2/addresses", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusForbidden)
	})

	got, resp, err := client.AddressHistory(context.Background())
	checkHasError(t, err)

	if resp.StatusCode != http.StatusForbidden {
		t.Error("should return HTTP 403 status")
	}

	if got != nil {
		t.Error("should not return any addresses")
	}
}
