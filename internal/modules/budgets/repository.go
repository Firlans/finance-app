package budgets

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	UpsertBudget(ctx context.Context, budget *Budget, categoryIDs []int) error
	GetBudgetSummaries(ctx context.Context, userID uuid.UUID, dateStr string) ([]BudgetSummaryResponse, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) UpsertBudget(ctx context.Context, budget *Budget, categoryIDs []int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if budget.ID == uuid.Nil {
		budget.ID = uuid.New()
	}

	var exists bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM budgets WHERE id = $1)", budget.ID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		query := `
			UPDATE budgets 
			SET name = $1, amount = $2, interval_type = $3, day = $4, date = $5, month = $6, repeat = $7, updated_at = CURRENT_TIMESTAMP
			WHERE id = $8 AND user_id = $9
		`
		_, err = tx.Exec(ctx, query, budget.Name, budget.Amount, budget.IntervalType, budget.Day, budget.Date, budget.Month, budget.Repeat, budget.ID, budget.UserID)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, "DELETE FROM budget_categories WHERE budget_id = $1", budget.ID)
		if err != nil {
			return err
		}
	} else {
		query := `
			INSERT INTO budgets (id, user_id, name, amount, interval_type, day, date, month, repeat)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`
		_, err = tx.Exec(ctx, query, budget.ID, budget.UserID, budget.Name, budget.Amount, budget.IntervalType, budget.Day, budget.Date, budget.Month, budget.Repeat)
		if err != nil {
			return err
		}
	}

	for _, catID := range categoryIDs {
		_, err = tx.Exec(ctx, "INSERT INTO budget_categories (budget_id, category_id) VALUES ($1, $2)", budget.ID, catID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func snapToEndOfMonth(year int, month time.Month, date int) time.Time {
	t := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)
	lastDay := t.Day()
	if date > lastDay {
		return time.Date(year, month, lastDay, 0, 0, 0, 0, time.UTC)
	}
	return time.Date(year, month, date, 0, 0, 0, 0, time.UTC)
}

func calculatePeriod(b Budget, target time.Time) (time.Time, time.Time) {
	switch b.IntervalType {
	case "weekly":
		targetDow := int(target.Weekday())
		if targetDow == 0 {
			targetDow = 7
		}
		budgetDay := 1
		if b.Day != nil {
			budgetDay = *b.Day
		}
		diff := targetDow - budgetDay
		if diff < 0 {
			diff += 7
		}
		start := target.AddDate(0, 0, -diff)
		end := start.AddDate(0, 0, 6)
		return start, end

	case "monthly":
		year, month, day := target.Date()
		budgetDate := 1
		if b.Date != nil {
			budgetDate = *b.Date
		}
		if day < budgetDate {
			month--
			if month == 0 {
				month = 12
				year--
			}
		}
		start := snapToEndOfMonth(year, month, budgetDate)
		
		endMonth := month + 1
		endYear := year
		if endMonth == 13 {
			endMonth = 1
			endYear++
		}
		endRaw := snapToEndOfMonth(endYear, endMonth, budgetDate)
		end := endRaw.AddDate(0, 0, -1)
		return start, end

	case "yearly":
		year := target.Year()
		budgetDate := 1
		if b.Date != nil {
			budgetDate = *b.Date
		}
		budgetMonth := time.January
		if b.Month != nil {
			budgetMonth = time.Month(*b.Month)
		}
		
		if target.Month() < budgetMonth || (target.Month() == budgetMonth && target.Day() < budgetDate) {
			year--
		}
		start := snapToEndOfMonth(year, budgetMonth, budgetDate)
		endRaw := snapToEndOfMonth(year+1, budgetMonth, budgetDate)
		end := endRaw.AddDate(0, 0, -1)
		return start, end
	}
	
	return target, target
}

func (r *repository) GetBudgetSummaries(ctx context.Context, userID uuid.UUID, dateStr string) ([]BudgetSummaryResponse, error) {
	targetDate, _ := time.Parse("2006-01-02", dateStr)

	// 1. Get all budgets
	rows, err := r.db.Query(ctx, "SELECT id, user_id, name, amount, interval_type, day, date, month, repeat FROM budgets WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var budgets []Budget
	for rows.Next() {
		var b Budget
		if err := rows.Scan(&b.ID, &b.UserID, &b.Name, &b.Amount, &b.IntervalType, &b.Day, &b.Date, &b.Month, &b.Repeat); err != nil {
			return nil, err
		}
		budgets = append(budgets, b)
	}
	rows.Close()

	var summaries []BudgetSummaryResponse
	
	// 2. Process each budget
	for _, b := range budgets {
		start, end := calculatePeriod(b, targetDate)
		
		// Check if history exists
		var histID uuid.UUID
		var histAmount float64
		err := r.db.QueryRow(ctx, "SELECT id, amount FROM budget_history WHERE budget_id = $1 AND start_date = $2", b.ID, start).Scan(&histID, &histAmount)
		
		if err != nil {
			// Not found, check if we should create it
			// We create it if it's repeating, OR if it's the very first time (no history exists at all for this budget)
			var count int
			r.db.QueryRow(ctx, "SELECT COUNT(*) FROM budget_history WHERE budget_id = $1", b.ID).Scan(&count)
			
			if b.Repeat || count == 0 {
				histAmount = b.Amount
				_, err = r.db.Exec(ctx, "INSERT INTO budget_history (budget_id, amount, start_date, end_date) VALUES ($1, $2, $3, $4)", b.ID, histAmount, start, end)
				if err != nil {
					return nil, err
				}
			} else {
				// Budget doesn't repeat and already had history, so it's expired for this period.
				// We can skip returning it for the current date.
				continue
			}
		}

		// Get categories
		catRows, _ := r.db.Query(ctx, "SELECT category_id FROM budget_categories WHERE budget_id = $1", b.ID)
		var catIDs []int
		for catRows.Next() {
			var catID int
			catRows.Scan(&catID)
			catIDs = append(catIDs, catID)
		}
		catRows.Close()

		// Get total spent
		var totalSpent float64
		
		if len(catIDs) > 0 {
			r.db.QueryRow(ctx, `
				SELECT COALESCE(SUM(t.amount), 0)
				FROM transactions t
				JOIN accounts a ON t.account_id = a.id
				WHERE a.user_id = $1 AND t.transaction_type = 'credit'
				AND t.category_id = ANY($2)
				AND t.transaction_date::DATE >= $3 AND t.transaction_date::DATE <= $4
			`, userID, catIDs, start, end).Scan(&totalSpent)
		} else {
			r.db.QueryRow(ctx, `
				SELECT COALESCE(SUM(t.amount), 0)
				FROM transactions t
				JOIN accounts a ON t.account_id = a.id
				WHERE a.user_id = $1 AND t.transaction_type = 'credit'
				AND t.transaction_date::DATE >= $2 AND t.transaction_date::DATE <= $3
			`, userID, start, end).Scan(&totalSpent)
		}

		summaries = append(summaries, BudgetSummaryResponse{
			ID:                 b.ID,
			Name:               b.Name,
			Amount:             histAmount,
			IntervalType:       b.IntervalType,
			Day:                b.Day,
			Date:               b.Date,
			Month:              b.Month,
			Repeat:             b.Repeat,
			CurrentPeriodStart: start,
			CurrentPeriodEnd:   end,
			TotalSpent:         totalSpent,
			CategoryIDs:        catIDs,
		})
	}

	return summaries, nil
}
