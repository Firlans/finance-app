-- Revert payment types for Debt back to the state of the very first migration
-- (debit -> increase, credit -> decrease)
UPDATE payments
SET type = CASE 
    WHEN (SELECT transaction_type FROM transactions WHERE transactions.id = payments.transaction_id) = 'debit' THEN 'increase'
    ELSE 'decrease'
END
WHERE transaction_id IS NOT NULL
  AND loan_id IN (SELECT id FROM loans WHERE loan_type = 'debt');
