CREATE TABLE processed_actions (
  idempotency_key CHAR(128) NOT NULL,
  metadata TEXT NULL,
  processed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY ux_idempotency_key (idempotency_key)
);
