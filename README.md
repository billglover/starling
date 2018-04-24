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
    opts := ClientOptions{BaseURL: baseURL}
    client := NewClientWithOptions(nil, opts)

    txns, _, _ := client.GetTransactions(ctx, dr)

    for i, txn := range txns.Transactions {
        fmt.Println(txn.Created, tx.Amount, txn.Currency, txn.Narrative)
    }
}
```

## Starling Developer Documentation

* [Developer Documentation](https://developer.starlingbank.com/)

## Features

### API Coverage

| Method | Resource                                                              | Status      |
|--------|-----------------------------------------------------------------------|------------:|
| PUT    | /api/v1/savings-goals/{savingsGoalUid}/add-money/{transferUid}        | Done        |
| GET    | /api/v1/savings-goals/{savingsGoalUid}                                | Done        |
| PUT    | /api/v1/savings-goals/{savingsGoalUid}                                | Done        |
| DELETE | /api/v1/savings-goals/{savingsGoalUid}                                | Done        |
| GET    | /api/v1/savings-goals                                                 | Done        |
| GET    | /api/v1/savings-goals/{savingsGoalUid}/photo                          | Done        |
| PUT    | /api/v1/savings-goals/{savingsGoalUid}/withdraw-money/{transferUid}   | Done        |
| GET    | /api/v1/savings-goals/{savingsGoalUid}/recurring-transfer             | Done        |
| PUT    | /api/v1/savings-goals/{savingsGoalUid}/recurring-transfer             | Done        |
| DELETE | /api/v1/savings-goals/{savingsGoalUid}/recurring-transfer             | Done        |
| GET    | /api/v1/accounts                                                      | Done        |
| GET    | /api/v1/accounts/balance                                              | Done        |
| GET    | /api/v1/addresses                                                     | Done        |
| GET    | /api/v1/contacts                                                      | Done        |
| POST   | /api/v1/contacts                                                      | Done        |
| GET    | /api/v1/contacts/{id}                                                 | Done        |
| DELETE | /api/v1/contacts/{id}                                                 | Done        |
| GET    | /api/v1/contacts/{id}/accounts                                        | Done        |
| GET    | /api/v1/contacts/{contactId}/accounts/{accountId}                     | Done        |
| GET    | /api/v1/customers                                                     | Done        |
| GET    | /api/v1/direct-debit/mandates                                         | Done        |
| GET    | /api/v1/direct-debit/mandates/{mandateUid}                            | Done        |
| DELETE | /api/v1/direct-debit/mandates/{mandateUid}                            | Done        |
| GET    | /api/v1/me                                                            | Done        |
| GET    | /api/v1/cards                                                         | Done        |
| GET    | /api/v1/merchants/{merchantUid}                                       | Done        |
| GET    | /api/v1/merchants/{merchantUid}/locations/{merchantLocationUid}       | Done        |
| POST   | /api/v1/payments/local                                                |             |
| GET    | /api/v1/payments/scheduled                                            |             |
| POST   | /api/v1/payments/scheduled                                            |             |
| GET    | /api/v1/transactions/direct-debit                                     | Done        |
| GET    | /api/v1/transactions/direct-debit/{transactionUid}                    | Done        |
| PUT    | /api/v1/transactions/direct-debit/{transactionUid}                    | Done        |
| GET    | /api/v1/transactions/fps/in                                           | Done        |
| GET    | /api/v1/transactions/fps/in/{transactionUid}                          |             |
| GET    | /api/v1/transactions/fps/out                                          |             |
| GET    | /api/v1/transactions/fps/out/{transactionUid}                         |             |
| GET    | /api/v1/transactions/mastercard                                       |             |
| GET    | /api/v1/transactions/mastercard/{transactionUid}                      |             |
| PUT    | /api/v1/transactions/mastercard/{transactionUid}                      |             |
| POST   | /api/v1/transactions/mastercard/{transactionUid}/receipt              |             |
| PUT    | /api/v1/transactions/mastercard/{transactionUid}/receipt/{receiptUid} |             |
| GET    | /api/v1/transactions                                                  | Done        |
| GET    | /api/v1/transactions/{transactionUid}                                 | Done        |

### Webhook Coverage

| Method | Resource                                                            | Status      |
|--------|---------------------------------------------------------------------|------------:|
| POST   | /your-registered-web-hook-address/card-transaction                  |             |
