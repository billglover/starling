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

func TestContacts(t *testing.T) {
	for _, tc := range contactsTestCases {
		t.Run(tc.name, func(st *testing.T) {
			testContacts(st, tc.name, tc.mock)
		})
	}
}

func testContacts(t *testing.T, name, mock string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/contacts", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		fmt.Fprint(w, mock)
	})

	got, _, err := client.Contacts(context.Background())
	checkNoError(t, err)

	hal := &halContacts{}
	json.Unmarshal([]byte(mock), hal)
	want := hal.Embedded

	if !reflect.DeepEqual(got, want.Contacts) {
		t.Error("should return a list of contacts matching the mock response", cross)
	}

	if len(got) == 0 {
		t.Errorf("should have at least one contact %s %d", cross, len(got))
	}
}

func TestContactAccountsError(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/contacts", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusForbidden)
	})

	got, resp, err := client.Contacts(context.Background())
	checkNoError(t, err)

	if resp.StatusCode != http.StatusForbidden {
		t.Error("should return HTTP 403 status")
	}

	if len(got) != 0 {
		t.Error("should not return any contacts")
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

func TestContact(t *testing.T) {
	for _, tc := range contactTestCases {
		t.Run(tc.name, func(st *testing.T) {
			testContact(st, tc.name, tc.uid, tc.mock)
		})
	}
}

func testContact(t *testing.T, name, uid, mock string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/contacts/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)

		reqUID := path.Base(r.URL.Path)
		if reqUID != uid {
			t.Error("should send a requestwith the correct UID", cross, reqUID)
		}

		fmt.Fprint(w, mock)
	})

	got, _, err := client.Contact(context.Background(), uid)
	checkNoError(t, err)

	want := &Contact{}
	json.Unmarshal([]byte(mock), want)

	if !reflect.DeepEqual(got, want) {
		t.Error("should return a single contact matching the mock response", cross)
	}

	if got.UID != want.UID {
		t.Error("should have the correct UID", cross, got.UID)
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
	for _, tc := range deleteContactTestCases {
		t.Run(tc.name, func(st *testing.T) {
			testDeleteContact(st, tc.name, tc.uid)
		})
	}
}

func testDeleteContact(t *testing.T, name, uid string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/contacts/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodDelete)

		reqUID := path.Base(r.URL.Path)
		if reqUID != uid {
			t.Error("should send a requestwith the correct UID", cross, reqUID)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.DeleteContact(context.Background(), uid)
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

var contactAccountsTestCases = []struct {
	name string
	uid  string
	mock string
}{
	{
		name: "sample customer contact accounts",
		uid:  "8a7d4b0c-e4a0-4687-86ae-2f859f75d17c",
		mock: `{
			"self": {
				"href": "api/v1/contacts/8a7d4b0c-e4a0-4687-86ae-2f859f75d17c/accounts",
				"templated": false
			},
			"contactAccounts": [
				{
					"self": {
						"href": "api/v1/contacts/8a7d4b0c-e4a0-4687-86ae-2f859f75d17c/accounts/64834e9a-a920-4329-b28d-24246d332f83",
						"templated": false
					},
					"id": "64834e9a-a920-4329-b28d-24246d332f83",
					"type": "UK_ACCOUNT_AND_SORT_CODE",
					"name": "UK account",
					"accountNumber": "00000825",
					"sortCode": "204514"
				}
			]
		}`,
	},
	{
		name: "sample customer contact accounts without HAL links",
		uid:  "8a7d4b0c-e4a0-4687-86ae-2f859f75d17c",
		mock: `{
			"contactAccounts": [
				{
					"id": "64834e9a-a920-4329-b28d-24246d332f83",
					"type": "UK_ACCOUNT_AND_SORT_CODE",
					"name": "UK account",
					"accountNumber": "00000825",
					"sortCode": "204514"
				}
			]
		}`,
	},
}

func TestContactAccounts(t *testing.T) {
	for _, tc := range contactAccountsTestCases {
		t.Run(tc.name, func(st *testing.T) {
			testContactAccounts(st, tc.name, tc.mock, tc.uid)
		})
	}
}

func testContactAccounts(t *testing.T, name, mock, uid string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/contacts/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)

		reqUID := path.Base(path.Dir(r.URL.Path))

		if reqUID != uid {
			t.Error("should send a request with the correct UID", cross, reqUID)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, mock)
	})

	got, _, err := client.ContactAccounts(context.Background(), uid)
	checkNoError(t, err)

	want := &contactAccounts{}
	json.Unmarshal([]byte(mock), want)

	if !reflect.DeepEqual(got, want.ContactAccounts) {
		t.Error("should return a list of contact accounts matching the mock response", cross)
	}

	if len(got) == 0 {
		t.Errorf("should have at least one contact account %s %d", cross, len(got))
	}
}

var contactAccountTestCases = []struct {
	name       string
	contactUID string
	accountUID string
	mock       string
}{
	{
		name:       "sample customer contact account",
		contactUID: "8a7d4b0c-e4a0-4687-86ae-2f859f75d17c",
		accountUID: "64834e9a-a920-4329-b28d-24246d332f83",
		mock: `{
			"self": {
				 "href": "api/v1/contacts/8a7d4b0c-e4a0-4687-86ae-2f859f75d17c/accounts/64834e9a-a920-4329-b28d-24246d332f83",
				 "templated": false
			},
			"id": "64834e9a-a920-4329-b28d-24246d332f83",
			"type": "UK_ACCOUNT_AND_SORT_CODE",
			"name": "UK account",
			"accountNumber": "00000825",
			"sortCode": "204514"
	  }`,
	},
	{
		name:       "sample customer contact account without HAL links",
		contactUID: "8a7d4b0c-e4a0-4687-86ae-2f859f75d17c",
		accountUID: "64834e9a-a920-4329-b28d-24246d332f83",
		mock: `{
			"id": "64834e9a-a920-4329-b28d-24246d332f83",
			"type": "UK_ACCOUNT_AND_SORT_CODE",
			"name": "UK account",
			"accountNumber": "00000825",
			"sortCode": "204514"
	  }`,
	},
}

func TestContactAccount(t *testing.T) {
	for _, tc := range contactAccountTestCases {
		t.Run(tc.name, func(st *testing.T) {
			testContactAccount(st, tc.name, tc.mock, tc.contactUID, tc.accountUID)
		})
	}
}

func testContactAccount(t *testing.T, name, mock, cUID, aUID string) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/contacts/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)

		reqAccountUID := path.Base(r.URL.Path)

		if reqAccountUID != aUID {
			t.Error("should send a request with the correct UID", cross, reqAccountUID)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, mock)
	})

	got, _, err := client.ContactAccount(context.Background(), cUID, aUID)
	checkNoError(t, err)

	want := &ContactAccount{}
	json.Unmarshal([]byte(mock), want)

	if got.AccountNumber == "" {
		t.Error("should have an account number", cross)
	}

	if got.SortCode == "" {
		t.Error("should have an sort code", cross)
	}

	if !reflect.DeepEqual(got, want) {
		t.Error("should return a contact account matching the mock response", cross)
	}
}

var createContactAcctTestCases = []struct {
	name       string
	contAct    ContactAccount
	respBody   string
	respStatus int
}{
	{
		name: "sample customer contact account",
		contAct: ContactAccount{
			UID:           "8cdab926-1d16-46a7-b4a9-6cb38f0c9b49",
			Name:          "Dave Bowman",
			Type:          "UK_ACCOUNT_AND_SORT_CODE",
			AccountNumber: "70872490",
			SortCode:      "404784",
		},
		respBody:   "",
		respStatus: http.StatusCreated,
	},
	{
		name: "sample customer contact account",
		contAct: ContactAccount{
			UID:           "8cdab926-1d16-46a7-b4a9-6cb38f0c9b49",
			Name:          "Dave Bowman",
			Type:          "UK_ACCOUNT_AND_SORT_CODE",
			AccountNumber: "12345678",
			SortCode:      "404784",
		},
		respBody: `[
    "INVALID_SORT_CODE_OR_ACCOUNT_NUMBER"
]`,
		respStatus: http.StatusBadRequest,
	},
}

func TestPostContactAccount(t *testing.T) {
	for _, tc := range createContactAcctTestCases {
		t.Run(tc.name, func(st *testing.T) {
			testPostContactAccount(st, tc.name, tc.contAct, tc.respBody, tc.respStatus)
		})
	}
}

func testPostContactAccount(t *testing.T, name string, ca ContactAccount, respBody string, respStatus int) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/contacts", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodPost)

		var reqCA = ContactAccount{}
		err := json.NewDecoder(r.Body).Decode(&reqCA)
		if err != nil {
			t.Fatal("should send a request that the API can parse", cross, err)
		}

		if !reflect.DeepEqual(ca, reqCA) {
			t.Error("should send a contact account that matches the mock", cross)
		}
		w.Header().Set("Location", reqCA.UID)
		w.WriteHeader(respStatus)
		fmt.Fprintln(w, respBody)
	})

	uid, resp, err := client.CreateContactAccount(context.Background(), ca)
	if respStatus <= 299 {
		checkNoError(t, err)
	} else {
		checkHasError(t, err)
	}

	if uid != ca.UID {
		t.Error("should return the UID in the Location header")
	}

	checkStatus(t, resp, respStatus)
}
