package valueobjects

import (
	"fmt"
	"github.com/shopspring/decimal"
)

// Money represents a monetary amount with currency awareness
type Money struct {
	amount   decimal.Decimal
	currency Currency
}

// NewMoney creates a new Money value object with validation
func NewMoney(amount decimal.Decimal, currency Currency) (Money, error) {
	if amount.IsNegative() {
		return Money{}, fmt.Errorf("money amount cannot be negative: %s", amount.String())
	}
	
	if !currency.IsValid() {
		return Money{}, fmt.Errorf("invalid currency: %s", currency)
	}
	
	return Money{
		amount:   amount,
		currency: currency,
	}, nil
}

// NewZeroMoney creates a zero Money value for the given currency
func NewZeroMoney(currency Currency) Money {
	return Money{
		amount:   decimal.Zero,
		currency: currency,
	}
}

// Add adds two Money amounts, ensuring same currency
func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("cannot add different currencies: %s and %s", m.currency, other.currency)
	}
	
	return Money{
		amount:   m.amount.Add(other.amount),
		currency: m.currency,
	}, nil
}

// Subtract subtracts two Money amounts, ensuring same currency
func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("cannot subtract different currencies: %s and %s", m.currency, other.currency)
	}
	
	result := m.amount.Sub(other.amount)
	if result.IsNegative() {
		return Money{}, fmt.Errorf("result would be negative: %s", result.String())
	}
	
	return Money{
		amount:   result,
		currency: m.currency,
	}, nil
}

// MultiplyByRate multiplies the money amount by a rate
func (m Money) MultiplyByRate(rate decimal.Decimal) Money {
	return Money{
		amount:   m.amount.Mul(rate),
		currency: m.currency,
	}
}

// MultiplyByPercent multiplies the money amount by a percentage (0-100)
func (m Money) MultiplyByPercent(percent decimal.Decimal) Money {
	rate := percent.Div(decimal.NewFromInt(100))
	return m.MultiplyByRate(rate)
}

// DivideBy divides the money amount by a divisor
func (m Money) DivideBy(divisor decimal.Decimal) (Money, error) {
	if divisor.IsZero() {
		return Money{}, fmt.Errorf("cannot divide by zero")
	}
	
	return Money{
		amount:   m.amount.Div(divisor),
		currency: m.currency,
	}, nil
}

// IsZero returns true if the amount is zero
func (m Money) IsZero() bool {
	return m.amount.IsZero()
}

// IsPositive returns true if the amount is positive
func (m Money) IsPositive() bool {
	return m.amount.IsPositive()
}

// GreaterThan checks if this money is greater than another
func (m Money) GreaterThan(other Money) (bool, error) {
	if m.currency != other.currency {
		return false, fmt.Errorf("cannot compare different currencies: %s and %s", m.currency, other.currency)
	}
	return m.amount.GreaterThan(other.amount), nil
}

// GreaterThanOrEqual checks if this money is greater than or equal to another
func (m Money) GreaterThanOrEqual(other Money) (bool, error) {
	if m.currency != other.currency {
		return false, fmt.Errorf("cannot compare different currencies: %s and %s", m.currency, other.currency)
	}
	return m.amount.GreaterThanOrEqual(other.amount), nil
}

// LessThan checks if this money is less than another
func (m Money) LessThan(other Money) (bool, error) {
	if m.currency != other.currency {
		return false, fmt.Errorf("cannot compare different currencies: %s and %s", m.currency, other.currency)
	}
	return m.amount.LessThan(other.amount), nil
}

// LessThanOrEqual checks if this money is less than or equal to another
func (m Money) LessThanOrEqual(other Money) (bool, error) {
	if m.currency != other.currency {
		return false, fmt.Errorf("cannot compare different currencies: %s and %s", m.currency, other.currency)
	}
	return m.amount.LessThanOrEqual(other.amount), nil
}

// Equal checks if this money equals another
func (m Money) Equal(other Money) bool {
	return m.currency == other.currency && m.amount.Equal(other.amount)
}

// Amount returns the decimal amount
func (m Money) Amount() decimal.Decimal {
	return m.amount
}

// Currency returns the currency
func (m Money) Currency() Currency {
	return m.currency
}

// String returns a formatted string representation
func (m Money) String() string {
	return fmt.Sprintf("%s %s", m.currency.Symbol(), m.amount.StringFixed(2))
}

// StringWithCode returns a string representation with currency code
func (m Money) StringWithCode() string {
	return fmt.Sprintf("%s %s", m.amount.StringFixed(2), m.currency.String())
}