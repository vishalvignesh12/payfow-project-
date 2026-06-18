CREATE TABLE fraud_flags (
    id             BIGSERIAL PRIMARY KEY,
    transaction_id UUID NOT NULL REFERENCES transactions(id),
    rule_name      VARCHAR(100) NOT NULL,
    rule_value     JSONB,
    severity       VARCHAR(20) NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE idempotency_keys (
    key            VARCHAR(255) PRIMARY KEY,
    transaction_id UUID NOT NULL REFERENCES transactions(id),
    response_body  JSONB NOT NULL,
    status_code    INT NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at     TIMESTAMPTZ NOT NULL DEFAULT NOW() + INTERVAL '24 hours'
);