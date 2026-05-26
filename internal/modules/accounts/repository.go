package accounts

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Save(ctx context.Context, account *CreateAccountRequest) error
	FindByID(ctx context.Context, id int) (*Account, error)
	FindByUserID(ctx context.Context, userID string) (*[]Account, error)
	Update(ctx context.Context, account *CreateAccountRequest) error
	Delete(ctx context.Context, id int) error
}

type repository struct{ *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db}
}

func (r *repository) Save(ctx context.Context, account *CreateAccountRequest) error {
	query := `INSERT INTO accounts (user_id, account_name, balance, description)
	VALUES ($1, $2, $3, $4) RETURNING id`

	var id int64
	err := r.QueryRow(ctx, query, account.UserID, account.AccountName, account.Balance, account.Description).Scan(&id)
	if err != nil {
		return err
	}

	account.ID = int(id)
	return nil
}

func (r *repository) FindByID(ctx context.Context, id int) (*Account, error) {
	query := `
		select 	a.id, 
				a.user_id, 
				a.account_name, 
				a.description,
				a.balance + coalesce(sum(
					case
						when t.transaction_type = 'debit' then t.amount
						when t.transaction_type = 'credit' then -t.amount
						else 0
					end
				), 0) as balance,
				a.balance as initial_balance,
				a.created_at 
		from accounts a 
		left join transactions t on a.id = t.account_id
		where a.id = $1`
	row := r.QueryRow(ctx, query, id)
	var account Account
	err := row.Scan(&account.ID, &account.UserID, &account.AccountName, &account.Description, &account.Balance, &account.InitialBalance, &account.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *repository) FindByUserID(ctx context.Context, userID string) (*[]Account, error) {
	query := `
	select 	a.id, 
			a.user_id, 
			a.account_name, 
			a.description, 
			a.balance + coalesce(sum(
				case
					when t.transaction_type = 'debit' then t.amount
					when t.transaction_type = 'credit' then -t.amount
					else 0
				end
			), 0) as balance,
			a.balance as initial_balance,
			a.created_at
	from accounts a
	left join transactions t on a.id = t.account_id
	where a.user_id = $1
	group by a.id, a.user_id, a.account_name, a.description, a.balance, a.created_at
	`
	accounts := make([]Account, 0)
	rows, err := r.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var account Account
		err := rows.Scan(&account.ID, &account.UserID, &account.AccountName, &account.Description, &account.Balance, &account.InitialBalance, &account.CreatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &accounts, nil
}

func (r *repository) Update(ctx context.Context, account *CreateAccountRequest) error {
	query := "UPDATE accounts SET account_name = $1, description = $2, balance = $3 WHERE id = $4"
	_, err := r.Exec(ctx, query, account.AccountName, account.Description, account.Balance, account.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM accounts WHERE id = $1"
	_, err := r.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
