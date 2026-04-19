package transaction

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, transaction *Transaction) error
	Update(ctx context.Context, transaction *Transaction) error
	FindByID(ctx context.Context, id string) (*Transaction, error)
	FindManyByAccountID(ctx context.Context, accountID string) ([]*Transaction, error)
	IsAccountOwnedByUser(ctx context.Context, accountID string, userID string) (bool, error)
	IsTransactionOwnedByUser(ctx context.Context, transactionID string, userID string) (bool, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, t *Transaction) error {
	query := `
		INSERT INTO transactions (id, user_id, account_id, amount, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query, t.ID, t.UserID, t.AccountID, t.Amount, t.CreatedAt)
	return err
}

func (r *repository) Update(ctx context.Context, t *Transaction) error {
	query := `
		UPDATE transactions 
		SET amount = $1, updated_at = $2
		WHERE id = $3
	`
	cmdTag, err := r.db.Exec(ctx, query, t.Amount, t.UpdatedAt, t.ID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrTransactionNotFound
	}

	return nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*Transaction, error) {
	query := `
		SELECT id, user_id, account_id, amount, created_at, updated_at
		FROM transactions
		WHERE id = $1
	`

	var t Transaction
	err := r.db.QueryRow(ctx, query, id).Scan(&t.ID, &t.UserID, &t.AccountID, &t.Amount, &t.CreatedAt, &t.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &t, nil
}

func (r *repository) FindManyByAccountID(ctx context.Context, accountID string) ([]*Transaction, error) {
	query := `
		SELECT id, user_id, account_id, amount, created_at, updated_at
		FROM transactions
		WHERE account_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.AccountID, &t.Amount, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, &t)
	}
	return transactions, nil
}

func (r *repository) IsAccountOwnedByUser(ctx context.Context, accountID string, userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM transactions WHERE account_id = $1 AND user_id = $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, accountID, userID).Scan(&exists)
	return exists, err
}

func (r *repository) IsTransactionOwnedByUser(ctx context.Context, transactionID string, userID string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM transactions t
			WHERE t.id = $1 AND t.user_id = $2
		)
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, transactionID, userID).Scan(&exists)
	return exists, err
}
