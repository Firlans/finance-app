CREATE TABLE IF NOT EXISTS loans(
    id SERIAL PRIMARY KEY,
    name VARCHAR(50),
    balance DECIMAL(15,4),
    user_id UUID,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT fk_user_loan FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS payments(
    id SERIAL PRIMARY KEY,
    transaction_id INT,
    loan_id INT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT fk_payment_transaction FOREIGN KEY(transaction_id)
        REFERENCES transactions(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_payment_loan FOREIGN KEY(loan_id)
        REFERENCES loans(id)
        ON DELETE CASCADE
)