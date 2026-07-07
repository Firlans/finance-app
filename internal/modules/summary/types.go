package summary

import "time"

type TransactionCategoryBalance struct {
	CategoryName string  `json:"category_name"`
	Balance      float64 `json:"balance"`
}

type TransactionCategoryBalanceListResponse struct {
	Message   string                       `json:"message"`
	Data      []TransactionCategoryBalance `json:"data"`
	RequestID string                       `json:"request_id"`
}

type Module string

const (
	ModuleDebit  Module = "debit"  // pemasukan
	ModuleCredit Module = "credit" // pengeluaran
)

type TransactionSummaryQuery struct {
	Module Module
	From   time.Time
	To     time.Time
}
