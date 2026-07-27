package loans

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	SaveLoan(ctx context.Context, loan *CreateLoanRequest) error
	FindLoanByID(ctx context.Context, id int) (*Loan, error)
	FindLoansByUserID(ctx context.Context, userID string) (*[]Loan, error)
	UpdateLoan(ctx context.Context, id int, loan *UpdateLoanRequest) error
	DeleteLoan(ctx context.Context, id int) error
}

type repository struct{ *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db}
}

// ===== Loan Operations =====

// SaveLoan creates a new loan
func (r *repository) SaveLoan(ctx context.Context, loan *CreateLoanRequest) error {
	query := `INSERT INTO loans (user_id, name, loan_type, created_at, updated_at)
	VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`

	var id int64
	err := r.QueryRow(ctx, query, loan.UserID, loan.Name, loan.LoanType).Scan(&id)
	if err != nil {
		return err
	}

	loan.ID = int(id)
	return nil
}

// FindLoanByID retrieves a loan by its ID
func (r *repository) FindLoanByID(ctx context.Context, id int) (*Loan, error) {
	query := `
		SELECT 
			loans.id, 
			loans.user_id, 
			loans.name, 
			loans.loan_type,
			COALESCE(SUM(CASE WHEN payments.type = 'increase' THEN payments.amount ELSE 0 END), 0) - COALESCE(SUM(CASE WHEN payments.type = 'decrease' THEN payments.amount ELSE 0 END), 0) as outstanding_amount,
			loans.created_at, 
			loans.updated_at
		FROM loans 
		LEFT JOIN payments ON loans.id = payments.loan_id
		WHERE loans.id = $1
		GROUP BY loans.id`
	row := r.QueryRow(ctx, query, id)

	var loan Loan
	err := row.Scan(&loan.ID, &loan.UserID, &loan.Name, &loan.LoanType, &loan.OutstandingAmount, &loan.CreatedAt, &loan.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &loan, nil
}

// FindLoansByUserID retrieves all loans for a specific user
func (r *repository) FindLoansByUserID(ctx context.Context, userID string) (*[]Loan, error) {
	query := `SELECT 
				loans.id, 
				loans.user_id, 
				loans.name, 
				loans.loan_type,
				COALESCE(SUM(CASE WHEN payments.type = 'increase' THEN payments.amount ELSE 0 END), 0) - COALESCE(SUM(CASE WHEN payments.type = 'decrease' THEN payments.amount ELSE 0 END), 0) as outstanding_amount,
				loans.created_at, 
				loans.updated_at 
			FROM loans
			LEFT JOIN payments ON loans.id = payments.loan_id
			WHERE loans.user_id = $1 
			GROUP BY loans.id
			ORDER BY loans.created_at DESC`
	loansList := make([]Loan, 0)
	rows, err := r.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var loan Loan
		err := rows.Scan(&loan.ID, &loan.UserID, &loan.Name, &loan.LoanType, &loan.OutstandingAmount, &loan.CreatedAt, &loan.UpdatedAt)
		if err != nil {
			return nil, err
		}
		loansList = append(loansList, loan)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &loansList, nil
}

// UpdateLoan updates an existing loan
func (r *repository) UpdateLoan(ctx context.Context, id int, loan *UpdateLoanRequest) error {
	query := `UPDATE loans SET name = COALESCE($1, name), loan_type = COALESCE($2, loan_type), updated_at = NOW() WHERE id = $3`
	_, err := r.Exec(ctx, query, loan.Name, loan.LoanType, id)
	if err != nil {
		return err
	}
	return nil
}

// DeleteLoan deletes a loan by its ID
func (r *repository) DeleteLoan(ctx context.Context, id int) error {
	query := `DELETE FROM loans WHERE id = $1`
	_, err := r.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
