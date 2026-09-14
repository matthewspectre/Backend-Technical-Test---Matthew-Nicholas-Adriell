CREATE TABLE IF NOT EXISTS public.inventory_movements (
    id SERIAL PRIMARY KEY,
    warehouse_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    movement_type VARCHAR(30) NOT NULL,
    quantity INTEGER NOT NULL,
    reference VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT inventory_movements_quantity_check CHECK (quantity > 0),
    CONSTRAINT inventory_movements_type_check CHECK (movement_type IN ('PURCHASE_RECEIPT')),
    CONSTRAINT inventory_movements_warehouse_fk FOREIGN KEY (warehouse_id)
        REFERENCES public.warehouse(id),
    CONSTRAINT inventory_movements_product_fk FOREIGN KEY (product_id)
        REFERENCES public.product(id)
);
