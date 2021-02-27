package starling_test

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/lildude/starling"
	"golang.org/x/oauth2"
)

// Note: the examples listed here are compiled but not executed while testing.
// See the documentation on [Testing](https://golang.org/pkg/testing/#hdr-Examples)
// for further details.

func Example_account() {
	godotenv.Load(".env")
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("STARLING_DEV_TOKEN")})
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)

	baseURL, _ := url.Parse(starling.ProdURL)
	opts := starling.ClientOptions{BaseURL: baseURL}
	client := starling.NewClientWithOptions(tc, opts)
	acct, _, err := client.Accounts(ctx)
	if err != nil {
		log.Fatalf("Whoops: %v", err)
	}

	// Last month
	since := time.Now().AddDate(0, -1, 0)
	txns, _, err := client.Feed(ctx, acct[0].UID, acct[0].DefaultCategory, since)
	if err != nil {
		log.Fatalf("Whoops: %v", err)
	}

	for _, txn := range txns {
		fmt.Println(txn.TransactionTime, txn.Amount, txn.Amount.Currency, txn.Direction)
	}
}

func Example_balance() {
	godotenv.Load(".env")
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("STARLING_DEV_TOKEN")})
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)

	baseURL, _ := url.Parse(starling.ProdURL)
	opts := starling.ClientOptions{BaseURL: baseURL}
	client := starling.NewClientWithOptions(tc, opts)
	acct, _, err := client.Accounts(ctx)
	if err != nil {
		log.Fatalf("Whoops: %v", err)
	}

	bal, _, err := client.AccountBalance(ctx, acct[0].UID)
	if err != nil {
		log.Fatalf("Whoops: %v", err)
	}
	fmt.Printf("%v", bal)
}