CREATE TABLE audit_log (
    id             BIGSERIAL PRIMARY KEY,
    transaction_id UUID NOT NULL REFERENCES transactions(id),
    from_state     VARCHAR(20),
    to_state       VARCHAR(20) NOT NULL,
    reason         TEXT,
    metadata       JSONB,
    actor          VARCHAR(100),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_txn ON audit_log(transaction_id, created_at);