-- Add related_transaction_id to transactions table
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS related_transaction_id INT NULL REFERENCES transactions(id) ON DELETE SET NULL;

-- Insert global category 'Transfer'
INSERT INTO categories (name, description, type_category, user_id, created_at)
VALUES ('Transfer', 'Kategori khusus untuk pemindahan saldo / tarik tunai', 'both', NULL, NOW())
ON CONFLICT (name) DO UPDATE SET user_id = NULL, type_category = 'both';

-- Insert global category 'fee'
INSERT INTO categories (name, description, type_category, user_id, created_at)
VALUES ('fee', 'Kategori khusus biaya admin / fee', 'credit', NULL, NOW())
ON CONFLICT (name) DO UPDATE SET user_id = NULL, type_category = 'credit';
