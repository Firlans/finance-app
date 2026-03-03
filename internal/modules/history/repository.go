package history

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, history *History) error
	Update(ctx context.Context, history *History) error
	FindByID(ctx context.Context, id string) (*History, error)
	FindManyByAccountID(ctx context.Context, accountID string) ([]*History, error)
	IsAccountOwnedByUser(ctx context.Context, accountID string, userID string) (bool, error)
	IsHistoryOwnedByUser(ctx context.Context, historyID string, userID string) (bool, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, h *History) error {
	query := `
		INSERT INTO transactions (id, user_id, account_id, amount, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query, h.ID, h.UserID, h.AccountID, h.Amount, h.CreatedAt)
	return err
}

func (r *repository) Update(ctx context.Context, h *History) error {
	query := `
		UPDATE transactions 
		SET amount = $1, updated_at = $2
		WHERE id = $3
	`
	cmdTag, err := r.db.Exec(ctx, query, h.Amount, h.UpdatedAt, h.ID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrHistoryNotFound
	}

	return nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*History, error) {
	query := `
		SELECT id, user_id, account_id, amount, created_at, updated_at
		FROM transactions
		WHERE id = $1
	`

	var h History
	err := r.db.QueryRow(ctx, query, id).Scan(&h.ID, &h.UserID, &h.AccountID, &h.Amount, &h.CreatedAt, &h.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &h, nil
}

func (r *repository) FindManyByAccountID(ctx context.Context, accountID string) ([]*History, error) {
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

	var histories []*History
	for rows.Next() {
		var h History
		if err := rows.Scan(&h.ID, &h.UserID, &h.AccountID, &h.Amount, &h.CreatedAt, &h.UpdatedAt); err != nil {
			return nil, err
		}
		histories = append(histories, &h)
	}
	return histories, nil
}

func (r *repository) IsAccountOwnedByUser(ctx context.Context, accountID string, userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM transactions WHERE account_id = $1 AND user_id = $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, accountID, userID).Scan(&exists)
	return exists, err
}

func (r *repository) IsHistoryOwnedByUser(ctx context.Context, historyID string, userID string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM transactions t
			WHERE t.id = $1 AND t.user_id = $2
		)
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, historyID, userID).Scan(&exists)
	return exists, err
}
