-- Add balance back to loans
ALTER TABLE loans ADD COLUMN balance DECIMAL(15,4);

-- Migrate 'increase' payments back to loans balance
UPDATE loans
SET balance = (
    SELECT COALESCE(SUM(amount), 0)
    FROM payments
    WHERE payments.loan_id = loans.id AND payments.type = 'increase'
);

-- Delete all 'increase' payments as they are now represented by loan balance
DELETE FROM payments WHERE type = 'increase';

-- Drop columns from payments
ALTER TABLE payments DROP COLUMN amount;
ALTER TABLE payments DROP COLUMN payment_date;
ALTER TABLE payments DROP COLUMN type;
