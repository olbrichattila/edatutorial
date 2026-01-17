-- episode 004
ALTER TABLE invoice_heads
    DROP INDEX ux_invoice_heads_order_uuid,
    DROP COLUMN order_uuid;
