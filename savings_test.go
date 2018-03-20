package starling

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

// TestGetSavingsGoals confirms that the client is able to parse a successful response from the API
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
