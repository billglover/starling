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
				"streetAddress": "3rd Floor",
				"city": "London",
				"country": "GBR",
				"postcode": " EC2M 2PP"
			},
			"previous": []
		}`,
	},
	{
		name: "single previous address",
		mock: `{
			"current": {
				"streetAddress": "3rd Floor",
				"city": "London",
				"country": "GBR",
				"postcode": " EC2M 2PP"
			},
			"previous": [{
				"streetAddress": "3rd Floor",
				"city": "London",
				"country": "GBR",
				"postcode": " EC2M 2PP"
			}]
		}`,
	},
	{
		name: "multiple previous addresses",
		mock: `{
			"current": {
				"streetAddress": "3rd Floor",
				"city": "London",
				"country": "GBR",
				"postcode": " EC2M 2PP"
			},
			"previous": [{
				"streetAddress": "3rd Floor",
				"city": "London",
				"country": "GBR",
				"postcode": " EC2M 2PP"
			},
			{
				"streetAddress": "3rd Floor",
				"city": "London",
				"country": "GBR",
				"postcode": " EC2M 2PP"
			}]
		}`,
	},
}

func TestGetAddresses(t *testing.T) {

	t.Log("Given the need to test fetching customer addresses:")

	for _, tc := range addressesTestCases {
		t.Run(tc.name, func(st *testing.T) {
			testGetAddresses(st, tc.name, tc.mock)
		})
	}
}

func testGetAddresses(t *testing.T, name, mock string) {
	t.Logf("\tWhen making a call to GetAddresses() with a %s:", name)

	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/addresses", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		fmt.Fprint(w, mock)
	})

	got, _, err := client.GetAddresses(context.Background())
	checkNoError(t, err)

	want := &Addresses{}
	json.Unmarshal([]byte(mock), want)

	if !reflect.DeepEqual(got, want) {
		t.Error("\t\tshould return addresses matching the mock response", cross)
	} else {
		t.Log("\t\tshould return addresses matching the mock response", tick)
	}
}
