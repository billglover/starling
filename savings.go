package starling

import (
	"context"
	"net/http"

	"github.com/google/uuid"
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

// SavingsGoalRequest is a request to create a new savings goal
type SavingsGoalRequest struct {
	Name               string            `json:"name"`     // Name of the savings goal
	Currency           string            `json:"currency"` // ISO-4217 3 character currency code of the savings goal
	Target             CurrencyAndAmount `json:"target"`
	Base64EncodedPhoto string            `json:"base64EncodedPhoto"` // A text (base 64) encoded picture to associate with the savings goal
}

// CreateOrUpdateSavingsGoalResponse represents a response after attempting to create a savings goal
type CreateOrUpdateSavingsGoalResponse struct {
	UID     string        `json:"savingsGoalUid"`
	Success bool          `json:"success"`
	Errors  []ErrorDetail `json:"errors"`
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
		return goals, resp, err
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
		return goal, resp, err
	}

	return goal, resp, nil
}

// PutSavingsGoal creates an individual savings goal based on a UID. It returns the http response
// in case this is required for further processing. An error will be returned if the API is unable
// to create the goal.
func (c *Client) PutSavingsGoal(ctx context.Context, uid string, sgreq SavingsGoalRequest) (*CreateOrUpdateSavingsGoalResponse, *http.Response, error) {
	req, err := c.NewRequest("PUT", "/api/v1/savings-goals/"+uid, sgreq)
	if err != nil {
		return nil, nil, err
	}

	var sgresp *CreateOrUpdateSavingsGoalResponse
	resp, err := c.Do(ctx, req, &sgresp)
	if err != nil {
		return sgresp, resp, err
	}

	return sgresp, resp, nil
}

// AddMoney transfers money into a savings goal. It returns the http response in case this is required for further
// processing. An error will be returned if the API is unable to transfer the amount into the savings goal.
func (c *Client) AddMoney(ctx context.Context, goalUID string, tuReq TopUpRequest) (*SavingsGoalTransferResponse, *http.Response, error) {
	txnUID, err := uuid.NewRandom()
	if err != nil {
		return nil, nil, err
	}

	req, err := c.NewRequest("PUT", "/api/v1/savings-goals/"+goalUID+"/add-money/"+txnUID.String(), tuReq)
	if err != nil {
		return nil, nil, err
	}

	var tuResp *SavingsGoalTransferResponse
	resp, err := c.Do(ctx, req, &tuResp)
	if err != nil {
		return tuResp, resp, err
	}
	return tuResp, resp, nil
}
