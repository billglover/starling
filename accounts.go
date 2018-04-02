package starling

import (
	"context"
	"net/http"
)

// Account represents bank account details
type Account struct {
	UID           string `json:"id"`
	Name          string `json:"name"`
	AccountNUmber string `json:"accountNumber"`
	SortCode      string `json:"sortCode"`
	Currency      string `json:"currency"`
	IBAN          string `json:"iban"`
	BIC           string `json:"bic"`
	CreatedAt     string `json:"createdAt"`
}

// GetAccount returns the the account details for the current customer.
func (c *Client) GetAccount(ctx context.Context) (*Account, *http.Response, error) {
	req, err := c.NewRequest("GET", "/api/v1/accounts", nil)
	if err != nil {
		return nil, nil, err
	}

	var act *Account
	resp, err := c.Do(ctx, req, &act)
	if err != nil {
		return act, resp, err
	}

	return act, resp, nil
}
