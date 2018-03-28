package main

import (
	"context"
	"fmt"
	"os"

	"github.com/billglover/starling"
	"github.com/fatih/color"
	"golang.org/x/oauth2"
)

func main() {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "qSXbG5pBtFhwSf9VrypqB0hXE5hhPUGmcw0u89DwhKPYKlT0NTMpZ5mBneXakBhh"},
	)
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)

	client := starling.NewClient(tc)
	txns, _, err := client.GetTransactions(ctx, nil)
	if err != nil {
		fmt.Println("unable to get transactions:", err)
		os.Exit(1)
	}

	for i, txn := range txns.Transactions {
		fmt.Printf("%s: %s %10.2f %s %s\n", color.CyanString("%2d", i), color.GreenString(txn.Created), txn.Amount, txn.Currency, txn.Narrative)
	}

}
