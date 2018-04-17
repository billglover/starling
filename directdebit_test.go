package starling

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var ddTestCases = []struct {
	name string
	mock string
}{
	{
		name: "empty dd list",
		mock: `{
	"_links": {
		 "self": {
			  "href": "/api/v1/direct-debit/mandates",
			  "templated": false
		 }
	},
	"_embedded": {
		 "mandates": []
	}
}`,
	},
	{
		name: "single dd",
		mock: `{
			"_links": {
				 "self": {
					  "href": "/api/v1/direct-debit/mandates",
					  "templated": false
				 }
			},
			"_embedded": {
				 "mandates": [
					  {
							"uid": "fa7998f6-07ce-42a9-ba5b-ce45ea8aff89",
							"reference": "VolcanoDisruptions",
							"status": "LIVE",
							"source": "ELECTRONIC",
							"created": "2018-04-17T07:23:59.173Z",
							"originatorName": "ANTIQUARIES",
							"originatorUid": "949404bd-d32e-4f1e-9759-4d6caee3137c"
					  }
				 ]
			}
	  }`,
	},
}

func TestDirectDebits(t *testing.T) {
	for _, tc := range ddTestCases {
		t.Run(tc.name, func(st *testing.T) {
			testDirectDebits(st, tc.name, tc.mock)
		})
	}
}

func testDirectDebits(t *testing.T, name, mock string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/direct-debit/mandates", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		fmt.Fprint(w, mock)
	})

	got, _, err := client.DirectDebitMandates(context.Background())
	checkNoError(t, err)

	hal := &halDirectDebitMandates{}
	json.Unmarshal([]byte(mock), hal)
	want := hal.Embedded

	if !reflect.DeepEqual(got, &want.Mandates) {
		t.Error("\t\tshould return a list of mandates matching the mock response", cross)
	}
}
