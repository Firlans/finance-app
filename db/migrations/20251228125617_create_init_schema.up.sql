CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY, -- Menggunakan tipe data UUID native
    username VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT users_email_unique UNIQUE (email),
    CONSTRAINT users_username_unique UNIQUE (username)
);
CREATE TABLE IF NOT EXISTS budgets (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    budget NUMERIC(15, 2) NOT NULL,
    month TIMESTAMP WITH TIME ZONE NOT NULL,
    year TIMESTAMP WITH TIME ZONE NOT NULL,
    start_date DATE,
    iteration VARCHAR(20),
    CONSTRAINT fk_user
    FOREIGN KEY(user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS accounts(
    id UUID PRIMARY KEY,
    account_name VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    account_id UUID,
    amount NUMERIC(15, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user
    FOREIGN KEY(user_id)
    REFERENCES users(id)
    ON DELETE CASCADE,
    CONSTRAINT fk_account
    FOREIGN KEY(account_id)
    REFERENCES accounts(id)
    ON DELETE CASCADE
);