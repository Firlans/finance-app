UPDATE categories SET type_category = 'both' WHERE LOWER(name) = 'hutang' AND user_id IS NULL;
UPDATE categories SET type_category = 'both' WHERE LOWER(name) = 'transfer' AND user_id IS NULL;
UPDATE categories SET type_category = 'credit' WHERE LOWER(name) = 'fee' AND user_id IS NULL;
