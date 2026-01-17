CREATE TABLE processed_actions (
  idempotency_key CHAR(85) NOT NULL,
  processed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY ux_idempotency_key (idempotency_key)
);
