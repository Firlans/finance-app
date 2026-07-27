-- Fix payment types based on loan type and transaction type

-- For Debt (Hutang):
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

-- For Receivable (Piutang):
-- Kas Keluar (credit) = Penambahan Piutang (increase)
-- Kas Masuk (debit) = Pembayaran Piutang (decrease)
UPDATE payments
SET type = CASE 
    WHEN (SELECT transaction_type FROM transactions WHERE transactions.id = payments.transaction_id) = 'credit' THEN 'increase'
    WHEN (SELECT transaction_type FROM transactions WHERE transactions.id = payments.transaction_id) = 'debit' THEN 'decrease'
    ELSE type
END
WHERE transaction_id IS NOT NULL
  AND loan_id IN (SELECT id FROM loans WHERE loan_type = 'receivable');
