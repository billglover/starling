package starling

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var userTestCases = []struct {
	name string
	mock string
}{
	{
		name: "sample user",
		mock: `{
		"customerUid": "6d2aa528-b9d1-4083-ae7c-53d460cd8d88",
		"expiresAt": "2018-04-03T09:33:01.137Z",
		"authenticated": true,
		"expiresInSeconds": 86373,
		"scopes": [
			"account:read",
			"balance:read",
			"address:read",
			"address:edit",
			"card:read",
			"customer:read",
			"mandate:read",
			"mandate:delete",
			"metadata:create",
			"metadata:edit",
			"payee:create",
			"payee:delete",
			"payee:edit",
			"payee:read",
			"pay-local:create",
			"pay-foreign:create",
			"transaction:read",
			"transaction:edit",
			"savings-goal:read",
			"savings-goal:create",
			"savings-goal:delete",
			"savings-goal-transfer:read",
			"savings-goal-transfer:create",
			"savings-goal-transfer:delete"
		]
	}`,
	},
	{
		name: "sample user and no scopes",
		mock: `{
		"customerUid": "6d2aa528-b9d1-4083-ae7c-53d460cd8d88",
		"expiresAt": "2018-04-03T09:33:01.137Z",
		"authenticated": true,
		"expiresInSeconds": 86373,
		"scopes": []
	}`,
	},
}

func TestCurrentUser(t *testing.T) {

	t.Log("Given the need to test fetching details on the current user:")

	// Run each of the test cases a subtest.
	for _, tc := range userTestCases {
		t.Run(tc.name, func(st *testing.T) {
			testCurrentUser(st, tc.name, tc.mock)
		})
	}
}

func testCurrentUser(t *testing.T, name, mock string) {
	t.Logf("\tWhen making a call to CurrentUser() with %s:", name)

	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/me", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, http.MethodGet)
		fmt.Fprint(w, mock)
	})

	got, _, err := client.CurrentUser(context.Background())
	checkNoError(t, err)

	want := &Identity{}
	json.Unmarshal([]byte(mock), want)

	if !reflect.DeepEqual(got, want) {
		t.Error("\t\tshould return an identity matching the mock response", cross)
	} else {
		t.Log("\t\tshould return an identity matching the mock response", tick)
	}
}
