-- Add columns to payments
ALTER TABLE payments ADD COLUMN amount DECIMAL(15,4);
ALTER TABLE payments ADD COLUMN payment_date TIMESTAMP;
ALTER TABLE payments ADD COLUMN type VARCHAR(20);

-- Migrate existing data
UPDATE payments 
SET amount = (SELECT amount FROM transactions WHERE transactions.id = payments.transaction_id),
    type = CASE 
        WHEN (SELECT loan_type FROM loans WHERE id = payments.loan_id) = 'debt' THEN
            CASE WHEN (SELECT transaction_type FROM transactions WHERE transactions.id = payments.transaction_id) = 'credit' THEN 'increase' ELSE 'decrease' END
        WHEN (SELECT loan_type FROM loans WHERE id = payments.loan_id) = 'receivable' THEN
            CASE WHEN (SELECT transaction_type FROM transactions WHERE transactions.id = payments.transaction_id) = 'debit' THEN 'increase' ELSE 'decrease' END
    END,
    payment_date = (SELECT transaction_date FROM transactions WHERE transactions.id = payments.transaction_id)
WHERE transaction_id IS NOT NULL;

-- Migrate initial balance from loans as 'increase' payments
INSERT INTO payments (loan_id, amount, payment_date, type, created_at, updated_at)
SELECT id, balance, created_at, 'increase', created_at, updated_at
FROM loans
WHERE balance > 0;

-- Drop balance from loans
ALTER TABLE loans DROP COLUMN balance;
