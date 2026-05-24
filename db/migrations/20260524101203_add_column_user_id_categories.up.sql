ALTER TABLE categories ADD COLUMN user_id UUID NULL;
ALTER TABLE categories ADD CONSTRAINT fk_categories_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;