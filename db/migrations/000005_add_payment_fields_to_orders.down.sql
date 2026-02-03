ALTER TABLE orders
DROP COLUMN IF EXISTS change_amount,
DROP COLUMN IF EXISTS paid_amount,
DROP COLUMN IF EXISTS payment_method;
