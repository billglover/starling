package starling

import (
	"context"
	"net/http"
)

// Receipt is a receipt for a transaction
type Receipt struct {
	UID                string        `json:"receiptUid"`
	EventUID           string        `json:"eventUid"`
	MetadataSource     string        `json:"metadataSource"`
	ReceiptIdentifier  string        `json:"receiptIdentifier"`
	MerchantIdentifier string        `json:"merchantIdentifier"`
	MerchantAddress    string        `json:"merchantAddress"`
	TotalAmount        float64       `json:"totalAmount"`
	TotalTax           float64       `json:"totalTax"`
	TaxReference       string        `json:"taxNumber"`
	AuthCode           string        `json:"authCode"`
	CardLast4          string        `json:"cardLast4"`
	ProviderName       string        `json:"providerName"`
	Items              []ReceiptItem `json:"items"`
	Notes              []ReceiptNote `json:"notes"`
}

// ReceiptItem is a single item on a Receipt
type ReceiptItem struct {
	UID         string  `json:"receiptItemUid"`
	Description string  `json:"description"`
	Quantity    int32   `json:"quantity"`
	Amount      float64 `json:"amount"`
	Tax         float64 `json:"tax"`
	URL         string  `json:"url"`
}

// ReceiptNote is a single item on a Receipt
type ReceiptNote struct {
	UID         string `json:"noteUid"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

// CreateReceipt creates a receipt for a given mastercard transaction.
func (c *Client) CreateReceipt(ctx context.Context, txnUID string, r Receipt) (*http.Response, error) {
	req, err := c.NewRequest("POST", "/api/v1/transactions/mastercard/"+txnUID+"/receipt", r)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(ctx, req, nil)
	return resp, err
}
