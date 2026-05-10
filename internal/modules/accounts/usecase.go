package accounts

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
)

type UseCase interface {
	Save(ctx context.Context, account *CreateAccountRequest) (int, error)
	getAccounts(ctx context.Context, userID string) (*[]Account, error)
	Update(ctx context.Context, account *CreateAccountRequest) (int, error)
	DeleteAccount(ctx context.Context, id int) error
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

func (uc *useCase) Save(ctx context.Context, account *CreateAccountRequest) (int, error) {
	if account == nil {
		return 0, errors.New("account is nil")
	}

	if err := uc.validate.Struct(account); err != nil {
		return 0, err
	}

	newAccount := &CreateAccountRequest{
		UserID:      account.UserID,
		AccountName: account.AccountName,
		Balance:     account.Balance,
		Description: account.Description,
	}

	if err := uc.repo.Save(ctx, newAccount); err != nil {
		return 0, err
	}

	return newAccount.ID, nil
}

func (uc *useCase) getAccounts(ctx context.Context, userID string) (*[]Account, error) {
	var accounts *[]Account
	accounts, err := uc.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (uc *useCase) Update(ctx context.Context, account *CreateAccountRequest) (int, error) {
	if account == nil {
		return 0, errors.New("account is nil")
	}

	if err := uc.validate.Struct(account); err != nil {
		return 0, err
	}

	if err := uc.repo.Update(ctx, account); err != nil {
		return 0, err
	}

	return account.ID, nil
}

func (uc *useCase) DeleteAccount(ctx context.Context, id int) error {
	account, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := uc.repo.Delete(ctx, account.ID); err != nil {
		return err
	}

	return nil
}
