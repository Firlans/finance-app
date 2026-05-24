ALTER TABLE transactions DROP COLUMN IF EXISTS category_id;
ALTER TABLE transactions DROP CONSTRAINT IF EXISTS transactions_category_id_fkey;