-- Fix script for live database to correct ALL payment types without losing new data
-- 1. Ensure initial balance payments (transaction_id IS NULL) are set to 'increase'
UPDATE payments
SET type = 'increase'
WHERE transaction_id IS NULL;

-- 2. Correct payments based on loan type and transaction type
UPDATE payments
SET type = CASE 
    -- For Debt:
    WHEN (SELECT loan_type FROM loans WHERE id = payments.loan_id) = 'debt' THEN
        CASE 
            WHEN (SELECT transaction_type FROM transactions WHERE transactions.id = payments.transaction_id) = 'credit' THEN 'increase'
            WHEN (SELECT transaction_type FROM transactions WHERE transactions.id = payments.transaction_id) = 'debit' THEN 'decrease'
            ELSE type 
        END
    -- For Receivable:
    WHEN (SELECT loan_type FROM loans WHERE id = payments.loan_id) = 'receivable' THEN
        CASE 
            WHEN (SELECT transaction_type FROM transactions WHERE transactions.id = payments.transaction_id) = 'debit' THEN 'increase'
            WHEN (SELECT transaction_type FROM transactions WHERE transactions.id = payments.transaction_id) = 'credit' THEN 'decrease'
            ELSE type 
        END
    ELSE type
END
WHERE transaction_id IS NOT NULL;
