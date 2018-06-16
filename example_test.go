package starling_test

import (
	"context"
	"fmt"
	"net/url"

	"github.com/billglover/starling"
	"golang.org/x/oauth2"
)

// Note: the examples listed here are compiled but not executed while testing.
// See the documentation on [Testing](https://golang.org/pkg/testing/#hdr-Examples)
// for further details.

func Example_sandbox() {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "{{ACCESS_TOKEN}}"})
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)

	client := starling.NewClient(tc)

	txns, _, _ := client.Transactions(ctx, nil)

	for _, txn := range txns {
		fmt.Println(txn.Created, txn.Amount, txn.Currency, txn.Narrative)
	}
}

func Example_live() {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "{{ACCESS_TOKEN}}"})
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)

	baseURL, _ := url.Parse(starling.ProdURL)
	opts := starling.ClientOptions{BaseURL: baseURL}
	client := starling.NewClientWithOptions(tc, opts)

	txns, _, _ := client.Transactions(ctx, nil)

	for _, txn := range txns {
		fmt.Println(txn.Created, txn.Amount, txn.Currency, txn.Narrative)
	}
}
