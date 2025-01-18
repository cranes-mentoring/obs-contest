package model

import (
	"time"

	"github.com/google/uuid"
)

// PurchaseRequest represents the data required to initiate a purchase transaction by a user.
type PurchaseRequest struct {
	UserID         uuid.UUID `json:"user_id"`
	Amount         float64   `json:"amount"`
	Currency       string    `json:"currency"`
	PaymentMethod  string    `json:"payment_method"`
	CardHolderName string    `json:"card_holder_name"`
	CardNumber     string    `json:"card_number"`
	CardExpiry     string    `json:"card_expiry"`
	CardCVC        string    `json:"card_cvc"`
	BillingAddress string    `json:"billing_address"`
	TransactionAt  time.Time `json:"transaction_date"`
}

// PurchaseResponse represents the response object for a purchase operation containing transaction details and status.
type PurchaseResponse struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}
