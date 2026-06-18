package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID            int64           `db:"id" json:"id"`
	TransactionID uuid.UUID       `db:"transaction_id" json:"transaction_id"`
	FromState     *string         `db:"from_state" json:"from_state,omitempty"`
	ToState       string          `db:"to_state" json:"to_state"`
	Reason        string          `db:"reason" json:"reason"`
	Metadata      json.RawMessage `db:"metadata" json:"metadata,omitempty"`
	Actor         string          `db:"actor" json:"actor"`
	CreatedAt     time.Time       `db:"created_at" json:"created_at"`
}
