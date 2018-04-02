package starling

import "time"

// ErrorDetail holds the details of an error message
type ErrorDetail struct {
	Message string
}

// SavingsGoalTransferResponse represents the response received after attempting to make an immediate or recurring transfer
// into/out of a savings goal.
type SavingsGoalTransferResponse struct {
	UID     string        `json:"transferUid"` // Unique identifier for the transfer
	Success bool          `json:"success"`     // True if the method completed successfully
	Errors  []ErrorDetail `json:"errors"`      // List of errors if the method request failed
}

// CurrencyAndAmount represents the value and currency of a monetary amount
type CurrencyAndAmount struct {
	Currency   string `json:"currency"`   // ISO-4217 3 character currency code
	MinorUnits int64  `json:"minorUnits"` // Amount in the minor units of the given currency; eg pence in GBP, cents in EUR
}

// TopUpRequest represents request to make an immediate transfer into a savings goal
type TopUpRequest struct {
	Amount CurrencyAndAmount `json:"amount"`
}

// RecurrenceRule defines the pattern for recurring events
type RecurrenceRule struct {
	StartDate string `json:"startDate"`
	Frequency string `json:"frequency"`
	Interval  int32  `json:"interval,omitempty"`
	Count     int32  `json:"count,omitempty"`
	UntilDate string `json:"untilDate,omitempty"`
	WeekStart string `json:"weekStart"`
}

// ScheduledSavingsPaymentRequest represents a request to create scheduled payment into a savings goal
type ScheduledSavingsPaymentRequest struct {
	RecurrenceRule    RecurrenceRule    `json:"recurrenceRule"`
	CurrencyAndAmount CurrencyAndAmount `json:"currencyAndAmount"`
}

// SavingsGoalPhoto is a photo associated to a savings goal
type SavingsGoalPhoto struct {
	Base64EncodedPhoto string `json:"base64EncodedPhoto"` // A text (base 64) encoded picture to associate with the savings goal
}

// WithdrawalRequest is a request to withdraw money from a savings goal
type WithdrawalRequest struct {
	Amount CurrencyAndAmount `json:"amount"`
}

// MastercardTransactionPayload is the webhook payload for mastercard transactions
type MastercardTransactionPayload struct {
	WebhookNotificationUID string            `json:"webhookNotificationUid"` // Unique identifier of the webhook dispatch event
	CustomerUID            string            `json:"customerUid"`            // Unique identifier of the customer
	WebhookType            string            `json:"webhookType"`            // The type of the event
	EventUID               string            `json:"eventUid"`               // Unique identifier of the customer transaction event
	TransactionAmount      CurrencyAndAmount `json:"transactionAmount"`
	SourceAmount           CurrencyAndAmount `json:"sourceAmount"`
	Direction              string            `json:"direction"`            // The cashflow direction of the card transaction
	Description            string            `json:"description"`          // The transaction description, usually the name of the merchant
	MerchantUID            string            `json:"merchantUid"`          // The unique identifier of the merchant
	MerchantLocationUID    string            `json:"merchantLocationUid"`  // The unique identifier of the merchant location
	Status                 string            `json:"status"`               // The status of the transaction
	TransactionMethod      string            `json:"transactionMethod"`    // The method of card usage
	TransactionTimestamp   string            `json:"transactionTimestamp"` // Timestamp of the card transaction
	MerchantPosData        MerchantPosData   `json:"merchantPosData"`
}

// MerchantPosData is data relating to the merchant at the point-of-sale
type MerchantPosData struct {
	PosTimestamp       string `json:"posTimestamp"`       // The transaction time as reported at the point of sale
	CardLast4          string `json:"cardLast4"`          // The last 4 digits of the mastercard PAN
	AuthorisationCode  string `json:"authorisationCode"`  // The authorisation code of the transaction, as reported at the point of sale
	Country            string `json:"country"`            // The country of the transaction, in ISO-3 format
	MerchantIdentifier string `json:"merchantIdentifier"` // The merchant identifier as reported by Mastercard AKA mid
}

// Address is the physical address of the customer
type Address struct {
	StreetAddress string `json:"streetAddress"`
	City          string `json:"city"`
	Country       string `json:"country"`
	Postcode      string `json:"postcode"`
}

// Addresses are the current and previous physical addresses
type Addresses struct {
	Current  Address   `json:"current"`
	Previous []Address `json:"previous"`
}

// OptionalContact identifies the presence of a contact
type OptionalContact optional

// Contact represents the details of a payee
type Contact struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Contacts are a list of payees
type Contacts struct {
	Contacts []Contact
}

// ContactAccount holds payee account details
type ContactAccount struct {
	UID           string `json:"id"`
	Type          string `json:"type"`
	Name          string `json:"name"`
	AccountNumber string `json:"accountNumber"`
	SortCode      string `json:"sortCode"`
}

// OptionalContactAccounts identifies the presence of a contact accounts
type OptionalContactAccounts optional

// ContactAccounts holds a list of accounts for a payee
type ContactAccounts struct {
	ContactAccounts []ContactAccount `json:"contactAccounts"`
}

// OptionalContactAccount identifies the presence of a contact account
type OptionalContactAccount optional

// DirectDebitMandate represents a single mandate
type DirectDebitMandate struct {
	UID            string `json:"uid"`
	Reference      string `json:"reference"`
	Status         string `json:"status"`
	Source         string `json:"source"`
	Created        string `json:"created"`
	Cancelled      string `json:"cancelled"`
	OriginatorName string `json:"originatorName"`
	OriginatorUID  string `json:"originatorUid"`
}

// DirectDebitMandates represents a list of mandates
type DirectDebitMandates struct {
	Mandates []DirectDebitMandate `json:"mandates"`
}

// OptionalCard identifies the presence of a card
type OptionalCard optional

// Card represents card details
type Card struct {
	UID                 string  `json:"id"`
	NameOnCard          string  `json:"nameOnCard"`
	Type                string  `json:"type"`
	Enabled             bool    `json:"enabled"`
	Cancelled           bool    `json:"cancelled"`
	ActivationRequested bool    `json:"activationRequested"`
	Activated           bool    `json:"activated"`
	DispatchDate        string  `json:"dispatchDate"`
	LastFourDigits      string  `json:"lastFourDigits"`
	Transactions        HALLink `json:"transactions"`
}

// HALLink is a link to another resource
type HALLink struct {
	HREF        string `json:"href"`
	Templated   bool   `json:"templated"`
	Type        string `json:"type"`
	Deprecation string `json:"deprecation"`
	Name        string `json:"name"`
	Profile     string `json:"profile"`
	Title       string `json:"title"`
	HREFLang    string `json:"hreflang"`
}

// Merchant represents details of a merchant
type Merchant struct {
	UID             string `json:"merchantUid"`
	Name            string `json:"name"`
	Website         string `json:"website"`
	PhoneNumber     string `json:"phoneNumber"`
	TwitterUsername string `json:"twitterUsername"`
}

// MerchantLocation represents details of a merchant location
type MerchantLocation struct {
	UID                            string  `json:"merchantLocationUid"`
	MerchantUID                    string  `json:"merchantUid"`
	Merchant                       HALLink `json:"merchant"`
	MerchantName                   string  `json:"merchantName"`
	LocationName                   string  `json:"locationName"`
	Address                        string  `json:"address"`
	PhoneNumber                    string  `json:"phoneNUmber"`
	GooglePlaceID                  string  `json:"googlePlaceId"`
	MastercardMerchantCategoryCode int32   `json:"mastercardMerchantCategoryCode"`
}

// PaymentAmount represents the currency and amount of a payment
type PaymentAmount struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

// ScheduledPayment represents a scheduled payment
type ScheduledPayment struct {
	LocalPayment
	RecurrenceRule RecurrenceRule `json:"recurrenceRule"`
}

// LocalPayment represents a local payment
type LocalPayment struct {
	Payment               PaymentAmount `json:"payment"`
	DestinationAccountUID string        `json:"destinationAccountUid"`
	Reference             string        `json:"reference"`
}

// PaymentOrder is a single PaymentOrder
type PaymentOrder struct {
	UID                        string         `json:"paymentOrderId"`
	Currency                   string         `json:"currency"`
	Amount                     float64        `json:"amount"`
	Reference                  string         `json:"reference"`
	ReceivingContactAccountUID string         `json:"receivingContactAccountId"`
	RecipientName              string         `json:"recipientName"`
	Immediate                  bool           `json:"immediate"`
	RecurrenceRule             RecurrenceRule `json:"recurrenceRule"`
	StartDate                  string         `json:"startDate"`
	NextDate                   string         `json:"nextDate"`
	CancelledAt                string         `json:"cancelledAt"`
	PaymentType                string         `json:"paymentType"`
	MandateUID                 string         `json:"mandateId"`
}

// PaymentOrders is a list of PaymentOrders
type PaymentOrders struct {
	NextPage      HALLink        `json:"nextPage"`
	PaymentOrders []PaymentOrder `json:"paymentOrders"`
}

// SpendingCategory is the category associated with a transaction
type SpendingCategory struct {
	SpendingCategory string `json:"spendingCategory"`
}

// DirectDebitTransaction represents details of a direct debit transaction
type DirectDebitTransaction struct {
	transaction
	MandateUID          string `json:"mandateId"`
	Type                string `json:"type"`
	MerchantUID         string `json:"merchantId"`
	MerchantLocationUID string `json:"merchantLocationId"`
	SpendingCategory    string `json:"spendingCategory"`
	Country             string `json:"country"`
}

// DirectDebitTransactions is a list of direct debit transactions
type DirectDebitTransactions struct {
	NextPage     HALLink                  `json:"nextPage"`
	Transactions []DirectDebitTransaction `json:"transactions"`
}

// FPSInTransaction represents details of an inbound faster payments transaction
type FPSInTransaction struct {
	transaction
	SendingContactUID        string  `json:"sendingContactId"`
	SendingContactAccountUID string  `json:"sendingContactAccountId"`
	SendingContactAccount    HALLink `json:"sendingContactAccount"`
}

// FPSInTransactions is a list of inbound faster payment transactions
type FPSInTransactions struct {
	NextPage     HALLink            `json:"nextPage"`
	Transactions []FPSInTransaction `json:"transactions"`
}

// FPSOutTransaction represents details of an outbound faster payments transaction
type FPSOutTransaction struct {
	transaction
	ReceivingContactUID        string  `json:"receivingContactId"`
	ReceivingContactAccountUID string  `json:"receivingContactAccountId"`
	ReceivingContactAccount    HALLink `json:"receivingContactAccount"`
}

// FPSOutTransactions is a list of outbound faster payment transactions
type FPSOutTransactions struct {
	NextPage     HALLink             `json:"nextPage"`
	Transactions []FPSOutTransaction `json:"transactions"`
}

// MastercardTransaction represents the details of a card transaction
type MastercardTransaction struct {
	transaction
	Method            string  `json:"mastercardTransactionMethod"`
	Status            string  `json:"status"`
	SourceAmount      float64 `json:"sourceAmount"`
	SourceCurrency    string  `json:"sourceCurrency"`
	MerchantUID       string  `json:"merchantId"`
	SpendingCategory  string  `json:"spendingCategory"`
	Country           string  `json:"country"`
	POSTimestamp      int64   `json:"posTimestamp"`
	AuthorisationCode string  `json:"authorisationCode"`
	EventUID          string  `json:"eventUid"`
	Receipt           Receipt `json:"receipt"`
	CardLast4         string  `json:"cardLast4"`
}

// MastercardTransactions is a list of Mastercard transactions
type MastercardTransactions struct {
	NextPage     HALLink                 `json:"nextPage"`
	Transactions []MastercardTransaction `json:"transactions"`
}

// Receipt is a receipt for a transaction
type Receipt struct {
	UID                string        `json:"receiptUid"`
	EventUID           string        `json:"eventUid"`
	MetadataSource     string        `json:"metadataSource"`
	ReceiptIdentifier  string        `json:"receiptIdentifier"`
	MerchantIdentifier string        `json:"merchantIdentifier"`
	TotalAmount        float64       `json:"totalAmount"`
	TotalTax           float64       `json:"totalTax"`
	AuthCode           string        `json:"authCode"`
	CardLast4          string        `json:"cardLast4"`
	Items              []ReceiptItem `json:"items"`
}

// ReceiptItem is a single item on a Receipt
type ReceiptItem struct {
	UID         string  `json:"receiptItemUid"`
	Description string  `json:"description"`
	Quantity    int32   `json:"quantity"`
	Amount      float64 `json:"amount"`
	Tax         float64 `json:"tax"`
}

// ReceiptUID is an un-used type
type ReceiptUID struct{}

// OptionalTransactionSummary indicates the presence of a TransactionSummary
type OptionalTransactionSummary optional

type optional struct {
	Present bool `json:"present"`
}

// DateRange holds two dates that represent a range. It is typically
// used when providing a range when querying the API.
type DateRange struct {
	From time.Time
	To   time.Time
}
