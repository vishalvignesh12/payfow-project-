package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/vishalvignesh12/payflow/internal/models"
)

type TransactionRepo struct {
	db *sqlx.DB
}

func NewTransactionRepo(db *sqlx.DB) *TransactionRepo {
	return &TransactionRepo{db: db}
}

func (r *TransactionRepo) Create(ctx context.Context, tx *sqlx.Tx, t *models.Transactions) error {
	query := `
		INSERT INTO transactions
            (id, idempotency_key, amount, currency, status, payer_id, payee_id,
             metadata, fraud_flagged, created_at, updated_at)
        VALUES
            (:id, :idempotency_key, :amount, :currency, :status, :payer_id, :payee_id, :metadata, :fraud_flagged, :created_at, :updated_at)`
	_, err := tx.NamedExecContext(ctx, query, t)
	return err
}

func (r *TransactionRepo) GetById(ctx context.Context, id uuid.UUID) (*models.Transactions, error) {
	var t models.Transactions

	err := r.db.GetContext(ctx, &t, "SELECT * FROM transactions WHERE id = $1", id)

	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TransactionRepo) Update(ctx context.Context, tx *sqlx.Tx, id uuid.UUID, currentstatus, newstatus models.Status, failureReason *string) (int64, error) {
	result, err := tx.ExecContext(ctx, `
        UPDATE transactions
        SET status = $1, failure_reason = $2, updated_at = NOW()
        WHERE id = $3 AND status = $4`,
		newstatus, currentstatus, failureReason, id)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()

}

func (r *TransactionRepo) UpdateRetry(ctx context.Context, retryconut int, id uuid.UUID, nextRetryAt time.Time) error {
	_, err := r.db.ExecContext(ctx, `
        UPDATE transactions
        SET retry_count = $1, next_retry_at = $2, updated_at = NOW()
        WHERE id = $3`,
		id, nextRetryAt, retryconut)

	return err
}

func (r *TransactionRepo) GetPendingQueries(ctx context.Context, maxretires int, limit int) ([]*models.Transactions, error) {
	var txns []*models.Transactions
	err := r.db.SelectContext(ctx, &txns, `
        SELECT * FROM transactions
        WHERE status = 'FAILED'
        AND next_retry_at <= NOW()
        AND retry_count < $1
        AND permanent_failure = FALSE
        ORDER BY next_retry_at ASC
        LIMIT $2`,
		maxretires, limit)
	return txns, err
}

func (r *TransactionRepo) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return r.db.BeginTxx(ctx, nil)
}
