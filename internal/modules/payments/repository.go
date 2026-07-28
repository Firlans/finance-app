package payments

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	SavePayment(ctx context.Context, payment *CreatePaymentRequest) error
	FindPaymentByID(ctx context.Context, id int) (*Payment, error)
	FindPaymentsByLoanID(ctx context.Context, loanID int) (*[]Payment, error)
	FindPaymentsByTransactionID(ctx context.Context, transactionID int) (*[]Payment, error)
	FindFirstPaymentByLoanID(ctx context.Context, loanID int) (*Payment, error)
	UpdatePayment(ctx context.Context, payment *Payment) error
	DeletePayment(ctx context.Context, id int) error
	GetLoanTypeByID(ctx context.Context, loanID int) (string, error)
	GetHutangCategoryID(ctx context.Context) (*int, error)
}

type repository struct{ *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db}
}

// SavePayment creates a new payment record
func (r *repository) SavePayment(ctx context.Context, payment *CreatePaymentRequest) error {
	query := `INSERT INTO payments (transaction_id, loan_id, amount, interest, payment_date, type, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW()) RETURNING id`

	var transactionID interface{}
	if payment.TransactionID != nil {
		transactionID = *payment.TransactionID
	}

	var id int64
	err := r.QueryRow(ctx, query, transactionID, payment.LoanID, payment.Amount, payment.Interest, payment.PaymentDate, payment.Type).Scan(&id)

	if err != nil {
		return err
	}

	payment.ID = int(id)
	return nil
}

// FindPaymentByID retrieves a payment by its ID
func (r *repository) FindPaymentByID(ctx context.Context, id int) (*Payment, error) {
	query := `SELECT id, transaction_id, loan_id, COALESCE(amount, 0), COALESCE(interest, 0), COALESCE(payment_date, created_at), COALESCE(type, 'decrease'), created_at, updated_at FROM payments WHERE id = $1`
	row := r.QueryRow(ctx, query, id)

	var payment Payment
	err := row.Scan(&payment.ID, &payment.TransactionID, &payment.LoanID, &payment.Amount, &payment.Interest, &payment.PaymentDate, &payment.Type, &payment.CreatedAt, &payment.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &payment, nil
}

// FindPaymentsByLoanID retrieves all payments for a specific loan
func (r *repository) FindPaymentsByLoanID(ctx context.Context, loanID int) (*[]Payment, error) {
	query := `SELECT id, transaction_id, loan_id, COALESCE(amount, 0), COALESCE(interest, 0), COALESCE(payment_date, created_at), COALESCE(type, 'decrease'), created_at, updated_at FROM payments WHERE loan_id = $1 ORDER BY created_at DESC`
	payments := make([]Payment, 0)
	rows, err := r.Query(ctx, query, loanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var payment Payment
		err := rows.Scan(&payment.ID, &payment.TransactionID, &payment.LoanID, &payment.Amount, &payment.Interest, &payment.PaymentDate, &payment.Type, &payment.CreatedAt, &payment.UpdatedAt)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &payments, nil
}

// FindPaymentsByTransactionID retrieves all payments for a specific transaction
func (r *repository) FindPaymentsByTransactionID(ctx context.Context, transactionID int) (*[]Payment, error) {
	query := `SELECT id, transaction_id, loan_id, COALESCE(amount, 0), COALESCE(interest, 0), COALESCE(payment_date, created_at), COALESCE(type, 'decrease'), created_at, updated_at FROM payments WHERE transaction_id = $1 ORDER BY created_at DESC`
	payments := make([]Payment, 0)
	rows, err := r.Query(ctx, query, transactionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var payment Payment
		err := rows.Scan(&payment.ID, &payment.TransactionID, &payment.LoanID, &payment.Amount, &payment.Interest, &payment.PaymentDate, &payment.Type, &payment.CreatedAt, &payment.UpdatedAt)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &payments, nil
}

func (r *repository) FindFirstPaymentByLoanID(ctx context.Context, loanID int) (*Payment, error) {
	query := `SELECT id, transaction_id, loan_id, COALESCE(amount, 0), COALESCE(interest, 0), COALESCE(payment_date, created_at), COALESCE(type, 'decrease'), created_at, updated_at FROM payments WHERE loan_id = $1 ORDER BY id ASC LIMIT 1`
	row := r.QueryRow(ctx, query, loanID)

	var payment Payment
	err := row.Scan(&payment.ID, &payment.TransactionID, &payment.LoanID, &payment.Amount, &payment.Interest, &payment.PaymentDate, &payment.Type, &payment.CreatedAt, &payment.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &payment, nil
}

func (r *repository) UpdatePayment(ctx context.Context, payment *Payment) error {
	query := `UPDATE payments SET transaction_id = $1, loan_id = $2, amount = $3, interest = $4, payment_date = $5, type = $6, updated_at = NOW() WHERE id = $7`

	var transactionID interface{}
	if payment.TransactionID != nil {
		transactionID = *payment.TransactionID
	}

	_, err := r.Exec(ctx, query, transactionID, payment.LoanID, payment.Amount, payment.Interest, payment.PaymentDate, payment.Type, payment.ID)
	return err

}

func (r *repository) DeletePayment(ctx context.Context, id int) error {
	query := `DELETE FROM payments WHERE id = $1`
	_, err := r.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetLoanTypeByID(ctx context.Context, loanID int) (string, error) {
	query := `SELECT loan_type FROM loans WHERE id = $1`
	var loanType string
	err := r.QueryRow(ctx, query, loanID).Scan(&loanType)
	return loanType, err
}

func (r *repository) GetHutangCategoryID(ctx context.Context) (*int, error) {
	query := `SELECT id FROM categories WHERE name = 'Hutang' AND user_id IS NULL LIMIT 1`
	var id int
	err := r.QueryRow(ctx, query).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}
