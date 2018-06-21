package starling

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
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

	if !reflect.DeepEqual(got, want.Mandates) {
		t.Error("should return a list of mandates matching the mock response", cross)
	}
}

func TestDDMandatesForbidden(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/direct-debit/mandates", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusForbidden)
	})

	got, resp, err := client.DirectDebitMandates(context.Background())
	checkHasError(t, err)

	if resp.StatusCode != http.StatusForbidden {
		t.Error("should return HTTP 403 status")
	}

	if got != nil {
		t.Error("should not return direct-debit mandates")
	}
}

var ddMandateCases = []struct {
	name string
	uid  string
	mock string
}{
	{
		name: "direct debit mandate",
		uid:  "fa7998f6-07ce-42a9-ba5b-ce45ea8aff89",
		mock: `{
			"uid": "fa7998f6-07ce-42a9-ba5b-ce45ea8aff89",
			"reference": "VolcanoDisruptions",
			"status": "LIVE",
			"source": "ELECTRONIC",
			"created": "2018-04-17T07:23:59.173Z",
			"originatorName": "ANTIQUARIES",
			"originatorUid": "949404bd-d32e-4f1e-9759-4d6caee3137c"
	  }`,
	},
}

func TestDDMandate(t *testing.T) {
	for _, tc := range ddMandateCases {
		t.Run(tc.name, func(st *testing.T) {
			testDDMandate(st, tc.name, tc.uid, tc.mock)
		})
	}
}

func testDDMandate(t *testing.T, name, uid, mock string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/direct-debit/mandates/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)

		reqUID := path.Base(r.URL.Path)
		if reqUID != uid {
			t.Error("should send a requestwith the correct UID", cross, reqUID)
		}

		fmt.Fprint(w, mock)
	})

	got, _, err := client.DirectDebitMandate(context.Background(), uid)
	checkNoError(t, err)

	want := &DirectDebitMandate{}
	json.Unmarshal([]byte(mock), want)

	if !reflect.DeepEqual(got, want) {
		t.Error("should return a single mandate matching the mock response", cross)
	}

	if got.UID != want.UID {
		t.Error("should have the correct UID", cross, got.UID)
	}
}

func TestDDMandateForbidden(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/direct-debit/mandates/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusForbidden)
	})

	got, resp, err := client.DirectDebitMandate(context.Background(), "949404bd-d32e-4f1e-9759-4d6caee3137c")
	checkHasError(t, err)

	if resp.StatusCode != http.StatusForbidden {
		t.Error("should return HTTP 403 status")
	}

	if got != nil {
		t.Error("should not return a direct-debit mandate")
	}
}

var deleteDDMandateCases = []struct {
	name string
	uid  string
}{
	{
		name: "sample direct debit mandate",
		uid:  "840e4030-b94c-4e71-a1d3-1319a233dd3c",
	},
}

func TestDeleteDDMandate(t *testing.T) {
	for _, tc := range deleteDDMandateCases {
		t.Run(tc.name, func(st *testing.T) {
			testDeleteDDMandate(st, tc.name, tc.uid)
		})
	}
}

func testDeleteDDMandate(t *testing.T, name, uid string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/direct-debit/mandates/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodDelete)

		reqUID := path.Base(r.URL.Path)
		if reqUID != uid {
			t.Error("should send a requestwith the correct UID", cross, reqUID)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.DeleteDirectDebitMandate(context.Background(), uid)
	checkNoError(t, err)

	if resp.StatusCode != http.StatusNoContent {
		t.Error("should return an HTTP 204 status", cross, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	checkNoError(t, err)

	if len(body) != 0 {
		t.Error("should return an empty body", cross, len(body))
	}
}

func TestDDMandateDeleteForbidden(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/direct-debit/mandates/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusForbidden)
	})

	resp, err := client.DeleteDirectDebitMandate(context.Background(), "949404bd-d32e-4f1e-9759-4d6caee3137c")
	checkHasError(t, err)

	if resp.StatusCode != http.StatusForbidden {
		t.Error("should return HTTP 403 status")
	}
}
