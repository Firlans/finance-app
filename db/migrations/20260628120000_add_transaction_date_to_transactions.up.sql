ALTER TABLE transactions
    ADD COLUMN IF NOT EXISTS transaction_date TIMESTAMP WITH TIME ZONE;

-- Backfill for existing rows
UPDATE transactions
SET transaction_date = created_at
WHERE transaction_date IS NULL;

-- Enforce NOT NULL for new/updated data
ALTER TABLE transactions
    ALTER COLUMN transaction_date SET NOT NULL;

