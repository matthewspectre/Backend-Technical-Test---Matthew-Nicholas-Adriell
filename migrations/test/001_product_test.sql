BEGIN;

CREATE TABLE IF NOT EXISTS product (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    sku VARCHAR,
    name VARCHAR,
    unit VARCHAR,
    is_active SMALLINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at VARCHAR DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE product
    ALTER COLUMN is_active SET DEFAULT 1,
    ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP,
    ALTER COLUMN updated_at DROP NOT NULL,
    ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'product_unique'
          AND conrelid = 'product'::regclass
    ) THEN
        ALTER TABLE product
            ADD CONSTRAINT product_unique UNIQUE (sku);
    END IF;
END;
$$;

INSERT INTO product (sku, name, unit, is_active)
VALUES
    ('TEST-001', 'Product Test One', 'pcs', 1),
    ('TEST-002', 'Product Test Two', 'box', 1)
ON CONFLICT (sku) DO NOTHING;

COMMIT;
