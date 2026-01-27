// backend/pkg/validators/validators.go
package validators

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Common validation errors
var (
	ErrInvalidUUID      = errors.New("invalid UUID format")
	ErrInvalidDecimal   = errors.New("invalid decimal format")
	ErrInvalidDate      = errors.New("invalid date format")
	ErrAmountTooLarge   = errors.New("amount exceeds maximum allowed (1 billion)")
	ErrDateTooFarFuture = errors.New("date cannot be more than 100 years in the future")
	ErrDateTooFarPast   = errors.New("date cannot be more than 100 years in the past")
)

// ValidateUUID checks if string is valid UUID
func ValidateUUID(id string) error {
	_, err := uuid.Parse(id)
	if err != nil {
		return ErrInvalidUUID
	}
	return nil
}

// ValidateDecimal validates and parses decimal string
func ValidateDecimal(value string, maxAmount decimal.Decimal) (decimal.Decimal, error) {
	dec, err := decimal.NewFromString(value)
	if err != nil {
		return decimal.Zero, ErrInvalidDecimal
	}

	if dec.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero, errors.New("amount must be greater than 0")
	}

	if dec.GreaterThan(maxAmount) {
		return decimal.Zero, ErrAmountTooLarge
	}

	return dec, nil
}

// ValidateDate validates date string and checks reasonable range
func ValidateDate(dateStr string) (time.Time, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, ErrInvalidDate
	}

	now := time.Now()

	// Check not more than 100 years in past
	hundredYearsAgo := now.AddDate(-100, 0, 0)
	if date.Before(hundredYearsAgo) {
		return time.Time{}, ErrDateTooFarPast
	}

	// Check not more than 100 years in future
	hundredYearsLater := now.AddDate(100, 0, 0)
	if date.After(hundredYearsLater) {
		return time.Time{}, ErrDateTooFarFuture
	}

	return date, nil
}

// SanitizeString removes potentially dangerous characters
func SanitizeString(input string, maxLen int) string {
	// Remove null bytes
	sanitized := ""
	for _, r := range input {
		if r != 0 {
			sanitized += string(r)
		}
	}

	// Limit length
	if len(sanitized) > maxLen {
		return sanitized[:maxLen]
	}

	return sanitized
}
