curl -X POST http://localhost:8080/api/v1/purchases \
-H "Content-Type: application/json" \
-H "Authorization: Bearer <your_oauth2_token>" \
-d '{
  "user_id": "4f8e3c79-9a0f-4a02-823b-9c6a21f46ad8",
  "amount": 100.50,
  "currency": "USD",
  "payment_method": "credit_card",
  "card_holder_name": "John Doe",
  "card_number": "4111111111111111",
  "card_expiry": "12/25",
  "card_cvc": "123",
  "billing_address": "123 Elm Street, Springfield",
  "transaction_date": "2025-01-19T15:04:05Z"
}'