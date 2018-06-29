package starling

import (
	"context"
	"net/http"
	"time"
)

// Feed is a slice of Items representing customer transactions
type feed struct {
	Items []Item `json:"feedItems"`
}

// Item is a single customer transaction in their feed
type Item struct {
	FeedItemUID              string    `json:"feedItemUid"`
	CategoryUID              string    `json:"categoryUid"`
	Amount                   Amount    `json:"amount"`
	SourceAmount             Amount    `json:"sourceAmount"`
	Direction                string    `json:"direction"`
	TransactionTime          time.Time `json:"transactionTime"`
	Source                   string    `json:"source"`
	SourceSubType            string    `json:"sourceSubType"`
	Status                   string    `json:"status"`
	CounterPartyType         string    `json:"counterPartyType"`
	CounterPartyUID          string    `json:"counterPartyUid"`
	CounterPartySubEntityUID string    `json:"counterPartySubEntityUid"`
	Reference                string    `json:"reference"`
	Country                  string    `json:"country"`
	SpendingCategory         string    `json:"spendingCategory"`
}

// Feed returns a slice of Items for a given account and category
func (c *Client) Feed(ctx context.Context, act, cat string) ([]Item, *http.Response, error) {
	req, err := c.NewRequest("GET", "/api/v2/feed/account/"+act+"/category/"+cat, nil)
	if err != nil {
		return nil, nil, err
	}

	var f feed
	resp, err := c.Do(ctx, req, &f)
	if err != nil {
		return nil, resp, err
	}
	return f.Items, resp, nil
}
