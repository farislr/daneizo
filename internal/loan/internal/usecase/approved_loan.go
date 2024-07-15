package usecase

import "context"

type (
	ApprovedLoan interface {
		Execute(ctx context.Context, in ApprovedLoanInput) error
	}

	ApprovedLoanInput struct {
		LoanID     uint64 `json:"loan_id"     validate:"required"`
		EmployeeID uint64 `json:"employee_id" validate:"required"`
	}
)
