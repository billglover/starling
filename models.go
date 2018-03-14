package starling

/*
- UID used as the identifier for each resource, renamed for consistency
- TODO: could we extract common fields for similar types e.g. Transaction?
*/

type ErrorDetail struct {
	Message string // The error message
}

// SavingsGoalTransferResponse represents the response received after attempting to make an immediate or recurring transfer
// into/out of a savings goal.
type SavingsGoalTransferResponse struct {
	UID     string        `json:"transferUid"` // Unique identifier for the transfer
	Success bool          `json:"success"`     // True if the method completed successfully
	Errors  []ErrorDetail `json:"errors"`      // List of errors if the method request failed
}

type CurrencyAndAmount struct {
	Currency   string `json:"currency"`   // ISO-4217 3 character currency code
	MinorUnits int64  `json:"minorUnits"` // Amount in the minor units of the given currency; eg pence in GBP, cents in EUR
}

// TopUpRequest represents request to make an immediate transfer into a savings goal
type TopUpRequest struct {
	Amount CurrencyAndAmount `json:"amount"`
}

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

// CreateOrUpdateSavingsGoalResponse represents a response after attempting to create a savings goal
type CreateOrUpdateSavingsGoalResponse struct {
	SavingsGoalUid string        `json:"savingsGoalUid"`
	success        bool          `json:"success"`
	errors         []ErrorDetail `json:"errors"`
}

// SavingsGoalRequest is a request to create a new savings goal
type SavingsGoalRequest struct {
	Name               string            `json:"name"`     // Name of the savings goal
	Currency           string            `json:"currency"` // ISO-4217 3 character currency code of the savings goal
	Target             CurrencyAndAmount `json:"target"`
	Base64EncodedPhoto string            `json:"base64EncodedPhoto"` // A text (base 64) encoded picture to associate with the savings goal
}

// SavingsGoal is a goal defined by a customer to hold savings
type SavingsGoal struct {
	UID             string            `json:"uid"`  // Unique identifier of the savings goal
	Name            string            `json:"name"` // Name of the savings goal
	Target          CurrencyAndAmount `json:"target"`
	TotalSaved      CurrencyAndAmount `json:"totalSaved"`
	SavedPercentage int32             `json:"savedPercentage"` // Percentage of target currently deposited in the savings goal
}

// SavingsGoals is a list containing all savings goals for customer
type SavingsGoals struct {
	SavingsGoalList []SavingsGoal `json:"savingsGoalList"`
}

// SavingsGoalPhoto is a photo associated to a savings goal
type SavingsGoalPhoto struct {
	base64EncodedPhoto string `json:"base64EncodedPhoto"` // A text (base 64) encoded picture to associate with the savings goal
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

// Account represents bank account details
type Account struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	AccountNUmber string `json:"accountNumber"`
	SortCode      string `json:"sortCode"`
	Currency      string `json:"currency"`
	IBAN          string `json:"iban"`
	BIC           string `json:"bic"`
	CreatedAt     string `json:"createdAt"`
}

type Balance struct {
	ClearedBalance      float64 `json:"clearedBalance"`
	EffectiveBalance    float64 `json:"effectiveBalance"`
	PendingTransactions float64 `json:"pendingTransactions"`
	AvailableToSpend    float64 `json:"availableToSpend"`
	AcceptedOverdraft   float64 `json:"acceptedOverdraft"`
	Currency            string  `json:"currency"`
	Amount              float64 `json:"amount"`
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

type Optional optional

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
	ID            string `json:"id"` // Unique account identifier of contact to be added
	Type          string `json:"type"`
	Name          string `json:"name"`          // Contact name
	AccountNumber string `json:"accountNumber"` // Contact account number
	SortCode      string `json:"sortCode"`      // Contact sort code
}

type OptionalContactAccounts optional

// ContactAccounts holds a list of accounts for a payee
type ContactAccounts struct {
	ContactAccounts []ContactAccount `json:"contactAccounts"`
}

type OptionalContactAccount optional

type OptionalCustomer optional

// Customer represents the personal details of a customer
type Customer struct {
	CustomerUID       string `json:"customerUid"`
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	DateOfBirth       string `json:"dateOfBirth"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	AccountHolderType string `json:"accountHolderType"`
}

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

type Identity struct {
	CustomerUID      string   `json:"customerUid"`
	ExpiresAt        string   `json:"expiresAt"`
	Authenticated    bool     `json:"authenticated"`
	ExpiresInSeconds int64    `json:"expiresInSeconds"`
	Scopes           []string `json:"scopes"`
}

type OptionalCard optional

// Card holds card details
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

type Merchant struct{}
type MerchantLocation struct{}
type PaymentAmount struct{}
type ScheduledPayment struct{}
type LocalPayment struct{}
type PaymentOrder struct{}
type PaymentOrders struct{}
type SpendingCategory struct{}
type DirectDebitTransaction struct{}
type DirectDebitTransactions struct{}
type FPSInTransaction struct{}
type FPSInTransactions struct{}
type FPSOutTransaction struct{}

// FPSOutTransactions is a list of FPSOutTransaction items
type FPSOutTransactions struct {
	Transactions []FPSOutTransaction `json:"transactions"`
}

// MastercardTransaction represents the details of a card transaction
type MastercardTransaction struct {
	UID               string  `json:"id"`
	Currency          string  `json:"currency"`
	Amount            float64 `json:"amount"`
	Direction         string  `json:"direction"`
	Created           string  `json:"created"`
	Narrative         string  `json:"narrative"`
	Source            string  `json:"source"`
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

// TransactionSummary represents an individual transaction
type TransactionSummary struct {
	UID       string  `json:"id"`
	Currency  string  `json:"currency"`
	Amount    float64 `json:"amount"`
	Direction string  `json:"direction"`
	Created   string  `json:"created"`
	Narrative string  `json:"narrative"`
	Source    string  `json:"source"`
	Balance   float64 `json:"balance"`
}

// Transactions is a list of transaction summaries
type Transactions struct {
	Transactions []TransactionSummary `json:"transactions"`
}

// OptionalTransactionSummary indicates the presence of a TransactionSummary
type OptionalTransactionSummary optional

type optional struct {
	Present bool `json:"present"`
}
