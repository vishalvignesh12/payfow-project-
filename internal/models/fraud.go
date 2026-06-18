package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Fraud struct {
	ID            int64           `db:"id" json:"id"`
	TransactionID uuid.UUID       `db:"transaction_id" json:"transaction_id"`
	RuleName      string          `db:"rule_name" json:"rule_name"`
	RuleValue     json.RawMessage `db:"rule_value" json:"rule_value,omitempty"`
	Severity      string          `db:"severity" json:"severity"`
	CreatedAt     time.Time       `db:"created_at" json:"created_at"`
}
