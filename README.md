# Starling

This is a Go client for the Starling Bank API.

Both the API itself and this client are under active development and cannot guarantee a stable interface.

## Installation

```shell
go get 'github.com/billglover/starling'
```

## Usage

```go
package main

import (
    "context"
    "fmt"

    "github.com/billglover/starling"
    "golang.org/x/oauth2"
)

func main() {
    ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "{{ACCESS_TOKEN}}"})
    ctx := context.Background()
    tc := oauth2.NewClient(ctx, ts)

    client := starling.NewClient(tc)

    txns, _, _ := client.GetTransactions(ctx, dr)

    for i, txn := range txns.Transactions {
			fmt.Println(txn.Created, tx.Amount, txn.Currency, txn.Narrative)
	}
}
```

If you want to use the production API rather than the sandbox, you need to create a client with additional options.


```go
package main

import (
    "context"
    "fmt"

    "github.com/billglover/starling"
    "golang.org/x/oauth2"
)

func main() {
    ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "{{ACCESS_TOKEN}}"})
    ctx := context.Background()
    tc := oauth2.NewClient(ctx, ts)

    baseURL, _ := url.Parse("https://dummyurl:4000")
	opts := ClientOptions{BaseURL: baseURL,
	}
    client := NewClientWithOptions(nil, opts)

    txns, _, _ := client.GetTransactions(ctx, dr)

    for i, txn := range txns.Transactions {
			fmt.Println(txn.Created, tx.Amount, txn.Currency, txn.Narrative)
    }
}
```

## Starling Developer Documentation

* https://developer.starlingbank.com/
