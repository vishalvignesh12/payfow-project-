package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/vishalvignesh12/payflow/internal/models"
)

type FraudRepo struct {
	db *sqlx.DB
}

func NewFraudRepo(db *sqlx.DB) *FraudRepo {
	return &FraudRepo{db: db}
}

func (r *FraudRepo) WriteFlag(ctx context.Context, tx *sqlx.Tx, flag *models.Fraud) error {
	query := `
        INSERT INTO fraud_flags (transaction_id, rule_name, rule_value, severity, created_at)
        VALUES (:transaction_id, :rule_name, :rule_value, :severity, NOW())`

	_, err := tx.NamedExecContext(ctx, query, flag)

	return err
}

func (r *FraudRepo) GetRecent(ctx context.Context, limit int) ([]*models.Fraud, error) {
	var flag []*models.Fraud
	err := r.db.SelectContext(ctx, &flag, "SELECT * FROM fraud_flags ORDER BY created_at DESC LIMIT $1", limit)
	return flag, err
}
