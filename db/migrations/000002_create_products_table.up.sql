BEGIN;

-- Tabel products untuk menyimpan informasi produk
CREATE TABLE IF NOT EXISTS products (
    product_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE, 
    description TEXT,
    photo_url VARCHAR(255),
    price NUMERIC(10, 2) NOT NULL CHECK (price >= 0), 
    stock INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0), 
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ 
);

COMMIT;