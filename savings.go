package starling

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// SavingsGoal is a goal defined by a customer to hold savings
type SavingsGoal struct {
	UID             string `json:"uid"`  // Unique identifier of the savings goal
	Name            string `json:"name"` // Name of the savings goal
	Target          Amount `json:"target"`
	TotalSaved      Amount `json:"totalSaved"`
	SavedPercentage int32  `json:"savedPercentage"` // Percentage of target currently deposited in the savings goal
}

// SavingsGoals is a list containing all savings goals for customer
type savingsGoals struct {
	SavingsGoals []SavingsGoal `json:"savingsGoalList"`
}

// SavingsGoalRequest is a request to create a new savings goal
type SavingsGoalRequest struct {
	Name               string `json:"name"`     // Name of the savings goal
	Currency           string `json:"currency"` // ISO-4217 3 character currency code of the savings goal
	Target             Amount `json:"target"`
	Base64EncodedPhoto string `json:"base64EncodedPhoto"` // A text (base 64) encoded picture to associate with the savings goal
}

// SavingsGoalResponse represents a response after attempting to create a savings goal
type savingsGoalResponse struct {
	UID     string        `json:"savingsGoalUid"`
	Success bool          `json:"success"`
	Errors  []ErrorDetail `json:"errors"`
}

// SavingsGoalTransferResponse represents the response received after attempting to make an immediate or recurring transfer
// into/out of a savings goal.
type savingsGoalTransferResponse struct {
	UID     string        `json:"transferUid"` // Unique identifier for the transfer
	Success bool          `json:"success"`     // True if the method completed successfully
	Errors  []ErrorDetail `json:"errors"`      // List of errors if the method request failed
}

// TopUpRequest represents request to make an immediate transfer into a savings goal
type topUpRequest struct {
	Amount `json:"amount"`
}

// SavingsGoals returns the savings goals for the current user. It also returns the http response
// in case this is required for further processing. It is possible that the user has no savings goals
// in which case a nil value will be returned. An error will be returned if unable to retrieve goals
// from the API.
func (c *Client) SavingsGoals(ctx context.Context) (*[]SavingsGoal, *http.Response, error) {
	req, err := c.NewRequest("GET", "/api/v1/savings-goals", nil)
	if err != nil {
		return nil, nil, err
	}

	var goals savingsGoals
	resp, err := c.Do(ctx, req, &goals)
	if err != nil {
		return &goals.SavingsGoals, resp, err
	}

	return &goals.SavingsGoals, resp, nil
}

// SavingsGoal returns an individual savings goal based on a UID. It also returns the http response
// in case this is required for further processing. An error will be returned if unable to retrieve
// goals from the API.
func (c *Client) SavingsGoal(ctx context.Context, uid string) (*SavingsGoal, *http.Response, error) {
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

// CreateSavingsGoal creates an individual savings goal based on a UID. It returns the http response
// in case this is required for further processing. An error will be returned if the API is unable
// to create the goal.
func (c *Client) CreateSavingsGoal(ctx context.Context, uid string, sgReq SavingsGoalRequest) (*http.Response, error) {
	req, err := c.NewRequest("PUT", "/api/v1/savings-goals/"+uid, sgReq)
	if err != nil {
		return nil, err
	}

	var sgResp *savingsGoalResponse
	resp, err := c.Do(ctx, req, &sgResp)
	if err != nil {
		return resp, err
	}

	ers := make([]string, len(sgResp.Errors))
	for i, v := range sgResp.Errors {
		ers[i] = v.Message
	}

	if sgResp.Success != true {
		return resp, fmt.Errorf(strings.Join(ers, ", "))
	}

	return resp, nil
}

// AddMoney transfers money into a savings goal. It returns the http response in case this is required for further
// processing. An error will be returned if the API is unable to transfer the amount into the savings goal.
func (c *Client) AddMoney(ctx context.Context, goalUID string, a Amount) (string, *http.Response, error) {
	txnUID, err := uuid.NewRandom()
	if err != nil {
		return "", nil, err
	}

	req, err := c.NewRequest("PUT", "/api/v1/savings-goals/"+goalUID+"/add-money/"+txnUID.String(), topUpRequest{Amount: a})
	if err != nil {
		return "", nil, err
	}

	var tuResp *savingsGoalTransferResponse
	resp, err := c.Do(ctx, req, &tuResp)
	if err != nil {
		return tuResp.UID, resp, err
	}
	return tuResp.UID, resp, nil
}

// DeleteSavingsGoal deletes a savings goal for the current customer. It returns http.StatusNoContent
// on success. No payload is returned.
func (c *Client) DeleteSavingsGoal(ctx context.Context, uid string) (*http.Response, error) {
	req, err := c.NewRequest("DELETE", "/api/v1/savings-goals/"+uid, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(ctx, req, nil)
	return resp, err
}
