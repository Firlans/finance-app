package budget

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, budget *MonthlyBudget) error
	Update(ctx context.Context, budget *MonthlyBudget) error
	FindByID(ctx context.Context, id string) (*MonthlyBudget, error)
	FindMany(ctx context.Context, userID string, startDate, endDate time.Time) ([]*MonthlyBudget, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r repository) Create(ctx context.Context, b *MonthlyBudget) error {
	query := `
		INSERT INTO budgets (id, user_id, budget, start_date, iteration) 
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query, b.ID, b.UserID, b.Budget, b.StartDate, b.Iteration)
	return err
}

func (r *repository) Update(ctx context.Context, b *MonthlyBudget) error {
	query := `
		UPDATE budgets 
		SET budget = $1, start_date = $2
		WHERE id = $3
	`
	cmdTag, err := r.db.Exec(ctx, query, b.Budget, b.StartDate, b.ID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrBudgetNotFound
	}

	return nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*MonthlyBudget, error) {
	query := `
		SELECT id, user_id, budget, start_date, iteration
		FROM budgets
		WHERE id = $1
	`

	var b MonthlyBudget
	err := r.db.QueryRow(ctx, query, id).Scan(&b.ID, &b.UserID, &b.Budget, &b.StartDate, &b.Iteration)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &b, nil
}

func (r *repository) FindMany(ctx context.Context, userID string, startDate, endDate time.Time) ([]*MonthlyBudget, error) {
	query := `
		SELECT id, user_id, budget, start_date, iteration
		FROM budgets
		WHERE user_id = $1 AND start_date >= $2 AND start_date <= $3
		ORDER BY start_date DESC
	`

	rows, err := r.db.Query(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var budgets []*MonthlyBudget
	for rows.Next() {
		var b MonthlyBudget
		if err := rows.Scan(&b.ID, &b.UserID, &b.Budget, &b.StartDate, &b.Iteration); err != nil {
			return nil, err
		}
		budgets = append(budgets, &b)
	}

	return budgets, nil
}
