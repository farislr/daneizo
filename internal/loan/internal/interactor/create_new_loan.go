package interactor

import (
	"context"

	"github.com/farislr/daneizo/internal/loan/internal/entity/sqlentity"
	"github.com/farislr/daneizo/internal/loan/internal/usecase"
	"github.com/farislr/daneizo/internal/pkg/pkguid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type (
	InsertLoanStore interface {
		InsertLoan(ctx context.Context, in sqlentity.Loan) error
	}

	CreateProposedLoan struct {
		store        InsertLoanStore
		logger       *zap.SugaredLogger
		snowflakeGen pkguid.Snowflake
	}
)

func NewCreateProposedLoan(
	store InsertLoanStore,
	logger *zap.SugaredLogger,
	snowflakeGen pkguid.Snowflake,
) *CreateProposedLoan {
	return &CreateProposedLoan{
		store:        store,
		logger:       logger,
		snowflakeGen: snowflakeGen,
	}
}

func (c *CreateProposedLoan) Execute(
	ctx context.Context,
	in usecase.CreateProposedLoanInput,
) error {

	if err := c.store.InsertLoan(ctx, sqlentity.Loan{
		ID:              c.snowflakeGen.Generate(),
		BorrowerID:      in.UserID,
		PrincipalAmount: in.Amount,
		InterestRate:    in.InterestRate,
		InvestedAmount:  decimal.Zero,
		Status:          sqlentity.Proposed,
	}); err != nil {
		c.logger.Errorw("failed to insert loan", "error", err)

		return err
	}

	return nil
}
