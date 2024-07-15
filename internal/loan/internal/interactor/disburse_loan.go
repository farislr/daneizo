package interactor

import (
	"context"
	"database/sql"
	"time"

	"github.com/farislr/daneizo/internal/loan/internal/entity/sqlentity"
	"github.com/farislr/daneizo/internal/loan/internal/gateway"
	"github.com/farislr/daneizo/internal/loan/internal/usecase"
	"github.com/farislr/daneizo/internal/pkg/pkgerror"
	"go.uber.org/zap"
)

type (
	DisburseLoanStore interface {
		UpdateLoan(
			ctx context.Context,
			in sqlentity.UpdateEntity,
			opts ...gateway.UpdateLoanOption,
		) error
		GetLoan(ctx context.Context, opts ...gateway.GetLoanOption) (sqlentity.Loans, error)
	}

	DisburseLoan struct {
		store  DisburseLoanStore
		logger *zap.SugaredLogger
	}
)

func NewDisburseLoan(
	store DisburseLoanStore,
	logger *zap.SugaredLogger,
) *DisburseLoan {
	return &DisburseLoan{
		store:  store,
		logger: logger,
	}
}

func (d *DisburseLoan) Execute(ctx context.Context, in usecase.DisburseLoanInput) error {
	loans, err := d.store.GetLoan(ctx, gateway.GetLoanWithLoanIDFilter(in.LoanID))
	if err != nil {
		d.logger.Errorw("failed to get loan", "error", err)

		return pkgerror.ServerErrorFrom(err)
	}

	var loan sqlentity.Loan
	if loan = loans.First(); loans.IsEmpty() {
		d.logger.Errorw("loan not found")

		return pkgerror.NewBusinessError("loan not found")
	}

	if loan.Status != sqlentity.Invested {
		d.logger.Errorw("loan not invested")

		return pkgerror.NewBusinessError("loan not invested")
	}

	if err := d.store.UpdateLoan(
		ctx,
		sqlentity.DisburseLoan{
			DisburesmentDate: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
		},
		gateway.UpdateLoanWithLoanIDFilter(in.LoanID),
	); err != nil {
		d.logger.Errorw("failed to update loan", "error", err)

		return err
	}

	return nil
}
