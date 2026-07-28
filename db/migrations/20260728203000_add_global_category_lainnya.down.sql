UPDATE categories 
SET type_category = 'expense' 
WHERE LOWER(name) = 'lainnya' OR LOWER(name) = 'lain-lain';
