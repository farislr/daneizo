package gateway

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/farislr/daneizo/internal/loan/internal/entity/sqlentity"
	"github.com/farislr/daneizo/internal/pkg/pkgsql"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type loanSQLGatewaySuite struct {
	db           *sql.DB
	dbmock       sqlmock.Sqlmock
	queryBuilder pkgsql.GoquBuilder

	loanTableName           string
	loanInvestmentTableName string
	userTableName           string

	suite.Suite
}

func TestLoanSQLGatewaySuite(t *testing.T) {
	suite.Run(t, new(loanSQLGatewaySuite))
}

func (ls *loanSQLGatewaySuite) SetupSuite() {
	var err error
	ls.db, ls.dbmock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		ls.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	ls.queryBuilder = goqu.New("mysql", ls.db)
	ls.loanTableName = "loans"
	ls.loanInvestmentTableName = "loan_investments"
	ls.userTableName = "users"
}

func (ls *loanSQLGatewaySuite) TestLoanSQLGateway_InsertLoan() {

	type args struct {
		ctx context.Context
		in  sqlentity.Loan
	}

	tests := []struct {
		name    string
		args    args
		mockFn  func(a args)
		wantErr bool
	}{
		{
			name: "error RowsAffected",
			args: args{
				ctx: context.Background(),
				in: sqlentity.Loan{
					ID:                 1,
					BorrowerID:         1,
					PrincipalAmount:    decimal.Decimal{},
					InvestedAmount:     decimal.Decimal{},
					InterestRate:       decimal.Decimal{},
					Status:             sqlentity.Proposed,
					ApprovalDate:       sql.NullTime{},
					ApprovalEmployeeID: sql.NullInt64{},
					DisbursementDate:   sql.NullTime{},
					CreatedAt:          time.Time{},
					UpdatedAt:          time.Time{},
				},
			},
			mockFn: func(a args) {
				query, _, err := ls.queryBuilder.Insert(ls.loanTableName).
					Cols(a.in.Columns()...).
					Vals(a.in.Values()).
					ToSQL()
				ls.NoError(err)

				ls.dbmock.ExpectExec(query).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("error")))
			},
			wantErr: true,
		},
		{
			name: "error exec",
			args: args{
				ctx: context.Background(),
				in: sqlentity.Loan{
					ID:                 1,
					BorrowerID:         1,
					PrincipalAmount:    decimal.Decimal{},
					InvestedAmount:     decimal.Decimal{},
					InterestRate:       decimal.Decimal{},
					Status:             sqlentity.Proposed,
					ApprovalDate:       sql.NullTime{},
					ApprovalEmployeeID: sql.NullInt64{},
					DisbursementDate:   sql.NullTime{},
					CreatedAt:          time.Time{},
					UpdatedAt:          time.Time{},
				},
			},
			mockFn: func(a args) {
				query, _, err := ls.queryBuilder.Insert(ls.loanTableName).
					Cols(a.in.Columns()...).
					Vals(a.in.Values()).
					ToSQL()
				ls.NoError(err)

				ls.dbmock.ExpectExec(query).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				in: sqlentity.Loan{
					ID:                 1,
					BorrowerID:         1,
					PrincipalAmount:    decimal.Decimal{},
					InvestedAmount:     decimal.Decimal{},
					InterestRate:       decimal.Decimal{},
					Status:             sqlentity.Proposed,
					ApprovalDate:       sql.NullTime{},
					ApprovalEmployeeID: sql.NullInt64{},
					DisbursementDate:   sql.NullTime{},
					CreatedAt:          time.Time{},
					UpdatedAt:          time.Time{},
				},
			},
			mockFn: func(a args) {
				query, _, err := ls.queryBuilder.Insert(ls.loanTableName).
					Cols(a.in.Columns()...).
					Vals(a.in.Values()).
					ToSQL()
				ls.NoError(err)

				ls.dbmock.ExpectExec(query).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		ls.Run(tt.name, func() {
			tt.mockFn(tt.args)

			r := NewLoanSQLGateway(
				ls.db,
				zap.NewNop().Sugar(),
				ls.queryBuilder,
			)
			if err := r.InsertLoan(tt.args.ctx, tt.args.in); (err != nil) != tt.wantErr {
				ls.T().Errorf("LoanSQLGateway.InsertLoan() error = %v, wantErr %v", err, tt.wantErr)
			}
			ls.NoError(ls.dbmock.ExpectationsWereMet())
		})
	}
}

func (ls *loanSQLGatewaySuite) TestLoanSQLGateway_UpdateLoan() {

	type args struct {
		ctx  context.Context
		in   sqlentity.Entity
		opts []UpdateLoanOption
	}
	tests := []struct {
		name    string
		args    args
		mockFn  func(a args)
		wantErr bool
	}{
		{
			name: "error RowsAffected",
			args: args{
				ctx: context.Background(),
				in: &sqlentity.UpdateApproveLoan{
					ApprovalDate: sql.NullTime{
						Valid: true,
						Time:  time.Now(),
					},
					ApprovalEmployeeID: sql.NullInt64{
						Valid: true,
						Int64: 1,
					},
				},
				opts: []UpdateLoanOption{
					UpdateLoanWithLoanIDFilter(1),
				},
			},
			mockFn: func(a args) {
				query, _, err := ls.queryBuilder.Update(ls.loanTableName).
					Set(a.in.MappedValues()).
					Where(goqu.Ex{"id": 1}).
					ToSQL()
				ls.NoError(err)

				ls.dbmock.ExpectExec(query).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("error")))
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				in: &sqlentity.UpdateApproveLoan{
					ApprovalDate: sql.NullTime{
						Valid: true,
						Time:  time.Now(),
					},
					ApprovalEmployeeID: sql.NullInt64{
						Valid: true,
						Int64: 1,
					},
				},
				opts: []UpdateLoanOption{
					UpdateLoanWithLoanIDFilter(1),
				},
			},
			mockFn: func(a args) {
				query, _, err := ls.queryBuilder.Update(ls.loanTableName).
					Set(a.in.MappedValues()).
					Where(goqu.Ex{"id": 1}).
					ToSQL()
				ls.NoError(err)

				ls.dbmock.ExpectExec(query).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		ls.Run(tt.name, func() {
			tt.mockFn(tt.args)

			r := NewLoanSQLGateway(
				ls.db,
				zap.NewNop().Sugar(),
				ls.queryBuilder,
			)
			if err := r.UpdateLoan(tt.args.ctx, tt.args.in, tt.args.opts...); (err != nil) != tt.wantErr {
				ls.T().Logf("LoanSQLGateway.UpdateLoan() error = %v, wantErr %v", err, tt.wantErr)

				ls.T().Errorf("LoanSQLGateway.UpdateLoan() error = %v, wantErr %v", err, tt.wantErr)
			}
			ls.NoError(ls.dbmock.ExpectationsWereMet())
		})
	}
}
