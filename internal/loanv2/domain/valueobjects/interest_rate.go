package valueobjects

import (
	"fmt"
	"github.com/shopspring/decimal"
)

// InterestRate represents an interest rate as a percentage with business validation and calculations
type InterestRate struct {
	rate decimal.Decimal // as percentage (0-100)
}

// NewInterestRate creates a new InterestRate with validation
func NewInterestRate(rate decimal.Decimal) (InterestRate, error) {
	if rate.IsNegative() {
		return InterestRate{}, fmt.Errorf("interest rate cannot be negative: %s", rate.String())
	}
	
	if rate.GreaterThan(decimal.NewFromFloat(50.0)) {
		return InterestRate{}, fmt.Errorf("interest rate cannot exceed 50%%: %s", rate.String())
	}
	
	return InterestRate{rate: rate}, nil
}

// NewInterestRateFromFloat creates an InterestRate from a float64
func NewInterestRateFromFloat(rate float64) (InterestRate, error) {
	return NewInterestRate(decimal.NewFromFloat(rate))
}

// CalculateMonthlyInterest calculates the monthly interest for a principal amount over a term
func (ir InterestRate) CalculateMonthlyInterest(principal Money, termMonths int) Money {
	if termMonths <= 0 {
		return NewZeroMoney(principal.Currency())
	}
	
	monthlyRate := ir.rate.Div(decimal.NewFromInt(100)).Div(decimal.NewFromInt(12))
	interest := principal.amount.Mul(monthlyRate).Mul(decimal.NewFromInt(int64(termMonths)))
	
	result, _ := NewMoney(interest, principal.currency)
	return result
}

// CalculateTotalInterest calculates the total interest for a principal amount over a term
func (ir InterestRate) CalculateTotalInterest(principal Money, termMonths int) Money {
	if termMonths <= 0 {
		return NewZeroMoney(principal.Currency())
	}
	
	annualRate := ir.rate.Div(decimal.NewFromInt(100))
	yearsDecimal := decimal.NewFromInt(int64(termMonths)).Div(decimal.NewFromInt(12))
	totalInterest := principal.amount.Mul(annualRate).Mul(yearsDecimal)
	
	result, _ := NewMoney(totalInterest, principal.currency)
	return result
}

// CalculateMonthlyPayment calculates the monthly payment using the standard loan formula
// Monthly Payment = P * [r(1+r)^n] / [(1+r)^n - 1]
func (ir InterestRate) CalculateMonthlyPayment(principal Money, termMonths int) Money {
	if termMonths <= 0 {
		return NewZeroMoney(principal.Currency())
	}
	
	// If interest rate is zero, just divide principal by term
	if ir.rate.IsZero() {
		monthlyPrincipal := principal.amount.Div(decimal.NewFromInt(int64(termMonths)))
		result, _ := NewMoney(monthlyPrincipal, principal.currency)
		return result
	}
	
	monthlyRate := ir.rate.Div(decimal.NewFromInt(100)).Div(decimal.NewFromInt(12))
	termDecimal := decimal.NewFromInt(int64(termMonths))
	
	// Calculate (1 + r)^n
	onePlusR := decimal.NewFromInt(1).Add(monthlyRate)
	onePlusRPowN := onePlusR.Pow(termDecimal)
	
	// Calculate numerator: r * (1+r)^n
	numerator := monthlyRate.Mul(onePlusRPowN)
	
	// Calculate denominator: (1+r)^n - 1
	denominator := onePlusRPowN.Sub(decimal.NewFromInt(1))
	
	// Calculate monthly payment: P * [numerator / denominator]
	monthlyPayment := principal.amount.Mul(numerator.Div(denominator))
	
	result, _ := NewMoney(monthlyPayment, principal.currency)
	return result
}

// AsMonthlyRate returns the interest rate as a monthly decimal rate
func (ir InterestRate) AsMonthlyRate() decimal.Decimal {
	return ir.rate.Div(decimal.NewFromInt(100)).Div(decimal.NewFromInt(12))
}

// AsAnnualRate returns the interest rate as an annual decimal rate
func (ir InterestRate) AsAnnualRate() decimal.Decimal {
	return ir.rate.Div(decimal.NewFromInt(100))
}

// Value returns the raw interest rate as a percentage
func (ir InterestRate) Value() decimal.Decimal {
	return ir.rate
}

// String returns a formatted string representation
func (ir InterestRate) String() string {
	return fmt.Sprintf("%s%%", ir.rate.StringFixed(2))
}

// Add adds two interest rates
func (ir InterestRate) Add(other InterestRate) (InterestRate, error) {
	newRate := ir.rate.Add(other.rate)
	return NewInterestRate(newRate)
}

// Subtract subtracts one interest rate from another
func (ir InterestRate) Subtract(other InterestRate) (InterestRate, error) {
	newRate := ir.rate.Sub(other.rate)
	return NewInterestRate(newRate)
}

// MultiplyBy multiplies the interest rate by a factor
func (ir InterestRate) MultiplyBy(factor decimal.Decimal) (InterestRate, error) {
	newRate := ir.rate.Mul(factor)
	return NewInterestRate(newRate)
}

// IsZero returns true if the interest rate is zero
func (ir InterestRate) IsZero() bool {
	return ir.rate.IsZero()
}

// GreaterThan checks if this rate is greater than another
func (ir InterestRate) GreaterThan(other InterestRate) bool {
	return ir.rate.GreaterThan(other.rate)
}

// LessThan checks if this rate is less than another
func (ir InterestRate) LessThan(other InterestRate) bool {
	return ir.rate.LessThan(other.rate)
}

// Equal checks if this rate equals another
func (ir InterestRate) Equal(other InterestRate) bool {
	return ir.rate.Equal(other.rate)
}