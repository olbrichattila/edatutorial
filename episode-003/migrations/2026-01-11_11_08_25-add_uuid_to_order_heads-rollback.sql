ALTER TABLE order_heads
    DROP INDEX ux_order_heads_uuid,
    DROP COLUMN uuid;
    