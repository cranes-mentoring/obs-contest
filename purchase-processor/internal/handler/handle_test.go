package handler

import (
	"encoding/json"
	"testing"

	"github.com/cranes-mentoring/obs-contest/purchase-processor/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestHandleMessage(t *testing.T) {
	payload := map[string]interface{}{
		"_id": map[string]interface{}{
			"$binary": "r0iMlOukR8y30j1zOQs5rg==",
			"$type":   "00",
		},
		"amount":           100.5,
		"billing_address":  "123 Elm Street, Springfield",
		"card_cvc":         "123",
		"card_expiry":      "12/25",
		"card_holder_name": "John Doessss",
		"card_number":      "4111111111111111",
		"created_at": map[string]interface{}{
			"$date": 1.738061520534e+12,
		},
		"currency":       "USD",
		"payment_method": "credit_card",
		"trace_id":       "249be7d06294bf65a10b27c50c560c3b",
		"transaction_at": map[string]interface{}{
			"$date": 1.737299045e+12,
		},
		"updated_at": map[string]interface{}{
			"$date": 1.738061520534e+12,
		},
		"user_id": map[string]interface{}{
			"$binary": "T448eZoPSgKCO5xqIfRq2A==",
			"$type":   "00",
		},
		"username": "first_user",
	}

	var purchase model.Purchase

	payloadBytes, err := json.Marshal(payload)
	assert.NoError(t, err)

	err = json.Unmarshal(payloadBytes, &purchase)
	assert.NoError(t, err)
}
