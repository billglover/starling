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

var contactsTestCases = []struct {
	name string
	mock string
}{
	{
		name: "sample customer contacts",
		mock: `{
			"_links": {
				"self": {
					"href": "api/v1/contacts",
					"templated": false
				}
			},
			"_embedded": {
				"contacts": [
					{
						"_links": {
							"accounts": {
								"href": "api/v1/contacts/840e4030-b94c-4e71-a1d3-1319a233dd3c/accounts",
								"templated": false
							},
							"photo": {
								"href": "api/v1/contacts/840e4030-b94c-4e71-a1d3-1319a233dd3c/photo",
								"templated": false
							},
							"self": {
								"href": "api/v1/contacts/840e4030-b94c-4e71-a1d3-1319a233dd3c",
								"templated": false
							}
						},
						"id": "840e4030-b94c-4e71-a1d3-1319a233dd3c",
						"name": "Mickey Mouse"
					},
					{
						"_links": {
							"accounts": {
								"href": "api/v1/contacts/8a7d4b0c-e4a0-4687-86ae-2f859f75d17c/accounts",
								"templated": false
							},
							"photo": {
								"href": "api/v1/contacts/8a7d4b0c-e4a0-4687-86ae-2f859f75d17c/photo",
								"templated": false
							},
							"self": {
								"href": "api/v1/contacts/8a7d4b0c-e4a0-4687-86ae-2f859f75d17c",
								"templated": false
							}
						},
						"id": "8a7d4b0c-e4a0-4687-86ae-2f859f75d17c",
						"name": "Mickey Mouse"
					},
					{
						"_links": {
							"accounts": {
								"href": "api/v1/contacts/30c9a5e0-6bc0-49f7-960d-c240eee24bcc/accounts",
								"templated": false
							},
							"photo": {
								"href": "api/v1/contacts/30c9a5e0-6bc0-49f7-960d-c240eee24bcc/photo",
								"templated": false
							},
							"self": {
								"href": "api/v1/contacts/30c9a5e0-6bc0-49f7-960d-c240eee24bcc",
								"templated": false
							}
						},
						"id": "30c9a5e0-6bc0-49f7-960d-c240eee24bcc",
						"name": "Mickey Mouse"
					},
					{
						"_links": {
							"accounts": {
								"href": "api/v1/contacts/157e8e67-c642-427a-a62c-b978fb6a6f55/accounts",
								"templated": false
							},
							"photo": {
								"href": "api/v1/contacts/157e8e67-c642-427a-a62c-b978fb6a6f55/photo",
								"templated": false
							},
							"self": {
								"href": "api/v1/contacts/157e8e67-c642-427a-a62c-b978fb6a6f55",
								"templated": false
							}
						},
						"id": "157e8e67-c642-427a-a62c-b978fb6a6f55",
						"name": "Mickey Mouse"
					}
				]
			}
		}`,
	},
}

func TestGetContacts(t *testing.T) {

	t.Log("Given the need to test fetching customer contacts:")

	for _, tc := range contactsTestCases {
		t.Run(tc.name, func(st *testing.T) {
			testGetContacts(st, tc.name, tc.mock)
		})
	}
}

func testGetContacts(t *testing.T, name, mock string) {
	t.Logf("\tWhen making a call to GetContacts() with %s:", name)

	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/contacts", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		fmt.Fprint(w, mock)
	})

	got, _, err := client.GetContacts(context.Background())
	checkNoError(t, err)

	hal := &HALContacts{}
	json.Unmarshal([]byte(mock), hal)
	want := hal.Embedded

	if !reflect.DeepEqual(got, want) {
		t.Error("\t\tshould return a list of contacts matching the mock response", cross)
	} else {
		t.Log("\t\tshould return a list of contacts matching the mock response", tick)
	}

	if len(got.Contacts) == 0 {
		t.Errorf("\t\tshould have at least one contact %s %d", cross, len(got.Contacts))
	} else {
		t.Log("\t\tshould have at least one contact", tick)
	}
}

var contactTestCases = []struct {
	name string
	uid  string
	mock string
}{
	{
		name: "sample customer contact",
		uid:  "840e4030-b94c-4e71-a1d3-1319a233dd3c",
		mock: `{
			"_links": {
				"accounts": {
					"href": "api/v1/contacts/840e4030-b94c-4e71-a1d3-1319a233dd3c/accounts",
					"templated": false
				},
				"photo": {
					"href": "api/v1/contacts/840e4030-b94c-4e71-a1d3-1319a233dd3c/photo",
					"templated": false
				},
				"self": {
					"href": "api/v1/contacts/840e4030-b94c-4e71-a1d3-1319a233dd3c",
					"templated": false
				}
			},
			"id": "840e4030-b94c-4e71-a1d3-1319a233dd3c",
			"name": "Mickey Mouse"
		}`,
	},
	{
		name: "sample customer contact without HAL links",
		uid:  "840e4030-b94c-4e71-a1d3-1319a233dd3c",
		mock: `{
			"id": "840e4030-b94c-4e71-a1d3-1319a233dd3c",
			"name": "Mickey Mouse"
		}`,
	},
}

func TestGetContact(t *testing.T) {

	t.Log("Given the need to test fetching an individual customer contact:")

	for _, tc := range contactTestCases {
		t.Run(tc.name, func(st *testing.T) {
			testGetContact(st, tc.name, tc.uid, tc.mock)
		})
	}
}

func testGetContact(t *testing.T, name, uid, mock string) {
	t.Logf("\tWhen making a call to GetContact() with a %s:", name)

	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/contacts/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)

		reqUID := path.Base(r.URL.Path)
		if reqUID != uid {
			t.Error("\t\tshould send a requestwith the correct UID", cross, reqUID)
		} else {
			t.Log("\t\tshould send a request with the correct UID", tick)
		}

		fmt.Fprint(w, mock)
	})

	got, _, err := client.GetContact(context.Background(), uid)
	checkNoError(t, err)

	t.Log("\tWhen parsing the response from the API:")

	want := &Contact{}
	json.Unmarshal([]byte(mock), want)

	if !reflect.DeepEqual(got, want) {
		t.Error("\t\tshould return a single contact matching the mock response", cross)
	} else {
		t.Log("\t\tshould return a single contact matching the mock response", tick)
	}

	if got.UID != want.UID {
		t.Error("\t\tshould have the correct UID", cross, got.UID)
	} else {
		t.Log("\t\tshould have the correct UID", tick)
	}
}

var deleteContactTestCases = []struct {
	name string
	uid  string
	mock string
}{
	{
		name: "sample customer contact",
		uid:  "840e4030-b94c-4e71-a1d3-1319a233dd3c",
	},
}

func TestDeleteContact(t *testing.T) {

	t.Log("Given the need to test deleting an individual customer contact:")

	for _, tc := range deleteContactTestCases {
		t.Run(tc.name, func(st *testing.T) {
			testDeleteContact(st, tc.name, tc.uid)
		})
	}
}

func testDeleteContact(t *testing.T, name, uid string) {
	t.Logf("\tWhen making a call to DeleteContact() with a %s:", name)

	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/contacts/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodDelete)

		reqUID := path.Base(r.URL.Path)
		if reqUID != uid {
			t.Error("\t\tshould send a requestwith the correct UID", cross, reqUID)
		} else {
			t.Log("\t\tshould send a request with the correct UID", tick)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.DeleteContact(context.Background(), uid)
	checkNoError(t, err)

	t.Log("\tWhen parsing the response from the API:")

	if resp.StatusCode != http.StatusNoContent {
		t.Error("\t\tshould return an HTTP 204 status", cross, resp.Status)
	} else {
		t.Log("\t\tshould return an HTTP 204 status", tick)
	}

	body, err := ioutil.ReadAll(resp.Body)
	checkNoError(t, err)

	if len(body) != 0 {
		t.Error("\t\tshould return an empty body", cross, len(body))
	} else {
		t.Log("\t\tshould return an empty body", tick)
	}
}
