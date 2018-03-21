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
	client, mux, _, teardown := setup()
	defer teardown()

	json := `{
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

	mux.HandleFunc("/api/v1/savings-goals", func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, "GET"; got != want {
			t.Errorf("request method: %v, want %v", got, want)
		}
		fmt.Fprint(w, json)
	})

	goals, _, err := client.GetSavingsGoals(context.Background())
	if err != nil {
		t.Error("unexpected error returned:", err)
	}

	want := &SavingsGoals{
		SavingsGoalList: []SavingsGoal{
			SavingsGoal{
				UID:  "e43d3060-2c83-4bb9-ac8c-c627b9c45f8b",
				Name: "Trip to Paris",
				Target: CurrencyAndAmount{
					Currency:   "GBP",
					MinorUnits: 11223344,
				},
				TotalSaved: CurrencyAndAmount{
					Currency:   "GBP",
					MinorUnits: 11223344,
				},
				SavedPercentage: 50,
			},
		},
	}

	if !reflect.DeepEqual(goals, want) {
		t.Errorf("GetSavingsGoals returned %+v, want %+v", goals, want)
	}

}

// TestGetSavingsGoals_Error confirms that the client is able to parse a successful error response from the API.
// It confirms that error messages are decoded from the JSON description and returned as the error string. It
// also confirms that nil is returned for the savings goals in the event of an error.
func TestGetSavingsGoals_Error(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	msg := ErrorDetail{
		Message: "this is an error message",
	}

	mux.HandleFunc("/api/v1/savings-goals", func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, "GET"; got != want {
			t.Errorf("request method: %v, want %v", got, want)
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)
	})

	goals, _, err := client.GetSavingsGoals(context.Background())
	if err == nil {
		t.Error("expected an error to be returned:")
	}

	if err.Error() != msg.Message {
		t.Errorf("GetSavingsGoals returned '%v', want '%v'", err, msg.Message)
	}

	if goals != nil {
		t.Errorf("unexpected goals returned: %+v", goals)
	}

}

// TestGetSavingsGoals confirms that the client is able to query a single savings goal.
func TestGetSavingsGoal(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	json := `{
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

	mux.HandleFunc("/api/v1/savings-goals/", func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, "GET"; got != want {
			t.Errorf("request method: %v, want %v", got, want)
		}
		fmt.Fprint(w, json)
	})

	goal, _, err := client.GetSavingsGoal(context.Background(), "e43d3060-2c83-4bb9-ac8c-c627b9c45f8b")
	if err != nil {
		t.Error("unexpected error returned:", err)
	}

	want := &SavingsGoal{
		UID:  "e43d3060-2c83-4bb9-ac8c-c627b9c45f8b",
		Name: "Trip to Paris",
		Target: CurrencyAndAmount{
			Currency:   "GBP",
			MinorUnits: 11223344,
		},
		TotalSaved: CurrencyAndAmount{
			Currency:   "GBP",
			MinorUnits: 11223344,
		},
		SavedPercentage: 50,
	}

	if !reflect.DeepEqual(goal, want) {
		t.Errorf("GetSavingsGoal returned %+v, want %+v", goal, want)
	}

}

// TestGetSavingsGoal_Error confirms that the client is able to parse a successful error response from the API.
// It confirms that error messages are decoded from the JSON description and returned as the error string. It
// also confirms that nil is returned for the savings goals in the event of an error.
func TestGetSavingsGoal_Error(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	msg := ErrorDetail{
		Message: "this is an error message",
	}

	mux.HandleFunc("/api/v1/savings-goals/", func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Method, "GET"; got != want {
			t.Errorf("request method: %v, want %v", got, want)
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)
	})

	goal, _, err := client.GetSavingsGoal(context.Background(), "e43d3060-2c83-4bb9-ac8c-c627b9c45f8b")
	if err == nil {
		t.Error("expected an error to be returned:")
	}

	if err.Error() != msg.Message {
		t.Errorf("GetSavingsGoals returned '%v', want '%v'", err, msg.Message)
	}

	if goal != nil {
		t.Errorf("unexpected goal returned: %+v", goal)
	}

}
