package starling

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
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

// Validate takes an http request and a web-hook secret and validates the
// request signature matches the signature provided in the X-Hook-Signature
// header. An error is returned if unable to parse the body of the request.
func Validate(r *http.Request, secret string) (bool, error) {
	if r.Body == nil {
		return false, fmt.Errorf("no body to validate")
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return false, err
	}

	body := ioutil.NopCloser(bytes.NewBuffer(buf))
	r.Body = body

	sha512 := sha512.New()
	sha512.Write([]byte(secret + string(buf)))
	recSig := base64.StdEncoding.EncodeToString(sha512.Sum(nil))
	reqSig := r.Header.Get("X-Hook-Signature")
	fmt.Println("reqSig:", reqSig)
	fmt.Println("recSig:", recSig)
	if reqSig != recSig {
		return false, nil
	}
	return true, nil
}
