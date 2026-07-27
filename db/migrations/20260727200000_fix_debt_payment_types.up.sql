-- Fix payment types for Debt loans based on transaction type
-- For Debt:
-- transaction_type = 'credit' -> type = 'increase'
-- transaction_type = 'debit' -> type = 'decrease'

UPDATE payments
SET type = CASE 
    WHEN (SELECT transaction_type FROM transactions WHERE transactions.id = payments.transaction_id) = 'credit' THEN 'increase'
    WHEN (SELECT transaction_type FROM transactions WHERE transactions.id = payments.transaction_id) = 'debit' THEN 'decrease'
    ELSE type
END
WHERE transaction_id IS NOT NULL
  AND loan_id IN (SELECT id FROM loans WHERE loan_type = 'debt');
