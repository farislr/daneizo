package sqlentity

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/shopspring/decimal"
)

type Loan struct {
	ID                 uint64
	BorrowerID         uint64
	PrincipalAmount    decimal.Decimal
	InvestedAmount     decimal.Decimal
	InterestRate       decimal.Decimal
	Status             LoanStatus
	ApprovalDate       sql.NullTime
	ApprovalEmployeeID sql.NullInt64
	DisbursementDate   sql.NullTime
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (l Loan) Columns() []any {
	return []any{
		"id",
		"borrower_id",
		"principal_amount",
		"invested_amount",
		"interest_rate",
		"status",
		"approval_date",
		"approval_employee_id",
		"disbursement_date",
		"created_at",
		"updated_at",
	}
}

func (l Loan) StringColumns() []string {
	vals := make([]string, len(l.Columns()))
	for i, col := range l.Columns() {
		c, ok := col.(string)
		if ok {
			vals[i] = c
		}
	}

	return vals
}

func (l *Loan) Values() []any {
	return []any{
		l.ID,
		l.BorrowerID,
		l.PrincipalAmount,
		l.InvestedAmount,
		l.InterestRate,
		l.Status,
		l.ApprovalDate,
		l.ApprovalEmployeeID,
		l.DisbursementDate,
		l.CreatedAt,
		l.UpdatedAt,
	}
}

type Loans []Loan

func (l Loan) DriverValues() []driver.Value {
	vals := make([]driver.Value, len(l.Values()))
	for i, v := range l.Values() {
		vals[i] = v
	}

	return vals
}

type LoanStatus int

const (
	UnknownStatus LoanStatus = iota
	Proposed
	Approved
	Invested
	Disbursed
)

func (ls LoanStatus) String() string {
	return [...]string{"UNKNOWN", "PROPOSED", "APPROVED", "INVESTED", "DISBURSED"}[ls]
}

func (ls LoanStatus) Value() (driver.Value, error) {
	return ls.String(), nil
}

func (ls LoanStatus) getMap() map[string]LoanStatus {
	return map[string]LoanStatus{
		"UNKNOWN":   UnknownStatus,
		"PROPOSED":  Proposed,
		"APPROVED":  Approved,
		"INVESTED":  Invested,
		"DISBURSED": Disbursed,
	}
}

func (ls *LoanStatus) Scan(value any) error {
	s, ok := value.(string)
	if ok {
		val := ls.getMap()[s]

		*ls = val
	}

	return nil
}
