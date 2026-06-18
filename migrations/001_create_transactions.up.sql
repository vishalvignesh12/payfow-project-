CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE transactions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idempotency_key VARCHAR(255) UNIQUE NOT NULL,
    amount          BIGINT NOT NULL,
    currency        VARCHAR(3) NOT NULL DEFAULT 'INR',
    status          VARCHAR(20) NOT NULL DEFAULT 'CREATED',
    payer_id        VARCHAR(255) NOT NULL,
    payee_id        VARCHAR(255) NOT NULL,
    metadata        JSONB,
    failure_reason  TEXT,
    retry_count     INT NOT NULL DEFAULT 0,
    next_retry_at   TIMESTAMPTZ,
    fraud_score     FLOAT,
    fraud_flagged   BOOLEAN NOT NULL DEFAULT FALSE,
    permanent_failure BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_payer ON transactions(payer_id, created_at DESC);
CREATE INDEX idx_transactions_retry ON transactions(next_retry_at)
    WHERE status = 'FAILED' AND permanent_failure = FALSE;