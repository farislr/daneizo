package sqlentity

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

type Loan struct {
	ID                         uint64
	BorrowerID                 uint64
	PrincipalAmount            decimal.Decimal
	InvestedAmount             decimal.Decimal
	InterestRate               decimal.Decimal
	Status                     LoanStatus
	ApprovalDate               sql.NullTime
	ApprovalEmployeeID         sql.NullInt64
	DisbursementDate           sql.NullTime
	AgreementLetterDocumentURL sql.NullString
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
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
		"agreement_letter_document_url",
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
		&l.ID,
		&l.BorrowerID,
		&l.PrincipalAmount,
		&l.InvestedAmount,
		&l.InterestRate,
		&l.Status,
		&l.ApprovalDate,
		&l.ApprovalEmployeeID,
		&l.DisbursementDate,
		&l.AgreementLetterDocumentURL,
	}
}

func (l Loan) DriverValues() []driver.Value {
	vals := make([]driver.Value, len(l.Values()))
	for i, v := range l.Values() {
		vals[i] = v
	}

	return vals
}

type Loans []Loan

func (l Loans) IsEmpty() bool {
	return l.Len() == 0
}

func (l Loans) Len() int {
	return len(l)
}

func (l Loans) First() Loan {
	if l.IsEmpty() {
		return Loan{}
	}

	return l[0]
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
	b, ok := value.([]byte)
	if ok {
		val := ls.getMap()[string(b)]

		*ls = val

		return nil
	}

	return errors.New("failed to scan loan status")
}

type ApproveLoan struct {
	ApprovalDate       sql.NullTime
	ApprovalEmployeeID sql.NullInt64
}

func (a ApproveLoan) Columns() []any {
	return []any{
		"status",
		"approval_date",
		"approval_employee_id",
	}
}

func (a ApproveLoan) StringColumns() []string {
	vals := make([]string, len(a.Columns()))
	for i, col := range a.Columns() {
		c, ok := col.(string)
		if ok {
			vals[i] = c
		}
	}

	return vals
}

func (a *ApproveLoan) Values() []any {
	return []any{
		Approved,
		a.ApprovalDate,
		a.ApprovalEmployeeID,
	}
}

func (a ApproveLoan) DriverValues() []driver.Value {
	vals := make([]driver.Value, len(a.Values()))
	for i, v := range a.Values() {
		vals[i] = v
	}

	return vals
}

func (a ApproveLoan) MappedValues() map[string]driver.Value {
	vals := make(map[string]driver.Value)
	cols := a.StringColumns()
	for i, col := range cols {
		vals[col] = a.DriverValues()[i]
	}

	return vals
}

type UpdateAmountLoan struct {
	Amount decimal.Decimal
}

func (a UpdateAmountLoan) Columns() []any {
	return []any{
		"invested_amount",
	}
}

func (a UpdateAmountLoan) StringColumns() []string {
	vals := make([]string, len(a.Columns()))
	for i, col := range a.Columns() {
		c, ok := col.(string)
		if ok {
			vals[i] = c
		}
	}

	return vals
}

func (a *UpdateAmountLoan) Values() []any {
	return []any{
		a.Amount,
	}
}

func (a UpdateAmountLoan) DriverValues() []driver.Value {
	vals := make([]driver.Value, len(a.Values()))
	for i, v := range a.Values() {
		vals[i] = v
	}

	return vals
}

func (a UpdateAmountLoan) MappedValues() map[string]driver.Value {
	vals := make(map[string]driver.Value)
	cols := a.StringColumns()
	for i, col := range cols {
		vals[col] = a.DriverValues()[i]
	}

	return vals
}

type UpdateLoanStatus struct {
	Status LoanStatus
}

func (a UpdateLoanStatus) Columns() []any {
	return []any{
		"status",
	}
}

func (a UpdateLoanStatus) StringColumns() []string {
	vals := make([]string, len(a.Columns()))
	for i, col := range a.Columns() {
		c, ok := col.(string)
		if ok {
			vals[i] = c
		}
	}

	return vals
}

func (a *UpdateLoanStatus) Values() []any {
	return []any{
		&a.Status,
	}
}

func (a UpdateLoanStatus) DriverValues() []driver.Value {
	vals := make([]driver.Value, len(a.Values()))
	for i, v := range a.Values() {
		vals[i] = v
	}

	return vals
}

func (a UpdateLoanStatus) MappedValues() map[string]driver.Value {
	vals := make(map[string]driver.Value)
	cols := a.StringColumns()
	for i, col := range cols {
		vals[col] = a.DriverValues()[i]
	}

	return vals
}

type DisburseLoan struct {
	DisburesmentDate           sql.NullTime
	AgreementLetterDocumentURL sql.NullString
}

func (a DisburseLoan) Columns() []any {
	return []any{
		"status",
		"disbursement_date",
		"agreement_letter_document_url",
	}
}

func (a DisburseLoan) StringColumns() []string {
	vals := make([]string, len(a.Columns()))
	for i, col := range a.Columns() {
		c, ok := col.(string)
		if ok {
			vals[i] = c
		}
	}

	return vals
}

func (a *DisburseLoan) Values() []any {
	return []any{
		Disbursed,
		a.DisburesmentDate,
		a.AgreementLetterDocumentURL,
	}
}

func (a DisburseLoan) DriverValues() []driver.Value {
	vals := make([]driver.Value, len(a.Values()))
	for i, v := range a.Values() {
		vals[i] = v
	}

	return vals
}

func (a DisburseLoan) MappedValues() map[string]driver.Value {
	vals := make(map[string]driver.Value)
	cols := a.StringColumns()
	for i, col := range cols {
		vals[col] = a.DriverValues()[i]
	}

	return vals
}
