package starling

import (
	"context"
	"net/http"
)

// SavingsGoal is a goal defined by a customer to hold savings
type SavingsGoal struct {
	UID             string            `json:"uid"`  // Unique identifier of the savings goal
	Name            string            `json:"name"` // Name of the savings goal
	Target          CurrencyAndAmount `json:"target"`
	TotalSaved      CurrencyAndAmount `json:"totalSaved"`
	SavedPercentage int32             `json:"savedPercentage"` // Percentage of target currently deposited in the savings goal
}

// SavingsGoals is a list containing all savings goals for customer
type SavingsGoals struct {
	SavingsGoalList []SavingsGoal `json:"savingsGoalList"`
}

// GetSavingsGoals returns the savings goals for the current user. It also returns the http response
// in case this is required for further processing. It is possible that the user has no savings goals
// in which case a nil value will be returned. An error will be returned if unable to retrieve goals
// from the API.
func (c *Client) GetSavingsGoals(ctx context.Context) (*SavingsGoals, *http.Response, error) {
	req, err := c.NewRequest("GET", "/api/v1/savings-goals", nil)
	if err != nil {
		return nil, nil, err
	}

	var goals *SavingsGoals
	resp, err := c.Do(ctx, req, &goals)
	if err != nil {
		return nil, resp, err
	}

	return goals, resp, nil
}

// GetSavingsGoal returns an individual savings goal based on a UID. It also returns the http response
// in case this is required for further processing. An error will be returned if unable to retrieve
// goals from the API.
func (c *Client) GetSavingsGoal(ctx context.Context, uid string) (*SavingsGoal, *http.Response, error) {
	req, err := c.NewRequest("GET", "/api/v1/savings-goals/"+uid, nil)
	if err != nil {
		return nil, nil, err
	}

	var goal *SavingsGoal
	resp, err := c.Do(ctx, req, &goal)
	if err != nil {
		return nil, resp, err
	}

	return goal, resp, nil
}
