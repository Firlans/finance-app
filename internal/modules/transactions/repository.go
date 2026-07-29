package transactions

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Save(ctx context.Context, transaction *Transaction) error
	GetTransactions(ctx context.Context, userID string, from string, to string, page int) ([]Transaction, error)
	GetTransactionByID(ctx context.Context, id int) (*Transaction, error)
	UpdateTransaction(ctx context.Context, transaction *Transaction) error
	DeleteTransaction(ctx context.Context, id int) error
	IsTransactionLinkedToPayment(ctx context.Context, transactionID int) (bool, error)
	CreateTransfer(ctx context.Context, userID string, req *CreateTransferRequest) error
	GetTransferByID(ctx context.Context, userID string, id int) (*CreateTransferRequest, error)
	UpdateTransfer(ctx context.Context, userID string, id int, req *CreateTransferRequest) error
	GetCategoryIDByName(ctx context.Context, name string) (*int, error)
}

type repository struct{ *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db}
}

func (r *repository) Save(ctx context.Context, transaction *Transaction) error {
	query := `INSERT INTO transactions (amount, transaction_type, description, category_id, account_id, transaction_date, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	var id int

	err := r.QueryRow(ctx, query,
		transaction.Amount,
		transaction.TransactionType,
		transaction.Description,
		transaction.CategoryID,
		transaction.AccountID,
		transaction.TransactionDate,
		transaction.CreatedAt,
		transaction.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return err
	}

	transaction.ID = id
	return nil
}

func (r *repository) GetTransactions(ctx context.Context, userID string, from string, to string, page int) ([]Transaction, error) {
	query := `SELECT t.id, t.amount, t.transaction_type, t.description, t.category_id, t.account_id, t.transaction_date, t.created_at, t.updated_at, EXISTS(SELECT 1 FROM payments p WHERE p.transaction_id = t.id) as is_loan, t.related_transaction_id FROM transactions t 
	JOIN accounts a ON t.account_id = a.id WHERE a.user_id = $1`
	var args []interface{}
	args = append(args, userID)
	argID := 2 // Argumen berikutnya dimulai dari index 2 ($2)

	// Jika filter 'from' diberikan, tambahkan ke query
	if from != "" {
		query += fmt.Sprintf(" AND t.transaction_date >= $%d", argID)
		args = append(args, from)
		argID++
	}

	// Jika filter 'to' diberikan, tambahkan ke query
	if to != "" {
		query += fmt.Sprintf(" AND t.transaction_date <= $%d", argID)
		args = append(args, to)
		argID++
	}

	// Tambahkan ORDER BY agar data transaksi selalu urut dari yang paling baru
	query += " ORDER BY t.transaction_date DESC"

	if page > 0 {
		limit := 100
		offset := (page - 1) * limit
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	// Gunakan args... agar semua parameter yang di-append masuk ke Query
	rows, err := r.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions = make([]Transaction, 0)
	for rows.Next() {
		var transaction Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.Amount,
			&transaction.TransactionType,
			&transaction.Description,
			&transaction.CategoryID,
			&transaction.AccountID,
			&transaction.TransactionDate,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
			&transaction.IsLoan,
			&transaction.RelatedTransactionID,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *repository) GetTransactionByID(ctx context.Context, id int) (*Transaction, error) {
	query := "SELECT id, amount, transaction_type, description, category_id, account_id, transaction_date, created_at, updated_at, EXISTS(SELECT 1 FROM payments p WHERE p.transaction_id = transactions.id) as is_loan, related_transaction_id FROM transactions WHERE id = $1"
	row := r.QueryRow(ctx, query, id)

	var transaction Transaction
	err := row.Scan(
		&transaction.ID,
		&transaction.Amount,
		&transaction.TransactionType,
		&transaction.Description,
		&transaction.CategoryID,
		&transaction.AccountID,
		&transaction.TransactionDate,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
		&transaction.IsLoan,
		&transaction.RelatedTransactionID,
	)
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *repository) UpdateTransaction(ctx context.Context, transaction *Transaction) error {
	query := "UPDATE transactions SET amount = $1, transaction_type = $2, description = $3, category_id = $4, account_id = $5, updated_at = $6, transaction_date = $7 WHERE id = $8"
	_, err := r.Exec(ctx, query, transaction.Amount, transaction.TransactionType, transaction.Description, transaction.CategoryID, transaction.AccountID, transaction.UpdatedAt, transaction.TransactionDate, transaction.ID)
	return err
}

func (r *repository) DeleteTransaction(ctx context.Context, id int) error {
	query := `DELETE FROM transactions 
	WHERE id = $1 
	   OR related_transaction_id = $1 
	   OR id = (SELECT related_transaction_id FROM transactions WHERE id = $1 AND related_transaction_id IS NOT NULL)`
	_, err := r.Exec(ctx, query, id)
	return err
}


func (r *repository) IsTransactionLinkedToPayment(ctx context.Context, transactionID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM payments WHERE transaction_id = $1)`
	var exists bool
	err := r.QueryRow(ctx, query, transactionID).Scan(&exists)
	return exists, err
}

func (r *repository) GetCategoryIDByName(ctx context.Context, name string) (*int, error) {
	query := `SELECT id FROM categories WHERE LOWER(name) = LOWER($1) AND user_id IS NULL LIMIT 1`
	var id int
	err := r.QueryRow(ctx, query, name).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *repository) CreateTransfer(ctx context.Context, userID string, req *CreateTransferRequest) error {
	tx, err := r.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Get Category IDs
	var transferCatID *int
	var feeCatID *int

	var id1 int
	err = tx.QueryRow(ctx, `SELECT id FROM categories WHERE LOWER(name) = 'transfer' AND user_id IS NULL LIMIT 1`).Scan(&id1)
	if err == nil {
		transferCatID = &id1
	}

	if req.AdminFee > 0 {
		var id2 int
		err = tx.QueryRow(ctx, `SELECT id FROM categories WHERE LOWER(name) = 'fee' AND user_id IS NULL LIMIT 1`).Scan(&id2)
		if err == nil {
			feeCatID = &id2
		}
	}

	// 2. Fetch Account Names for Description & Verify Ownership
	var fromAccountName string
	err = tx.QueryRow(ctx, `SELECT account_name FROM accounts WHERE id = $1 AND user_id = $2`, req.FromAccountID, userID).Scan(&fromAccountName)
	if err != nil {
		return fmt.Errorf("akun asal tidak ditemukan atau bukan milik Anda")
	}

	var toAccountName string
	err = tx.QueryRow(ctx, `SELECT account_name FROM accounts WHERE id = $1 AND user_id = $2`, req.ToAccountID, userID).Scan(&toAccountName)
	if err != nil {
		return fmt.Errorf("akun tujuan tidak ditemukan atau bukan milik Anda")
	}

	now := time.Now().UTC()
	descFrom := fmt.Sprintf("Transfer ke %s", toAccountName)
	if req.Description != "" {
		descFrom = fmt.Sprintf("Transfer ke %s - %s", toAccountName, req.Description)
	}

	descTo := fmt.Sprintf("Transfer dari %s", fromAccountName)
	if req.Description != "" {
		descTo = fmt.Sprintf("Transfer dari %s - %s", fromAccountName, req.Description)
	}

	// 3. Insert Transaction 1 (Credit at From Account)
	queryInsert := `INSERT INTO transactions (amount, transaction_type, description, category_id, account_id, transaction_date, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	var tx1ID int
	err = tx.QueryRow(ctx, queryInsert, req.Amount, "credit", descFrom, transferCatID, req.FromAccountID, req.TransactionDate, now, now).Scan(&tx1ID)
	if err != nil {
		return err
	}

	// 4. Insert Transaction 2 (Debit at To Account, linked to Tx1)
	queryInsertWithRelated := `INSERT INTO transactions (amount, transaction_type, description, category_id, account_id, transaction_date, created_at, updated_at, related_transaction_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	var tx2ID int
	err = tx.QueryRow(ctx, queryInsertWithRelated, req.Amount, "debit", descTo, transferCatID, req.ToAccountID, req.TransactionDate, now, now, tx1ID).Scan(&tx2ID)
	if err != nil {
		return err
	}

	// Link Tx1 to Tx2
	_, err = tx.Exec(ctx, `UPDATE transactions SET related_transaction_id = $1 WHERE id = $2`, tx2ID, tx1ID)
	if err != nil {
		return err
	}

	// 5. Insert Transaction 3 (Fee Credit at From Account, if AdminFee > 0)
	if req.AdminFee > 0 {
		descFee := fmt.Sprintf("Biaya Admin Transfer ke %s", toAccountName)
		if req.Description != "" {
			descFee = fmt.Sprintf("Biaya Admin Transfer ke %s - %s", toAccountName, req.Description)
		}
		_, err = tx.Exec(ctx, queryInsertWithRelated, req.AdminFee, "credit", descFee, feeCatID, req.FromAccountID, req.TransactionDate, now, now, tx1ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *repository) GetTransferByID(ctx context.Context, userID string, id int) (*CreateTransferRequest, error) {
	tx, err := r.GetTransactionByID(ctx, id)
	if err != nil || tx == nil {
		return nil, err
	}

	var mainTxID int = tx.ID
	var relatedID int
	if tx.RelatedTransactionID != nil {
		relatedID = *tx.RelatedTransactionID
	}

	var fromAccountID, toAccountID int
	var amount float64
	var adminFee float64
	var description string
	var txDate time.Time

	if tx.TransactionType == "credit" {
		fromAccountID = tx.AccountID
		amount = tx.Amount
		txDate = tx.TransactionDate
		description = tx.Description

		if relatedID > 0 {
			_ = r.QueryRow(ctx, `SELECT account_id FROM transactions WHERE id = $1 AND transaction_type = 'debit'`, relatedID).Scan(&toAccountID)
			_ = r.QueryRow(ctx, `SELECT amount FROM transactions WHERE related_transaction_id = $1 AND category_id = (SELECT id FROM categories WHERE LOWER(name) = 'fee' AND user_id IS NULL LIMIT 1)`, mainTxID).Scan(&adminFee)
		}
	} else if tx.TransactionType == "debit" {
		toAccountID = tx.AccountID
		amount = tx.Amount
		txDate = tx.TransactionDate
		description = tx.Description

		if relatedID > 0 {
			_ = r.QueryRow(ctx, `SELECT account_id, description FROM transactions WHERE id = $1 AND transaction_type = 'credit'`, relatedID).Scan(&fromAccountID, &description)
			_ = r.QueryRow(ctx, `SELECT amount FROM transactions WHERE related_transaction_id = $1 AND category_id = (SELECT id FROM categories WHERE LOWER(name) = 'fee' AND user_id IS NULL LIMIT 1)`, relatedID).Scan(&adminFee)
		}
	}

	cleanDesc := description
	if idx := strings.Index(description, " - "); idx != -1 {
		cleanDesc = description[idx+3:]
	} else if strings.HasPrefix(description, "Transfer ke ") || strings.HasPrefix(description, "Transfer dari ") {
		cleanDesc = ""
	}

	res := &CreateTransferRequest{
		FromAccountID:   fromAccountID,
		ToAccountID:     toAccountID,
		Amount:          amount,
		AdminFee:        adminFee,
		Description:     cleanDesc,
		TransactionDate: txDate,
	}

	return res, nil
}

func (r *repository) UpdateTransfer(ctx context.Context, userID string, id int, req *CreateTransferRequest) error {
	tx, err := r.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var tx1ID, tx2ID int
	var txType string
	var relatedID *int

	err = tx.QueryRow(ctx, `SELECT t.id, t.transaction_type, t.related_transaction_id FROM transactions t JOIN accounts a ON t.account_id = a.id WHERE t.id = $1 AND a.user_id = $2`, id, userID).Scan(&tx1ID, &txType, &relatedID)
	if err != nil {
		return fmt.Errorf("transaksi transfer tidak ditemukan: %w", err)
	}

	if txType == "debit" && relatedID != nil {
		tx2ID = tx1ID
		tx1ID = *relatedID
	} else if relatedID != nil {
		tx2ID = *relatedID
	}

	var transferCatID *int
	var feeCatID *int
	var id1, id2 int
	if err := tx.QueryRow(ctx, `SELECT id FROM categories WHERE LOWER(name) = 'transfer' AND user_id IS NULL LIMIT 1`).Scan(&id1); err == nil {
		transferCatID = &id1
	}
	if err := tx.QueryRow(ctx, `SELECT id FROM categories WHERE LOWER(name) = 'fee' AND user_id IS NULL LIMIT 1`).Scan(&id2); err == nil {
		feeCatID = &id2
	}

	var fromAccountName, toAccountName string
	if err := tx.QueryRow(ctx, `SELECT account_name FROM accounts WHERE id = $1 AND user_id = $2`, req.FromAccountID, userID).Scan(&fromAccountName); err != nil {
		return fmt.Errorf("akun asal tidak ditemukan")
	}
	if err := tx.QueryRow(ctx, `SELECT account_name FROM accounts WHERE id = $1 AND user_id = $2`, req.ToAccountID, userID).Scan(&toAccountName); err != nil {
		return fmt.Errorf("akun tujuan tidak ditemukan")
	}

	now := time.Now().UTC()
	descFrom := fmt.Sprintf("Transfer ke %s", toAccountName)
	if req.Description != "" {
		descFrom = fmt.Sprintf("Transfer ke %s - %s", toAccountName, req.Description)
	}

	descTo := fmt.Sprintf("Transfer dari %s", fromAccountName)
	if req.Description != "" {
		descTo = fmt.Sprintf("Transfer dari %s - %s", fromAccountName, req.Description)
	}

	_, err = tx.Exec(ctx, `UPDATE transactions SET amount = $1, account_id = $2, description = $3, category_id = $4, transaction_date = $5, updated_at = $6 WHERE id = $7`,
		req.Amount, req.FromAccountID, descFrom, transferCatID, req.TransactionDate, now, tx1ID)
	if err != nil {
		return err
	}

	if tx2ID > 0 {
		_, err = tx.Exec(ctx, `UPDATE transactions SET amount = $1, account_id = $2, description = $3, category_id = $4, transaction_date = $5, updated_at = $6 WHERE id = $7`,
			req.Amount, req.ToAccountID, descTo, transferCatID, req.TransactionDate, now, tx2ID)
		if err != nil {
			return err
		}
	}

	var feeTxID int
	_ = tx.QueryRow(ctx, `SELECT id FROM transactions WHERE related_transaction_id = $1 AND category_id = $2`, tx1ID, feeCatID).Scan(&feeTxID)

	if req.AdminFee > 0 {
		descFee := fmt.Sprintf("Biaya Admin Transfer ke %s", toAccountName)
		if req.Description != "" {
			descFee = fmt.Sprintf("Biaya Admin Transfer ke %s - %s", toAccountName, req.Description)
		}
		if feeTxID > 0 {
			_, err = tx.Exec(ctx, `UPDATE transactions SET amount = $1, account_id = $2, description = $3, transaction_date = $4, updated_at = $5 WHERE id = $6`,
				req.AdminFee, req.FromAccountID, descFee, req.TransactionDate, now, feeTxID)
			if err != nil {
				return err
			}
		} else {
			queryInsertWithRelated := `INSERT INTO transactions (amount, transaction_type, description, category_id, account_id, transaction_date, created_at, updated_at, related_transaction_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
			_, err = tx.Exec(ctx, queryInsertWithRelated, req.AdminFee, "credit", descFee, feeCatID, req.FromAccountID, req.TransactionDate, now, now, tx1ID)
			if err != nil {
				return err
			}
		}
	} else if feeTxID > 0 {
		_, err = tx.Exec(ctx, `DELETE FROM transactions WHERE id = $1`, feeTxID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}


