package sqlentity

import (
	"database/sql/driver"

	"github.com/shopspring/decimal"
)

type LoanInvestment struct {
	ID         uint64
	LoanID     uint64
	InvestorID uint64
	Amount     decimal.Decimal
}

func (l LoanInvestment) Columns() []any {
	return []any{
		"id",
		"loan_id",
		"investor_id",
		"amount",
	}
}

func (l LoanInvestment) StringColumns() []string {
	vals := make([]string, len(l.Columns()))
	for i, col := range l.Columns() {
		c, ok := col.(string)
		if ok {
			vals[i] = c
		}
	}

	return vals
}

func (l *LoanInvestment) Values() []any {
	return []any{
		l.ID,
		l.LoanID,
		l.InvestorID,
		l.Amount,
	}
}

func (l *LoanInvestment) DriverValues() []driver.Value {
	vals := make([]driver.Value, len(l.Values()))
	for i, v := range l.Values() {
		vals[i] = v
	}

	return vals
}
