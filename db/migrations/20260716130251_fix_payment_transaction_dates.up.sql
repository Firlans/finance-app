-- Fix transaction_date for payments that defaulted to zero time ('0001-01-01')
UPDATE transactions
SET transaction_date = created_at
WHERE transaction_date < '2000-01-01'
  AND id IN (SELECT transaction_id FROM payments WHERE transaction_id IS NOT NULL);
