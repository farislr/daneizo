package usecase

import (
	"context"

	"github.com/shopspring/decimal"
)

type (
	InvestLoan interface {
		Execute(ctx context.Context, in InvestLoanInput) error
	}

	InvestLoanInput struct {
		LoanID     uint64          `json:"loan_id"`
		InvestorID uint64          `json:"investor_id" validate:"required"`
		Amount     decimal.Decimal `json:"amount"      validate:"required"`
	}
)
