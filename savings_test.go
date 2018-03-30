package starling

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

// TestGetSavingsGoals confirms that the client is able to query a list of savings goals.
func TestGetSavingsGoals(t *testing.T) {

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

	got, _, err := client.GetSavingsGoals(context.Background())
	if err != nil {
		t.Fatal("\t\tshould be able to make the request", cross, err)
	} else {
		t.Log("\t\tshould be able to make the request", tick)
	}

	t.Log("\tWhen parsing the response from the API:")

	want := &SavingsGoals{}
	json.Unmarshal([]byte(mock), want)

	if got == nil {
		t.Fatal("\t\tshould not return 'nil'", cross)
	} else {
		t.Log("\t\tshould not return 'nil'", tick)
	}

	if !reflect.DeepEqual(got, want) {
		t.Error("\t\tshould return a list matching the mock response", cross)
	} else {
		t.Log("\t\tshould return a transaction list matching the mock response", tick)
	}

	if len(got.SavingsGoalList) == 0 {
		t.Errorf("\t\tshould have at least one transaction %s %d", cross, len(got.SavingsGoalList))
	} else {
		t.Log("\t\tshould have at least one transaction", tick)
	}

}

// TestGetSavingsGoals confirms that the client is able to query a single savings goal.
func TestGetSavingsGoal(t *testing.T) {
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

	// Define our request and execute the tests
	got, _, err := client.GetSavingsGoal(context.Background(), "e43d3060-2c83-4bb9-ac8c-c627b9c45f8b")
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
		t.Error("\t\tshould return a list matching the mock response", cross)
	} else {
		t.Log("\t\tshould return a transaction list matching the mock response", tick)
	}

	if len(got.UID) == 0 {
		t.Errorf("\t\tshould have a UID %s %d", cross, len(got.UID))
	} else {
		t.Log("\t\tshould have a UID", tick)
	}

}

// TestPutSavingsGoal confirms that the client is able to create a savings goal.
func TestPutSavingsGoal(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	uid := "e43d3060-2c83-4bb9-ac8c-c627b9c45f8b"
	sgr := SavingsGoalRequest{
		Name:     "test",
		Currency: "GBP",
		Target: CurrencyAndAmount{
			Currency:   "GBP",
			MinorUnits: 10000,
		},
		Base64EncodedPhoto: "",
	}

	mux.HandleFunc("/api/v1/savings-goals/", func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, "PUT"; got != want {
			t.Errorf("request method: %v, want %v", got, want)
		}

		var sg = SavingsGoalRequest{}
		err := json.NewDecoder(r.Body).Decode(&sg)
		if err != nil {
			t.Errorf("unable to decode savings goal request: %v", err)
		}

		if !reflect.DeepEqual(sgr, sg) {
			t.Errorf("PutSavingsGoal got %+v, want %+v", sg, sgr)
		}

		w.WriteHeader(http.StatusOK)
	})

	_, resp, err := client.PutSavingsGoal(context.Background(), uid, sgr)
	if err != nil {
		t.Error("unexpected error returned:", err)
	}

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		t.Errorf("response status: %v, want %v", got, want)
	}

}

// TestPutSavingsGoal_ValidateError confirms that the client is able to handle validation
// errors when creating savings goals.
func TestPutSavingsGoal_ValidateError(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	// Define our mock response and handler
	mock := `{
		"savingsGoalUid": "d8770f9d-4ee9-4cc1-86e1-83c26bcfcc4f",
		"success": true,
		"errors": [
			{
				"message": "Something about the validation error"
			}
		]
	}`

	mux.HandleFunc("/api/v1/savings-goals/", func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, "PUT"; got != want {
			t.Errorf("request method: %v, want %v", got, want)
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, mock)
	})

	// Define our request and execute the tests
	sgr := SavingsGoalRequest{
		Name:     "test",
		Currency: "GBP",
		Target: CurrencyAndAmount{
			Currency:   "GBP",
			MinorUnits: 10000,
		},
		Base64EncodedPhoto: "",
	}

	got, _, err := client.PutSavingsGoal(context.Background(), "d8770f9d-4ee9-4cc1-86e1-83c26bcfcc4f", sgr)
	if err == nil {
		t.Error("expected an error to be returned")
	}

	want := CreateOrUpdateSavingsGoalResponse{}
	json.Unmarshal([]byte(mock), &want)

	if !reflect.DeepEqual(*got, want) {
		t.Errorf("GetSavingsGoal returned %+v, want %+v", got, want)
	}
}

// TestAddMoney confirms that the client is able to make a request to add money to a savings goal and parse
// the successful response from the API.
func TestAddMoney(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	sgUID := "28dff346-dd48-426f-96df-d7f33d29c379"
	mock := `{"transferUid":"28dff346-dd48-426f-96df-d7f33d29c379","success":true,"errors":[]}`

	tuReq := TopUpRequest{
		Amount: CurrencyAndAmount{
			Currency:   "GBP",
			MinorUnits: 1050,
		},
	}

	mux.HandleFunc("/api/v1/savings-goals/", func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, "PUT"; got != want {
			t.Errorf("request method: %v, want %v", got, want)
		}

		var tu = TopUpRequest{}
		err := json.NewDecoder(r.Body).Decode(&tu)
		if err != nil {
			t.Errorf("unable to decode top up request: %v", err)
		}

		if !reflect.DeepEqual(tu, tuReq) {
			t.Errorf("AddMoney got %+v, want %+v", tu, tuReq)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, mock)
	})

	got, _, err := client.AddMoney(context.Background(), sgUID, tuReq)
	if err != nil {
		t.Error("unexpected error returned:", err)
	}

	want := &SavingsGoalTransferResponse{}
	json.Unmarshal([]byte(mock), want)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("AddMoney returned \n%+v, want \n%+v", got, want)
	}
}

func checkMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("\t\tshould send a %s request to the API %s %s", want, cross, got)
	} else {
		t.Logf("\t\tshould send a %s request to the API %s", want, tick)
	}
}
