package goals

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, goal *Goal) error
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]Goal, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Goal, error)
	Update(ctx context.Context, goal *Goal) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, goal *Goal) error {
	query := `
		INSERT INTO goals (id, user_id, name, target_amount, current_amount)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query, goal.ID, goal.UserID, goal.Name, goal.TargetAmount, goal.CurrentAmount)
	return err
}

func (r *repository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]Goal, error) {
	query := `
		SELECT id, user_id, name, target_amount, current_amount, created_at, updated_at
		FROM goals
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []Goal
	for rows.Next() {
		var g Goal
		if err := rows.Scan(&g.ID, &g.UserID, &g.Name, &g.TargetAmount, &g.CurrentAmount, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, err
		}
		goals = append(goals, g)
	}

	return goals, nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Goal, error) {
	query := `
		SELECT id, user_id, name, target_amount, current_amount, created_at, updated_at
		FROM goals
		WHERE id = $1 AND user_id = $2
	`
	var g Goal
	err := r.db.QueryRow(ctx, query, id, userID).Scan(&g.ID, &g.UserID, &g.Name, &g.TargetAmount, &g.CurrentAmount, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *repository) Update(ctx context.Context, goal *Goal) error {
	query := `
		UPDATE goals
		SET name = $1, target_amount = $2, current_amount = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 AND user_id = $5
	`
	_, err := r.db.Exec(ctx, query, goal.Name, goal.TargetAmount, goal.CurrentAmount, goal.ID, goal.UserID)
	return err
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM goals WHERE id = $1 AND user_id = $2", id, userID)
	return err
}
