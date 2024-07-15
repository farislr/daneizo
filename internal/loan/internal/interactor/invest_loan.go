package interactor

import (
	"context"

	"github.com/farislr/daneizo/internal/loan/internal/entity/sqlentity"
	"github.com/farislr/daneizo/internal/loan/internal/gateway"
	"github.com/farislr/daneizo/internal/loan/internal/usecase"
	"github.com/farislr/daneizo/internal/pkg/pkgerror"
	"github.com/farislr/daneizo/internal/pkg/pkguid"
	"go.uber.org/zap"
)

type (
	InvestLoanStore interface {
		InsertLoanInvestment(ctx context.Context, in sqlentity.LoanInvestment) error
		UpdateLoan(
			ctx context.Context,
			in sqlentity.UpdateEntity,
			opts ...gateway.UpdateLoanOption,
		) error

		GetLoan(ctx context.Context, opts ...gateway.GetLoanOption) (sqlentity.Loans, error)
	}

	InvestLoan struct {
		store        InvestLoanStore
		logger       *zap.SugaredLogger
		snowflakeGen pkguid.Snowflake
	}
)

func NewInvestLoan(
	store InvestLoanStore,
	logger *zap.SugaredLogger,
	snowflakeGen pkguid.Snowflake,
) *InvestLoan {
	return &InvestLoan{
		store:        store,
		logger:       logger,
		snowflakeGen: snowflakeGen,
	}
}

func (i *InvestLoan) Execute(ctx context.Context, in usecase.InvestLoanInput) error {
	loans, err := i.store.GetLoan(ctx, gateway.GetLoanWithLoanIDFilter(in.LoanID))
	if err != nil {
		i.logger.Errorw("failed to get loan", "error", err)

		return err
	}

	var loan sqlentity.Loan
	if loan = loans.First(); loans.IsEmpty() {
		i.logger.Errorw("loan not found")

		return pkgerror.NewBusinessError("loan not found")
	}

	if loan.Status == sqlentity.Invested {
		i.logger.Errorw("loan already invested")

		return pkgerror.NewBusinessError("loan already invested")
	}

	loanInvestmentID := i.snowflakeGen.Generate()

	if err := i.store.InsertLoanInvestment(
		ctx,
		sqlentity.LoanInvestment{
			ID:         loanInvestmentID,
			LoanID:     in.LoanID,
			InvestorID: in.InvestorID,
			Amount:     in.Amount,
		},
	); err != nil {
		i.logger.Errorw("failed to insert loan investment", "error", err)

		return err
	}

	totalInvestedAmount := loan.InvestedAmount.Add(in.Amount)

	if err := i.store.UpdateLoan(
		ctx,
		sqlentity.UpdateAmountLoan{
			Amount: totalInvestedAmount,
		},
		gateway.UpdateLoanWithLoanIDFilter(in.LoanID),
	); err != nil {
		i.logger.Errorw("failed to update loan", "error", err)

		return err
	}

	if totalInvestedAmount.GreaterThanOrEqual(loan.PrincipalAmount) {
		if err := i.store.UpdateLoan(
			ctx,
			sqlentity.UpdateLoanStatus{
				Status: sqlentity.Invested,
			},
			gateway.UpdateLoanWithLoanIDFilter(in.LoanID),
		); err != nil {
			i.logger.Errorw("failed to update loan", "error", err)

			return err
		}
	}

	return nil
}
