# Starling

[![Build Status](https://travis-ci.com/billglover/starling.svg?branch=master)](https://travis-ci.com/billglover/starling)

This is an unofficial Go client for the Starling Bank API.

Both the Starling Bank API itself and this package are under active development and, whilst we try to keep breaking changes to a minimum, we cannot guarantee a stable interface. We use [Semantic Versioning](https://semver.org) to quantify changes from one release to the next.

> "Major version zero (0.y.z) is for initial development. Anything may change at any time. The public API should not be considered stable."

## Installation

Use Go to fetch the latest version of the package.

```shell
go get -u 'github.com/billglover/starling'
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
    opts := ClientOptions{BaseURL: baseURL}
    client := NewClientWithOptions(nil, opts)

    txns, _, _ := client.GetTransactions(ctx, dr)

    for i, txn := range txns.Transactions {
        fmt.Println(txn.Created, tx.Amount, txn.Currency, txn.Narrative)
    }
}
```

## Starling Bank Developer Documentation

* [Developer Documentation](https://developer.starlingbank.com/)
