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
	UpdateLoanStore interface {
		UpdateLoan(
			ctx context.Context,
			in sqlentity.UpdateEntity,
			opts ...gateway.UpdateLoanOption,
		) error
	}

	ApproveLoan struct {
		store UpdateLoanStore

		logger *zap.SugaredLogger
	}
)

func NewApproveLoan(
	store UpdateLoanStore,
	logger *zap.SugaredLogger,
) *ApproveLoan {
	return &ApproveLoan{
		store:  store,
		logger: logger,
	}
}

func (a *ApproveLoan) Execute(
	ctx context.Context,
	in usecase.ApprovedLoanInput,
) error {
	if err := a.store.UpdateLoan(ctx, sqlentity.UpdateApproveLoan{
		ApprovalDate: sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		},
		ApprovalEmployeeID: sql.NullInt64{
			Valid: true,
			Int64: int64(in.EmployeeID),
		},
	}, gateway.UpdateLoanWithLoanIDFilter(in.LoanID)); err != nil {
		a.logger.Errorw("failed to update loan", "error", err)

		return pkgerror.ServerErrorFrom(err)
	}

	return nil
}
