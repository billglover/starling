package starling

import "time"

// Feed is a slice of Items representing customer transactions
type Feed struct {
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
