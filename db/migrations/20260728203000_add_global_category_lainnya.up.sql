UPDATE categories 
SET user_id = NULL, 
    type_category = 'both' 
WHERE LOWER(name) = 'lainnya' OR LOWER(name) = 'lain-lain';
