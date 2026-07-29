ALTER TABLE transactions DROP COLUMN IF EXISTS related_transaction_id;
DELETE FROM categories WHERE name IN ('Transfer', 'fee') AND user_id IS NULL;
