package starling

import (
	"context"
	"net/http"
	"time"
)

// GetTransactions returns a list of transaction summaries for the current user. It accepts optional
// time.Time values to request transactions within a given date range. If these values are not provided
// the API returns the last 100 transactions.
func (c *Client) GetTransactions(ctx context.Context, from, to time.Time) (*Transactions, *http.Response, error) {

	req, err := c.NewRequest("GET", "/api/v1/transactions", nil)
	if err != nil {
		return nil, nil, err
	}

	q := req.URL.Query()
	q.Add("from", from.Format("2006-01-02"))
	q.Add("to", to.Format("2006-01-02"))
	req.URL.RawQuery = q.Encode()

	var txns *Transactions
	resp, err := c.Do(ctx, req, &txns)
	if err != nil {
		return txns, resp, err
	}

	return txns, resp, nil
}
