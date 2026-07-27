-- Fix payment types for Debt (Hutang) based on the latest corrected API logic.
-- This migration corrects any incorrect data introduced by earlier migrations on the live server.
-- For Debt:
-- Kas Masuk (debit) = Penambahan Hutang (increase)
-- Kas Keluar (credit) = Pembayaran Hutang (decrease)

UPDATE payments
SET type = CASE 
    WHEN (SELECT transaction_type FROM transactions WHERE transactions.id = payments.transaction_id) = 'debit' THEN 'increase'
    WHEN (SELECT transaction_type FROM transactions WHERE transactions.id = payments.transaction_id) = 'credit' THEN 'decrease'
    ELSE type
END
WHERE transaction_id IS NOT NULL
  AND loan_id IN (SELECT id FROM loans WHERE loan_type = 'debt');
