package usecase

import "context"

type (
	DisburseLoan interface {
		Execute(ctx context.Context, in DisburseLoanInput) error
	}

	DisburseLoanInput struct {
		LoanID uint64 `json:"loan_id"`
	}
)
