-- Insert or update global category 'Hutang'
INSERT INTO categories (name, description, type_category, user_id, created_at)
VALUES ('Hutang', 'Kategori khusus untuk transaksi hutang & piutang', 'both', NULL, NOW())
ON CONFLICT (name) DO UPDATE SET user_id = NULL, type_category = 'both';

-- Update all existing transactions linked to loan payments to use the 'Hutang' category
UPDATE transactions
SET category_id = (SELECT id FROM categories WHERE name = 'Hutang' AND user_id IS NULL LIMIT 1)
WHERE id IN (SELECT transaction_id FROM payments WHERE transaction_id IS NOT NULL);
