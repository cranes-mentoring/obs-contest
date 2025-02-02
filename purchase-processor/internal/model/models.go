package model

type DebeziumMessage struct {
	Schema  Schema  `json:"schema"`
	Payload Payload `json:"payload"`
}

type Schema struct {
	Type     string  `json:"type"`
	Fields   []Field `json:"fields"`
	Optional bool    `json:"optional"`
	Name     string  `json:"name"`
}

type Field struct {
	Type        string            `json:"type"`
	Optional    bool              `json:"optional"`
	Name        string            `json:"name"`
	Version     int               `json:"version,omitempty"`
	Parameters  map[string]string `json:"parameters,omitempty"`
	Default     string            `json:"default,omitempty"`
	Field       string            `json:"field,omitempty"`
	Items       *Field            `json:"items,omitempty"`
	FieldsInner []Field           `json:"fields,omitempty"`
}

type Payload struct {
	Before            *string            `json:"before"`
	After             *string            `json:"after"`
	UpdateDescription *UpdateDescription `json:"updateDescription"`
	Source            Source             `json:"source"`
	Op                *string            `json:"op"`
	TsMs              *int64             `json:"ts_ms"`
	Transaction       *Transaction       `json:"transaction"`
}

type UpdateDescription struct {
	RemovedFields   []string             `json:"removedFields"`
	UpdatedFields   *string              `json:"updatedFields"`
	TruncatedArrays []TruncatedArrayItem `json:"truncatedArrays"`
}

type TruncatedArrayItem struct {
	Field string `json:"field"`
	Size  int32  `json:"size"`
}

type Source struct {
	Version    string  `json:"version"`
	Connector  string  `json:"connector"`
	Name       string  `json:"name"`
	TsMs       int64   `json:"ts_ms"`
	Snapshot   string  `json:"snapshot"`
	Db         string  `json:"db"`
	Sequence   *string `json:"sequence"`
	TsUs       *int64  `json:"ts_us"`
	TsNs       *int64  `json:"ts_ns"`
	Collection string  `json:"collection"`
	Ord        int32   `json:"ord"`
	Lsid       *string `json:"lsid"`
	TxnNumber  *int64  `json:"txnNumber"`
	WallTime   *int64  `json:"wallTime"`
}

type Transaction struct {
	ID                  string `json:"id"`
	TotalOrder          int64  `json:"total_order"`
	DataCollectionOrder int64  `json:"data_collection_order"`
}

type Purchase struct {
	Amount         float64 `json:"amount"`
	PaymentMethod  string  `json:"payment_method"`
	CardExpiry     string  `json:"card_expiry"`
	CreatedAt      int64   `json:"created_at"`
	Currency       string  `json:"currency"`
	CardHolderName string  `json:"card_holder_name"`
	CardNumber     string  `json:"card_number"`
	UpdatedAt      int64   `json:"updated_at"`
	UserID         string  `json:"user_id"`
	Username       string  `json:"username"`
	CardCVC        string  `json:"card_cvc"`
	TransactionAt  int64   `json:"transaction_at"`
	TraceID        string  `json:"trace_id"`
	BillingAddress string  `json:"billing_address"`
}
