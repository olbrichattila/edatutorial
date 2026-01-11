ALTER TABLE order_heads
    ADD COLUMN uuid CHAR(36) NOT NULL,
    ADD UNIQUE INDEX ux_order_heads_uuid (uuid);
