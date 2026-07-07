package summary

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetSummaryByModule(ctx context.Context, userID string, module Module, from time.Time, to time.Time) ([]TransactionCategoryBalance, error)
}

type repository struct{ *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db}
}

func (r *repository) GetSummaryByModule(ctx context.Context, userID string, module Module, from time.Time, to time.Time) ([]TransactionCategoryBalance, error) {
	var query string
	if module == ModuleCredit {
		query = `SELECT
			c.name AS category_name,
			COALESCE(SUM(CASE WHEN t.transaction_type = 'credit' THEN -t.amount ELSE 0 END), 0) AS balance
		FROM transactions t
		JOIN accounts a ON t.account_id = a.id
		JOIN categories c ON c.id = t.category_id
		WHERE a.user_id = $1
		  AND t.transaction_date BETWEEN $2 AND $3
		GROUP BY c.name`
	} else {
		query = `SELECT
			c.name AS category_name,
			COALESCE(SUM(CASE WHEN t.transaction_type = 'debit' THEN t.amount ELSE 0 END), 0) AS balance
		FROM transactions t
		JOIN accounts a ON t.account_id = a.id
		JOIN categories c ON c.id = t.category_id
		WHERE a.user_id = $1
		  AND t.transaction_date BETWEEN $2 AND $3
		GROUP BY c.name`
	}

	rows, err := r.Query(ctx, query, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	balances := make([]TransactionCategoryBalance, 0)
	for rows.Next() {
		var item TransactionCategoryBalance
		if err := rows.Scan(&item.CategoryName, &item.Balance); err != nil {
			return nil, err
		}
		balances = append(balances, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return balances, nil
}
