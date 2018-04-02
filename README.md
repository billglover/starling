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

* https://developer.starlingbank.com/

## Features

### Savings Goals

| Method | Resource                                                            | Status      |
|--------|---------------------------------------------------------------------|------------:|
| PUT    | /api/v1/savings-goals/{savingsGoalUid}/add-money/{transferUid}      | Done        |
| GET    | /api/v1/savings-goals/{savingsGoalUid}                              | Done        |
| PUT    | /api/v1/savings-goals/{savingsGoalUid}                              | Done        |
| DELETE | /api/v1/savings-goals/{savingsGoalUid                               |             |
| GET    | /api/v1/savings-goals                                               | Done        |
| GET    | /api/v1/savings-goals/{savingsGoalUid}/photo                        |             |
| PUT    | /api/v1/savings-goals/{savingsGoalUid}/withdraw-money/{transferUid} |             |
| GET    | /api/v1/savings-goals/{savingsGoalUid}/recurring-transfer           |             |
| PUT    | /api/v1/savings-goals/{savingsGoalUid}/recurring-transfer           |             |
| DELETE | /api/v1/savings-goals/{savingsGoalUid}/recurring-transfer           |             |

### Webhooks

| Method | Resource                                                            | Status      |
|--------|---------------------------------------------------------------------|------------:|
| POST   | /your-registered-web-hook-address/card-transaction                  |             |

### Accounts

| Method | Resource                                                            | Status      |
|--------|---------------------------------------------------------------------|------------:|
| GET    | /api/v1/accounts                                                    |             |
| GET    | /api/v1/accounts/balance                                            |             |

### Addresses

| Method | Resource                                                            | Status      |
|--------|---------------------------------------------------------------------|------------:|
| GET    | /api/v1/addresses                                                   |             |

### Contacts

| Method | Resource                                                            | Status      |
|--------|---------------------------------------------------------------------|------------:|
| GET    | /api/v1/contacts                                                    |             |
| POST   | /api/v1/contacts                                                    |             |
| GET    | /api/v1/contacts/{id}                                               |             |
| DELETE | /api/v1/contacts/{id}                                               |             |
| GET    | /api/v1/contacts/{id}/accounts                                      |             |
| GET    | /api/v1/contacts/{contactId}/accounts/{accountId}                   |             |


### Customers

| Method | Resource                                                            | Status      |
|--------|---------------------------------------------------------------------|------------:|
| GET    | /api/v1/customers                                                   |             |

### Direct Debit Mandates

| Method | Resource                                                            | Status      |
|--------|---------------------------------------------------------------------|------------:|
| GET    | /api/v1/direct-debit/mandates                                       |             |
| GET    | /api/v1/direct-debit/mandates/{mandateUid}                          |             |
| DELETE | /api/v1/direct-debit/mandates/{mandateUid}                          |             |

### Who am I?

| Method | Resource                                                            | Status      |
|--------|---------------------------------------------------------------------|------------:|
| GET    | /api/v1/me                                                          |             |

### Cards

| Method | Resource                                                            | Status      |
|--------|---------------------------------------------------------------------|------------:|
| GET    | /api/v1/cards                                                       |             |

### Merchants

| Method | Resource                                                            | Status      |
|--------|---------------------------------------------------------------------|------------:|
| GET    | /api/v1/merchants/{merchantUid}                                     |             |
| GET    | /api/v1/merchants/{merchantUid}/locations/{merchantLocationUid}     |             |

### Payments

| Method | Resource                                                            | Status      |
|--------|---------------------------------------------------------------------|------------:|
| POST   | /api/v1/payments/local                                              |             |
| GET    | /api/v1/payments/scheduled                                          |             |
| POST   | /api/v1/payments/scheduled                                          |             |

### Transactions - Direct Debits

| Method | Resource                                                            | Status      |
|--------|---------------------------------------------------------------------|------------:|
| GET    | /api/v1/transactions/direct-debit                                   |             |
| GET    | /api/v1/transactions/direct-debit/{transactionUid}                  |             |
| PUT    | /api/v1/transactions/direct-debit/{transactionUid}                  |             |

### Transactions - Faster Payments Service

| Method | Resource                                                            | Status      |
|--------|---------------------------------------------------------------------|------------:|
| GET    | /api/v1/transactions/fps/in                                         |             |
| GET    | /api/v1/transactions/fps/in/{transactionUid}                        |             |
| GET    | /api/v1/transactions/fps/out                                        |             |
| GET    | /api/v1/transactions/fps/out/{transactionUid}                       |             |

### Transactions - Card

| Method | Resource                                                              | Status      |
|--------|-----------------------------------------------------------------------|------------:|
| GET    | /api/v1/transactions/mastercard                                       |             |
| GET    | /api/v1/transactions/mastercard/{transactionUid}                      |             |
| PUT    | /api/v1/transactions/mastercard/{transactionUid}                      |             |
| POST   | /api/v1/transactions/mastercard/{transactionUid}/receipt              |             |
| PUT    | /api/v1/transactions/mastercard/{transactionUid}/receipt/{receiptUid} |             |

### Transactions - Any

| Method | Resource                                                            | Status      |
|--------|---------------------------------------------------------------------|------------:|
| GET    | /api/v1/transactions                                                | Done        |
| GET    | /api/v1/transactions/{transactionUid}                               |             |
