package starling

import (
	"time"
)

// WebHookPayload defines the structure of the Starling web hook payload
type WebHookPayload struct {
	WebhookNotificationUID string         `json:"webhookNotificationUid"`
	Timestamp              time.Time      `json:"timestamp"`
	Content                WebHookContent `json:"content"`
	AccountHolderUID       string         `json:"accountHolderUid"`
	WebhookType            string         `json:"webhookType"`
	CustomerUID            string         `json:"customerUid"`
	UID                    string         `json:"uid"`
}

// WebHookContent defines the structure of the Starling web hook content
type WebHookContent struct {
	Class          string  `json:"class"`
	TransactionUID string  `json:"transactionUid"`
	Amount         float64 `json:"amount"`
	SourceCurrency string  `json:"sourceCurrency"`
	SourceAmount   float64 `json:"sourceAmount"`
	CounterParty   string  `json:"counterParty"`
	Reference      string  `json:"reference"`
	Type           string  `json:"type"`
	ForCustomer    string  `json:"forCustomer"`
}
