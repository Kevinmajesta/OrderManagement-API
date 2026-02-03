BEGIN;

CREATE TABLE IF NOT EXISTS receipts (
    receipt_id UUID PRIMARY KEY,
    order_id UUID NOT NULL,
    user_id UUID NOT NULL,
    subtotal NUMERIC(10,2) NOT NULL,
    tax_amount NUMERIC(10,2) NOT NULL,
    total_amount NUMERIC(10,2) NOT NULL,
    payment_method VARCHAR(20),
    payment_status VARCHAR(20),
    receipt_number VARCHAR(20) UNIQUE NOT NULL,
    cashier_name VARCHAR(255),
    store_name VARCHAR(255),
    store_address TEXT,
    store_phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT receipts_order_fk FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE,
    CONSTRAINT receipts_user_fk FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS receipt_items (
    receipt_item_id UUID PRIMARY KEY,
    receipt_id UUID NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    unit_price NUMERIC(10,2) NOT NULL,
    total_price NUMERIC(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT receipt_items_receipt_fk FOREIGN KEY (receipt_id) REFERENCES receipts(receipt_id) ON DELETE CASCADE
);

COMMIT;
