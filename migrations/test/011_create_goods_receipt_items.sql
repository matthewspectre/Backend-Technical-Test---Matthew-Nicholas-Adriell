CREATE TABLE IF NOT EXISTS public.goods_receipt_items (
    id SERIAL PRIMARY KEY,
    goods_receipt_id INTEGER NOT NULL,
    purchase_order_item_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    received_quantity INTEGER NOT NULL,
    CONSTRAINT goods_receipt_items_quantity_check CHECK (received_quantity > 0),
    CONSTRAINT goods_receipt_items_receipt_fk FOREIGN KEY (goods_receipt_id)
        REFERENCES public.goods_receipts(id) ON DELETE CASCADE,
    CONSTRAINT goods_receipt_items_order_item_fk FOREIGN KEY (purchase_order_item_id)
        REFERENCES public.purchase_order_items(id),
    CONSTRAINT goods_receipt_items_product_fk FOREIGN KEY (product_id)
        REFERENCES public.product(id),
    CONSTRAINT goods_receipt_items_unique_product UNIQUE (goods_receipt_id, product_id)
);
