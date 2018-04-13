package starling

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

// TestSavingsGoals confirms that the client is able to query a list of savings goals.
func TestSavingsGoals(t *testing.T) {

	t.Log("Given the need to test fetching savings goals:")

	client, mux, _, teardown := setup()
	defer teardown()

	mock := `{
		"savingsGoalList": [
		  {
			"uid": "e43d3060-2c83-4bb9-ac8c-c627b9c45f8b",
			"name": "Trip to Paris",
			"target": {
			  "currency": "GBP",
			  "minorUnits": 11223344
			},
			"totalSaved": {
			  "currency": "GBP",
			  "minorUnits": 11223344
			},
			"savedPercentage": 50
		  }
		]
	  }`

	t.Log("\tWhen making a call to GetSavingsGoals():")

	mux.HandleFunc("/api/v1/savings-goals", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		fmt.Fprint(w, mock)
	})

	got, _, err := client.SavingsGoals(context.Background())
	if err != nil {
		t.Fatal("\t\tshould be able to make the request", cross, err)
	} else {
		t.Log("\t\tshould be able to make the request", tick)
	}

	t.Log("\tWhen parsing the response from the API:")

	want := &savingsGoals{}
	json.Unmarshal([]byte(mock), want)

	if got == nil {
		t.Fatal("\t\tshould not return 'nil'", cross)
	} else {
		t.Log("\t\tshould not return 'nil'", tick)
	}

	if !reflect.DeepEqual(got, &want.SavingsGoals) {
		t.Error("\t\tshould return a list of savings goals matching the mock response", cross)
	} else {
		t.Log("\t\tshould return a list of savings goals matching the mock response", tick)
	}

	if len(*got) == 0 {
		t.Errorf("\t\tshould return a list with at least one savings goal %s %d", cross, len(*got))
	} else {
		t.Log("\t\tshould return a list with at least one savings goal", tick)
	}

}

// TestGetSavingsGoals confirms that the client is able to query a single savings goal.
func TestGetSavingsGoal(t *testing.T) {

	t.Log("Given the need to test fetching an individual savings goal:")

	client, mux, _, teardown := setup()
	defer teardown()

	// Define our mock response and handler
	mock := `{
		"uid": "e43d3060-2c83-4bb9-ac8c-c627b9c45f8b",
		"name": "Trip to Paris",
		"target": {
			"currency": "GBP",
			"minorUnits": 11223344
		},
		"totalSaved": {
			"currency": "GBP",
			"minorUnits": 11223344
		},
		"savedPercentage": 50
	}`

	t.Log("\tWhen making a call to GetSavingsGoal():")

	mux.HandleFunc("/api/v1/savings-goals/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "GET")
		fmt.Fprint(w, mock)
	})

	got, _, err := client.SavingsGoal(context.Background(), "e43d3060-2c83-4bb9-ac8c-c627b9c45f8b")
	if err != nil {
		t.Fatal("\t\tshould be able to make the request", cross, err)
	} else {
		t.Log("\t\tshould be able to make the request", tick)
	}

	t.Log("\tWhen parsing the response from the API:")

	want := &SavingsGoal{}
	json.Unmarshal([]byte(mock), want)

	if got == nil {
		t.Fatal("\t\tshould not return 'nil'", cross)
	} else {
		t.Log("\t\tshould not return 'nil'", tick)
	}

	if !reflect.DeepEqual(got, want) {
		t.Error("\t\tshould return a savings goal that matches the mock response", cross)
	} else {
		t.Log("\t\tshould return a savings goal that matches the mock response", tick)
	}

	if len(got.UID) == 0 {
		t.Errorf("\t\tshould return a savings goal that has a UID %s %d", cross, len(got.UID))
	} else {
		t.Log("\t\tshould return a savings goal that has a UID", tick)
	}

}

// TestPutSavingsGoal confirms that the client is able to create a savings goal.
func TestPutSavingsGoal(t *testing.T) {

	t.Log("Given the need to test creating a savings goal:")

	client, mux, _, teardown := setup()
	defer teardown()

	uid := "e43d3060-2c83-4bb9-ac8c-c627b9c45f8b"
	mockReq := SavingsGoalRequest{
		Name:     "test",
		Currency: "GBP",
		Target: Amount{
			Currency:   "GBP",
			MinorUnits: 10000,
		},
		Base64EncodedPhoto: "",
	}
	mockResp := `{
		"savingsGoalUid": "e43d3060-2c83-4bb9-ac8c-c627b9c45f8b",
		"success": true,
		"errors": []
	  }`

	t.Log("\tWhen making a call to PutSavingsGoal():")

	mux.HandleFunc("/api/v1/savings-goals/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "PUT")

		var sg = SavingsGoalRequest{}
		err := json.NewDecoder(r.Body).Decode(&sg)
		if err != nil {
			t.Fatal("\t\tshould send a request that the API can parse", cross, err)
		} else {
			t.Log("\t\tshould send a request that the API can parse", tick)
		}

		if !reflect.DeepEqual(mockReq, sg) {
			t.Error("\t\tshould send a savings goal that matches the mock", cross)
		} else {
			t.Log("\t\tshould send a savings goal that matches the mock", tick)
		}

		reqUID := path.Base(r.URL.Path)
		if reqUID != uid {
			t.Error("\t\tshould send a savings goal with the correct UID", cross, reqUID)
		} else {
			t.Log("\t\tshould send a savings goal with the correct UID", tick)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, mockResp)
	})

	resp, err := client.CreateSavingsGoal(context.Background(), uid, mockReq)
	if err != nil {
		t.Fatal("\t\tshould be able to make the request", cross, err)
	} else {
		t.Log("\t\tshould be able to make the request", tick)
	}

	t.Log("\tWhen parsing the response from the API:")

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("\t\tshould receive a %d status code %s %d", want, cross, got)
	} else {
		t.Logf("\t\tshould receive a %d status code %s", want, tick)
	}
}

// TestPutSavingsGoal_ValidateError confirms that the client is able to handle validation
// errors when creating savings goals.
func TestPutSavingsGoal_ValidateError(t *testing.T) {

	t.Log("Given the need to test handling validation errors when creating a savings goal:")

	client, mux, _, teardown := setup()
	defer teardown()

	uid := "d8770f9d-4ee9-4cc1-86e1-83c26bcfcc4f"
	mockReq := SavingsGoalRequest{
		Name:     "test",
		Currency: "GBP",
		Target: Amount{
			Currency:   "GBP",
			MinorUnits: 10000,
		},
		Base64EncodedPhoto: "",
	}
	mockResp := `{
		"savingsGoalUid": "d8770f9d-4ee9-4cc1-86e1-83c26bcfcc4f",
		"success": true,
		"errors": [
			{
				"message": "Something about the validation error"
			}
		]
	}`

	t.Log("\tWhen making a call to PutSavingsGoal():")

	mux.HandleFunc("/api/v1/savings-goals/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "PUT")

		var sg = SavingsGoalRequest{}
		err := json.NewDecoder(r.Body).Decode(&sg)
		if err != nil {
			t.Fatal("\t\tshould send a request that the API can parse", cross, err)
		} else {
			t.Log("\t\tshould send a request that the API can parse", tick)
		}

		if !reflect.DeepEqual(mockReq, sg) {
			t.Error("\t\tshould send a savings goal that matches the mock", cross)
		} else {
			t.Log("\t\tshould send a savings goal that matches the mock", tick)
		}

		reqUID := path.Base(r.URL.Path)
		if reqUID != uid {
			t.Error("\t\tshould send a savings goal with the correct UID", cross, reqUID)
		} else {
			t.Log("\t\tshould send a savings goal with the correct UID", tick)
		}

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, mockResp)
	})

	resp, err := client.CreateSavingsGoal(context.Background(), uid, mockReq)
	if err == nil {
		t.Fatal("\t\texpected an error to be returned", cross)
	} else {
		t.Log("\t\texpected an error to be returned", tick)
	}

	t.Log("\tWhen parsing the response from the API:")

	if got, want := resp.StatusCode, http.StatusBadRequest; got != want {
		t.Errorf("\t\tshould receive a %d status code %s %d", want, cross, got)
	} else {
		t.Logf("\t\tshould receive a %d status code %s", want, tick)
	}
}

// TestAddMoney confirms that the client is able to make a request to add money to a savings goal and parse
// the successful response from the API.
func TestAddMoney(t *testing.T) {

	t.Log("Given the need to test adding money to a savings goal:")

	client, mux, _, teardown := setup()
	defer teardown()

	goalUID := "d8770f9d-4ee9-4cc1-86e1-83c26bcfcc4f"
	txnUID := "28dff346-dd48-426f-96df-d7f33d29c379"
	mockResp := `{"transferUid":"28dff346-dd48-426f-96df-d7f33d29c379","success":true,"errors":[]}`

	mockReq := TopUpRequest{
		Amount: Amount{
			Currency:   "GBP",
			MinorUnits: 1050,
		},
	}

	t.Log("\tWhen making a call to AddMoney():")

	mux.HandleFunc("/api/v1/savings-goals/", func(w http.ResponseWriter, r *http.Request) {
		checkMethod(t, r, "PUT")

		var tu = TopUpRequest{}
		err := json.NewDecoder(r.Body).Decode(&tu)
		if err != nil {
			t.Fatal("\t\tshould send a request that the API can parse", cross, err)
		} else {
			t.Log("\t\tshould send a request that the API can parse", tick)
		}

		if !reflect.DeepEqual(mockReq, tu) {
			t.Error("\t\tshould send a top-up request that matches the mock", cross)
		} else {
			t.Log("\t\tshould send a top-up request that matches the mock", tick)
		}

		reqUID, err := uuid.Parse(path.Base(r.URL.Path))
		if err != nil {
			t.Error("\t\tshould send a top-up request with a valid UID", cross, reqUID)
		} else {
			t.Log("\t\tshould send a top-up request with a valid UID", tick)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, mockResp)
	})

	tuResp, resp, err := client.AddMoney(context.Background(), goalUID, mockReq)
	if err != nil {
		t.Fatal("\t\tshould be able to make the request", cross, err)
	} else {
		t.Log("\t\tshould be able to make the request", tick)
	}

	t.Log("\tWhen parsing the response from the API:")

	want := &SavingsGoalTransferResponse{}
	json.Unmarshal([]byte(mockResp), want)

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("\t\tshould receive a %d status code %s %d", want, cross, got)
	} else {
		t.Logf("\t\tshould receive a %d status code %s", want, tick)
	}

	if got, want := tuResp.UID, txnUID; got != want {
		t.Fatal("\t\tshould be receive the UID assigned to the transaction", cross, got)
	} else {
		t.Log("\t\tshould be receive the UID assigned to the transaction", tick)
	}

	if got, want := tuResp.Success, true; got != want {
		t.Fatal("\t\tshould be receive a success status", cross, got)
	} else {
		t.Log("\t\tshould be receive a success status", tick)
	}

	if got, want := len(tuResp.Errors), 0; got != want {
		t.Fatal("\t\tshould be receive no validation errors", cross, got)
	} else {
		t.Log("\t\tshould be receive no validation errors", tick)
	}
}
