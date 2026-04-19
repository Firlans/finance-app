CREATE TABLE IF NOT EXISTS password_reset_tokens (
    token UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_password_reset_user FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);
