DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'purchase_orders_status_check'
          AND conrelid = 'public.purchase_orders'::regclass
    ) THEN
        ALTER TABLE public.purchase_orders
            DROP CONSTRAINT purchase_orders_status_check;
    END IF;

    ALTER TABLE public.purchase_orders
        ADD CONSTRAINT purchase_orders_status_check
        CHECK (status IN ('DRAFT', 'ORDERED', 'PARTIALLY_RECEIVED', 'RECEIVED', 'CANCELLED'));
END
$$;
