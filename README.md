# Starling

![Build Status](https://github.com/lildude/starling/actions/workflows/test.yml/badge.svg)

This is an unofficial Go client for the Starling Bank API and a fork of the original package at <https://github.com/billglover/starling>.

Both the Starling Bank API itself and this package are under active development and, whilst we try to keep breaking changes to a minimum, we cannot guarantee a stable interface. We use [Semantic Versioning](https://semver.org) to quantify changes from one release to the next.

> "Major version zero (0.y.z) is for initial development. Anything may change at any time. The public API should not be considered stable."

## Installation

Use Go to fetch the latest version of the package.

```shell
go get -u 'github.com/lildude/starling'
```

## Usage

It is assumed that you are able to provide an OAuth access-token when establishing the Starling client. Depending on your use case, it pay be as simple as passing in the personal access-token provided by Starling when you create an applicaiton. See the section on Personal Access Tokens in the [Starling Developer Docs](https://developer.starlingbank.com/docs) for more information on how to do this.

```go
package main

import (
    "context"
    "fmt"

    "github.com/lildude/starling"
    "golang.org/x/oauth2"
)

func main() {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "{{ACCESS_TOKEN}}"})
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)

	client := NewClient(tc)

	txns, _, _ := client.Transactions(ctx, nil)

	for _, txn := range txns {
		fmt.Println(txn.Created, txn.Amount, txn.Currency, txn.Narrative)
	}
}
```

If you want to use the production API rather than the sandbox, you need to create a client with additional options.

```go
package main

import (
    "context"
    "fmt"
    "net/url"

    "github.com/lildude/starling"
    "golang.org/x/oauth2"
)

func main() {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "{{ACCESS_TOKEN}}"})
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)

	baseURL, _ := url.Parse(ProdURL)
	opts := ClientOptions{BaseURL: baseURL}
	client := NewClientWithOptions(tc, opts)

	txns, _, _ := client.Transactions(ctx, nil)

	for _, txn := range txns {
		fmt.Println(txn.Created, txn.Amount, txn.Currency, txn.Narrative)
	}
}
```

## Starling Bank Developer Documentation

* [Developer Documentation](https://developer.starlingbank.com/)

## Contributors

* [@lildude](https://github.com/lildude/starling/commits?author=lildude)
* [@tuckerwales](https://github.com/lildude/starling/commits?author=tuckerwales)
* [@billglover](https://github.com/lildude/starling/commits?author=billglover)
