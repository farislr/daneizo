package usecase

import (
	"context"

	"github.com/shopspring/decimal"
)

type (
	CreateProposedLoan interface {
		Execute(ctx context.Context, in CreateProposedLoanInput) error
	}

	CreateProposedLoanInput struct {
		UserID       uint64          `json:"user_id"       validate:"required"` // UserID should get from authorization
		InterestRate decimal.Decimal `json:"interest_rate" validate:"required"`
		Amount       decimal.Decimal `json:"amount"        validate:"required"`
	}
)
