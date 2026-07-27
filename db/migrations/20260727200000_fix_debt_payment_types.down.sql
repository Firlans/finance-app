-- Revert payment types for Debt loans to previous logic
-- previous logic was: debit -> increase, credit -> decrease

UPDATE payments
SET type = CASE 
    WHEN (SELECT transaction_type FROM transactions WHERE transactions.id = payments.transaction_id) = 'credit' THEN 'decrease'
    WHEN (SELECT transaction_type FROM transactions WHERE transactions.id = payments.transaction_id) = 'debit' THEN 'increase'
    ELSE type
END
WHERE transaction_id IS NOT NULL
  AND loan_id IN (SELECT id FROM loans WHERE loan_type = 'debt');
