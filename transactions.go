package starling

import (
	"context"
	"net/http"
)

// Transaction represents the details of a transaction.
type Transaction struct {
	UID       string  `json:"id"`
	Currency  string  `json:"currency"`
	Amount    float64 `json:"amount"`
	Direction string  `json:"direction"`
	Created   string  `json:"created"`
	Narrative string  `json:"narrative"`
	Source    string  `json:"source"`
	Balance   float64 `json:"balance,omitempty"`
}

// Transactions is a list of transaction summaries.
type Transactions struct {
	Transactions []Transaction `json:"transactions"`
}

// HALTransactions is a HAL wrapper around the Transactions type.
type HALTransactions struct {
	Embedded *Transactions `json:"_embedded"`
}

// Transactions returns a list of transaction summaries for the current user. It accepts optional
// time.Time values to request transactions within a given date range. If these values are not provided
// the API returns the last 100 transactions.
func (c *Client) Transactions(ctx context.Context, dr *DateRange) (*[]Transaction, *http.Response, error) {

	req, err := c.NewRequest("GET", "/api/v1/transactions", nil)
	if err != nil {
		return nil, nil, err
	}

	if dr != nil {
		q := req.URL.Query()
		q.Add("from", dr.From.Format("2006-01-02"))
		q.Add("to", dr.To.Format("2006-01-02"))
		req.URL.RawQuery = q.Encode()
	}

	var halResp *HALTransactions
	var txns *Transactions
	resp, err := c.Do(ctx, req, &halResp)
	if err != nil {
		return &txns.Transactions, resp, err
	}

	if halResp.Embedded != nil {
		txns = halResp.Embedded
	}

	return &txns.Transactions, resp, nil
}
