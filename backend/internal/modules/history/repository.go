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
	FindManyByBudgetID(ctx context.Context, budgetID string) ([]*History, error)
	IsBudgetOwnedByUser(ctx context.Context, budgetID string, userID string) (bool, error)
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
		INSERT INTO histories (id, budget_id, amount, date, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query, h.ID, h.BudgetID, h.Amount, h.Date, h.CreatedAt)
	return err
}

func (r *repository) Update(ctx context.Context, h *History) error {
	query := `
		UPDATE histories 
		SET amount = $1, date = $2
		WHERE id = $3
	`
	cmdTag, err := r.db.Exec(ctx, query, h.Amount, h.Date, h.ID)
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
		SELECT id, budget_id, amount, date, created_at
		FROM histories
		WHERE id = $1
	`

	var h History
	err := r.db.QueryRow(ctx, query, id).Scan(&h.ID, &h.BudgetID, &h.Amount, &h.Date, &h.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &h, nil
}

func (r *repository) FindManyByBudgetID(ctx context.Context, budgetID string) ([]*History, error) {
	query := `
		SELECT id, budget_id, amount, date, created_at
		FROM histories
		WHERE budget_id = $1
		ORDER BY date DESC
	`
	rows, err := r.db.Query(ctx, query, budgetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []*History
	for rows.Next() {
		var h History
		if err := rows.Scan(&h.ID, &h.BudgetID, &h.Amount, &h.Date, &h.CreatedAt); err != nil {
			return nil, err
		}
		histories = append(histories, &h)
	}
	return histories, nil
}

func (r *repository) IsBudgetOwnedByUser(ctx context.Context, budgetID string, userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM monthly_budgets WHERE id = $1 AND user_id = $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, budgetID, userID).Scan(&exists)
	return exists, err
}

func (r *repository) IsHistoryOwnedByUser(ctx context.Context, historyID string, userID string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM histories h
			INNER JOIN monthly_budgets b ON h.budget_id = b.id
			WHERE h.id = $1 AND b.user_id = $2
		)
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, historyID, userID).Scan(&exists)
	return exists, err
}
