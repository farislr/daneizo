package interactor

import (
	"context"
	"errors"
	"testing"

	"github.com/farislr/daneizo/internal/loan/internal/entity/sqlentity"
	loanmocks "github.com/farislr/daneizo/internal/loan/internal/mocks"
	"github.com/farislr/daneizo/internal/loan/internal/usecase"
	"github.com/farislr/daneizo/internal/pkg/pkgmocks"
	"github.com/farislr/daneizo/internal/pkg/pkguid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func TestCreateProposedLoan_Execute(t *testing.T) {
	store := loanmocks.NewMockInsertLoanStore(t)
	logger := zap.NewNop().Sugar()
	snowflakeGen := pkgmocks.NewMockSnowflake(t)

	type fields struct {
		store        InsertLoanStore
		logger       *zap.SugaredLogger
		snowflakeGen pkguid.Snowflake
	}
	type args struct {
		ctx context.Context
		in  usecase.CreateProposedLoanInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mockFn  func(a args)
		wantErr bool
	}{
		{
			name: "error when insert loan",
			fields: fields{
				store:        store,
				logger:       logger,
				snowflakeGen: snowflakeGen,
			},
			args: args{
				ctx: context.Background(),
				in: usecase.CreateProposedLoanInput{
					UserID: 1,
					Amount: decimal.NewFromInt(1_000_000),
				},
			},
			mockFn: func(a args) {
				loanID := 1

				snowflakeGen.EXPECT().Generate().Return(uint64(loanID)).Once()

				store.EXPECT().InsertLoan(a.ctx, sqlentity.Loan{
					ID:              uint64(loanID),
					BorrowerID:      a.in.UserID,
					PrincipalAmount: a.in.Amount,
					InvestedAmount:  decimal.Zero,
					Status:          sqlentity.Proposed,
				}).Return(errors.New("any error")).Once()
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				store:        store,
				logger:       logger,
				snowflakeGen: snowflakeGen,
			},
			args: args{
				ctx: context.Background(),
				in: usecase.CreateProposedLoanInput{
					UserID: 1,
					Amount: decimal.NewFromInt(1_000_000),
				},
			},
			mockFn: func(a args) {
				loanID := 1

				snowflakeGen.EXPECT().Generate().Return(uint64(loanID)).Once()

				store.EXPECT().InsertLoan(a.ctx, sqlentity.Loan{
					ID:              uint64(loanID),
					BorrowerID:      a.in.UserID,
					PrincipalAmount: a.in.Amount,
					InvestedAmount:  decimal.Zero,
					Status:          sqlentity.Proposed,
				}).Return(nil).Once()
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)

			c := NewCreateProposedLoan(
				tt.fields.store,
				tt.fields.logger,
				tt.fields.snowflakeGen,
			)
			if err := c.Execute(tt.args.ctx, tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("CreateProposedLoan.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
