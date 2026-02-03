ALTER TABLE orders
ADD COLUMN payment_method VARCHAR(20) DEFAULT 'midtrans',
ADD COLUMN paid_amount NUMERIC,
ADD COLUMN change_amount NUMERIC;
