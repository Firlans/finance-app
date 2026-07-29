UPDATE categories SET type_category = 'loan' WHERE LOWER(name) = 'hutang' AND user_id IS NULL;
UPDATE categories SET type_category = 'transfer' WHERE LOWER(name) = 'transfer' AND user_id IS NULL;
UPDATE categories SET type_category = 'fee' WHERE LOWER(name) = 'fee' AND user_id IS NULL;
