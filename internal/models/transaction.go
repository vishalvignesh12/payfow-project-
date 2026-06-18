package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusCreated    Status = "Created"
	StatusPending    Status = "Pending"
	StatusProcessing Status = "Processing"
	StatusSuccess    Status = "Success"
	StatusFailed     Status = "Failed"
)

type Transactions struct {
	ID               uuid.UUID       `db:"id" json:"id"`
	IdempotencyKey   string          `db:"idempotency_key" json:"-"`
	Amount           int64           `db:"amount" json:"amount"`
	Currency         string          `db:"currency" json:"currency"`
	Status           Status          `db:"status" json:"status"`
	PayerID          string          `db:"payer_id" json:"payer_id"`
	PayeeID          string          `db:"payee_id" json:"payee_id"`
	Metadata         json.RawMessage `db:"metadata" json:"metadata,omitempty"`
	FailureReason    *string         `db:"failure_reason" json:"failure_reason,omitempty"`
	RetryCount       int             `db:"retry_count" json:"retry_count"`
	NextRetryAt      *time.Time      `db:"next_retry_at" json:"next_retry_at,omitempty"`
	FraudScore       *float64        `db:"fraud_score" json:"fraud_score,omitempty"`
	FraudFlagged     bool            `db:"fraud_flagged" json:"fraud_flagged"`
	PermanentFailure bool            `db:"permanent_failure" json:"permanent_failure"`
	CreatedAt        time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time       `db:"updated_at" json:"updated_at"`
}
