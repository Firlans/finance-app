ALTER TABLE categories DROP CONSTRAINT IF EXISTS fk_categories_user_id;
ALTER TABLE categories DROP COLUMN IF EXISTS user_id;