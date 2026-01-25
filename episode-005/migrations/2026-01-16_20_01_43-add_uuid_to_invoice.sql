-- episode 004
ALTER TABLE invoice_heads
    ADD COLUMN order_uuid CHAR(36) NOT NULL,
    ADD UNIQUE INDEX ux_invoice_heads_order_uuid (order_uuid);
