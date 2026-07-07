package transactions

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Save(ctx context.Context, transaction *Transaction) error
	GetTransactions(ctx context.Context, userID string) ([]Transaction, error)
	GetTransactionByID(ctx context.Context, id int) (*Transaction, error)
	UpdateTransaction(ctx context.Context, transaction *Transaction) error
	DeleteTransaction(ctx context.Context, id int) error
}

type repository struct{ *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db}
}

func (r *repository) Save(ctx context.Context, transaction *Transaction) error {
	query := `INSERT INTO transactions (amount, transaction_type, description, category_id, account_id, transaction_date, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	var id int

	err := r.QueryRow(ctx, query,
		transaction.Amount,
		transaction.TransactionType,
		transaction.Description,
		transaction.CategoryID,
		transaction.AccountID,
		transaction.TransactionDate,
		transaction.CreatedAt,
		transaction.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return err
	}

	transaction.ID = id
	return nil
}

func (r *repository) GetTransactions(ctx context.Context, userID string) ([]Transaction, error) {
	query := `SELECT t.id, t.amount, t.transaction_type, t.description, t.category_id, t.account_id, t.transaction_date, t.created_at, t.updated_at FROM transactions t 
	JOIN accounts a ON t.account_id = a.id WHERE a.user_id = $1`
	rows, err := r.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions = make([]Transaction, 0)
	for rows.Next() {
		var transaction Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.Amount,
			&transaction.TransactionType,
			&transaction.Description,
			&transaction.CategoryID,
			&transaction.AccountID,
			&transaction.TransactionDate,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *repository) GetTransactionByID(ctx context.Context, id int) (*Transaction, error) {
	query := "SELECT id, amount, transaction_type, description, category_id, account_id, transaction_date, created_at, updated_at FROM transactions WHERE id = $1"
	row := r.QueryRow(ctx, query, id)

	var transaction Transaction
	err := row.Scan(
		&transaction.ID,
		&transaction.Amount,
		&transaction.TransactionType,
		&transaction.Description,
		&transaction.CategoryID,
		&transaction.AccountID,
		&transaction.TransactionDate,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *repository) UpdateTransaction(ctx context.Context, transaction *Transaction) error {
	query := "UPDATE transactions SET amount = $1, transaction_type = $2, description = $3, category_id = $4, account_id = $5, updated_at = $6 WHERE id = $7"
	_, err := r.Exec(ctx, query, transaction.Amount, transaction.TransactionType, transaction.Description, transaction.CategoryID, transaction.AccountID, transaction.UpdatedAt, transaction.ID)
	return err
}

func (r *repository) DeleteTransaction(ctx context.Context, id int) error {
	query := "DELETE FROM transactions WHERE id = $1"
	_, err := r.Exec(ctx, query, id)
	return err
}
