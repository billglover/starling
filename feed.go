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

// FeedOpts defines options that can be passed when requesting a feed
type FeedOpts struct {
	Since time.Time
}

// Feed returns a slice of Items for a given account and category. It returns an error if unable
// to retrieve the feed.
// Note: Feed uses the v2 API which is still under active development.
func (c *Client) Feed(ctx context.Context, act, cat string, opts *FeedOpts) ([]Item, *http.Response, error) {
	req, err := c.NewRequest("GET", "/api/v2/feed/account/"+act+"/category/"+cat, nil)
	if err != nil {
		return nil, nil, err
	}

	if opts != nil {
		q := req.URL.Query()
		q.Add("changesSince", opts.Since.Format(time.RFC3339Nano))
		req.URL.RawQuery = q.Encode()
	}

	var f feed
	resp, err := c.Do(ctx, req, &f)
	if err != nil {
		return nil, resp, err
	}
	return f.Items, resp, nil
}

// FeedItem returns a feed Item for a given account and category. It returns an error if unable to
// retrieve the feed Item.
// Note: FeedItem uses the v2 API which is still under active development.
func (c *Client) FeedItem(ctx context.Context, act, cat, itm string) (*Item, *http.Response, error) {
	req, err := c.NewRequest("GET", "/api/v2/feed/account/"+act+"/category/"+cat+"/"+itm, nil)
	if err != nil {
		return nil, nil, err
	}

	var i Item
	resp, err := c.Do(ctx, req, &i)
	if err != nil {
		return nil, resp, err
	}
	return &i, resp, nil
}
