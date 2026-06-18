package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/vishalvignesh12/payflow/internal/models"
)

type AuditRepo struct {
	db *sqlx.DB
}

func NewAuditRepo(db *sqlx.DB) *AuditRepo {
	return &AuditRepo{db: db}
}

func (r *AuditRepo) Write(ctx context.Context, tx sqlx.Tx, entry *models.AuditLog) error {
	query := `
        INSERT INTO audit_log (transaction_id, from_state, to_state, reason, metadata, actor, created_at)
        VALUES (:transaction_id, :from_state, :to_state, :reason, :metadata, :actor, NOW())`

	_, err := r.db.NamedExecContext(ctx, query, entry)

	return err
}

func (r *AuditRepo) GetByTransactionID(ctx context.Context, txid uuid.UUID) ([]*models.AuditLog, error) {
	var entries []*models.AuditLog

	err := r.db.SelectContext(ctx, &entries, "SELECT * FROM audit_log WHERE transaction_id = $1 ORDER BY created_at ASC", txid)

	return entries, err
}
