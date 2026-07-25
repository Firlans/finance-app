package goals

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UseCase interface {
	CreateGoal(ctx context.Context, userID uuid.UUID, req *CreateGoalRequest) error
	GetGoals(ctx context.Context, userID uuid.UUID) ([]GoalResponse, error)
	UpdateGoal(ctx context.Context, userID uuid.UUID, id uuid.UUID, req *UpdateGoalRequest) error
	DeleteGoal(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
}

type useCase struct {
	repo     Repository
	validate *validator.Validate
}

func NewUseCase(repo Repository, validate *validator.Validate) UseCase {
	return &useCase{
		repo:     repo,
		validate: validate,
	}
}

func (u *useCase) CreateGoal(ctx context.Context, userID uuid.UUID, req *CreateGoalRequest) error {
	if err := u.validate.Struct(req); err != nil {
		return err
	}

	goal := &Goal{
		ID:            uuid.New(),
		UserID:        userID,
		Name:          req.Name,
		TargetAmount:  req.TargetAmount,
		CurrentAmount: req.CurrentAmount,
	}

	return u.repo.Create(ctx, goal)
}

func (u *useCase) GetGoals(ctx context.Context, userID uuid.UUID) ([]GoalResponse, error) {
	goals, err := u.repo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var responses []GoalResponse
	for _, g := range goals {
		responses = append(responses, GoalResponse{
			ID:            g.ID,
			Name:          g.Name,
			TargetAmount:  g.TargetAmount,
			CurrentAmount: g.CurrentAmount,
			IsCompleted:   g.CurrentAmount >= g.TargetAmount,
		})
	}

	return responses, nil
}

func (u *useCase) UpdateGoal(ctx context.Context, userID uuid.UUID, id uuid.UUID, req *UpdateGoalRequest) error {
	if err := u.validate.Struct(req); err != nil {
		return err
	}

	goal := &Goal{
		ID:            id,
		UserID:        userID,
		Name:          req.Name,
		TargetAmount:  req.TargetAmount,
		CurrentAmount: req.CurrentAmount,
	}

	return u.repo.Update(ctx, goal)
}

func (u *useCase) DeleteGoal(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	return u.repo.Delete(ctx, id, userID)
}
