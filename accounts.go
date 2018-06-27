package starling

import (
	"context"
	"net/http"
)

// Account represents bank account details
type Account struct {
	UID           string `json:"id"`
	Name          string `json:"name"`
	AccountNumber string `json:"accountNumber"`
	SortCode      string `json:"sortCode"`
	Currency      string `json:"currency"`
	IBAN          string `json:"iban"`
	BIC           string `json:"bic"`
	CreatedAt     string `json:"createdAt"`
}

// AccountSummary represents the basic account details
type AccountSummary struct {
	UID             string `json:"accountUid"`
	DefaultCategory string `json:"defaultCategory"`
	Currency        string `json:"currency"`
	CreatedAt       string `json:"createdAt"`
}

// Accounts is a list containing all accounts for a customer
type accounts struct {
	Accounts []AccountSummary `json:"accounts"`
}

// Account returns the the account details for the current customer.
func (c *Client) Account(ctx context.Context) (*Account, *http.Response, error) {
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

// Accounts returns the the accounts held by the current user.
func (c *Client) Accounts(ctx context.Context) ([]AccountSummary, *http.Response, error) {
	req, err := c.NewRequest("GET", "/api/v2/accounts", nil)
	if err != nil {
		return nil, nil, err
	}

	var acts accounts
	resp, err := c.Do(ctx, req, &acts)
	if err != nil {
		return nil, resp, err
	}
	return acts.Accounts, resp, nil
}

// AccountID represents the identifiers for an individual account
type AccountID struct {
	ID     string `json:"accountIdentifier"`
	BankID string `json:"bankIdentifier"`
	IBAN   string `json:"iban"`
	BIC    string `json:"bic"`
}

// AccountID returns the identifiers for an individual account
func (c *Client) AccountID(ctx context.Context, uid string) (*AccountID, *http.Response, error) {
	req, err := c.NewRequest("GET", "/api/v2/accounts/"+uid+"/identifiers", nil)
	if err != nil {
		return nil, nil, err
	}

	var actID *AccountID
	resp, err := c.Do(ctx, req, &actID)
	return actID, resp, err
}

// Balance represents the balance on an account
type Balance struct {
	Cleared     float64 `json:"clearedBalance"`
	Effective   float64 `json:"effectiveBalance"`
	PendingTxns float64 `json:"pendingTransactions"`
	Available   float64 `json:"availableToSpend"`
	Overdraft   float64 `json:"acceptedOverdraft"`
	Currency    string  `json:"currency"`
	Amount      float64 `json:"amount"`
}

// AccountBalance returns the the account balance for the current customer.
func (c *Client) AccountBalance(ctx context.Context) (*Balance, *http.Response, error) {
	req, err := c.NewRequest("GET", "/api/v1/accounts/balance", nil)
	if err != nil {
		return nil, nil, err
	}

	var b *Balance
	resp, err := c.Do(ctx, req, &b)
	if err != nil {
		return b, resp, err
	}
	return b, resp, nil
}
