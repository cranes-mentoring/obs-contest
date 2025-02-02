package entity

import (
	"time"

	"github.com/google/uuid"
)

// Purchase represents a transaction record for a product or service purchase by a user.
type Purchase struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	Username       string    `json:"username"`
	Amount         float64   `json:"amount"`
	Currency       string    `json:"currency"`
	PaymentMethod  string    `json:"payment_method"`
	CardHolderName string    `json:"card_holder_name"`
	CardNumber     string    `json:"-"`
	CardExpiry     string    `json:"-"`
	CardCVC        string    `json:"-"`
	BillingAddress string    `json:"billing_address"`
	TransactionAt  time.Time `json:"transaction_date"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
