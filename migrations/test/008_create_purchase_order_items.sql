CREATE TABLE IF NOT EXISTS public.purchase_order_items (
    id SERIAL PRIMARY KEY,
    purchase_order_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    ordered_quantity INTEGER NOT NULL,
    received_quantity INTEGER NOT NULL DEFAULT 0,
    CONSTRAINT purchase_order_items_ordered_quantity_check CHECK (ordered_quantity > 0),
    CONSTRAINT purchase_order_items_received_quantity_check CHECK (received_quantity >= 0),
    CONSTRAINT purchase_order_items_order_product_unique UNIQUE (purchase_order_id, product_id),
    CONSTRAINT purchase_order_items_order_fk FOREIGN KEY (purchase_order_id)
        REFERENCES public.purchase_orders(id) ON DELETE CASCADE,
    CONSTRAINT purchase_order_items_product_fk FOREIGN KEY (product_id)
        REFERENCES public.product(id)
);
