package gateway

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/farislr/daneizo/internal/loan/internal/entity/sqlentity"
	"github.com/farislr/daneizo/internal/pkg/pkgsql"

	"go.uber.org/zap"
)

type LoanSQLGateway struct {
	db           pkgsql.SQL
	logger       *zap.SugaredLogger
	queryBuilder pkgsql.GoquBuilder

	loanTableName           string
	loanInvestmentTableName string
	userTableName           string
}

func NewLoanSQLGateway(
	db *sql.DB,
	logger *zap.SugaredLogger,
	queryBuilder pkgsql.GoquBuilder,
) *LoanSQLGateway {
	return &LoanSQLGateway{
		db:           db,
		logger:       logger,
		queryBuilder: queryBuilder,

		loanTableName:           "loans",
		loanInvestmentTableName: "loan_investments",
		userTableName:           "users",
	}
}

func (r *LoanSQLGateway) InsertLoan(ctx context.Context, in sqlentity.Loan) error {
	query := r.queryBuilder.Insert(r.loanTableName).Cols(in.Columns()...).Vals(in.Values())
	sql, _, err := query.ToSQL()
	if err != nil {
		r.logger.Errorw("failed to build query", "error", err)

		return err
	}

	res, err := r.db.ExecContext(ctx, sql)
	if err != nil {
		r.logger.Errorw("failed to execute query", "error", err)

		return err
	}

	row, err := res.RowsAffected()
	if err != nil {
		r.logger.Errorw("failed to get last insert id", "error", err)

		return err
	}

	if row == 0 {
		return fmt.Errorf("failed to insert loan")
	}

	return nil
}

type GetLoanOption func(*goqu.SelectDataset) *goqu.SelectDataset

func WithBorrowerIDFilter(borrowerID uint64) GetLoanOption {
	return func(query *goqu.SelectDataset) *goqu.SelectDataset {
		return query.Where(goqu.Ex{"borrower_id": borrowerID})
	}
}

func WithStatusFilter(status sqlentity.LoanStatus) GetLoanOption {
	return func(query *goqu.SelectDataset) *goqu.SelectDataset {
		return query.Where(goqu.Ex{"status": status})
	}
}

func (r *LoanSQLGateway) GetLoan(
	ctx context.Context,
	opts ...GetLoanOption,
) (sqlentity.Loans, error) {
	var loan sqlentity.Loan
	query := r.queryBuilder.Select(loan.Columns()...).From(r.loanTableName)

	for _, opt := range opts {
		query = opt(query)
	}

	sql, _, err := query.ToSQL()
	if err != nil {
		r.logger.Errorw("failed to build query", "error", err)

		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, sql)
	if err != nil {
		r.logger.Errorw("failed to execute query", "error", err)

		return nil, err
	}

	var loans sqlentity.Loans
	for rows.Next() {
		err := rows.Scan(loan.Values()...)
		if err != nil {
			r.logger.Errorw("failed to scan row", "error", err)

			return nil, err
		}

		loans = append(loans, loan)
	}

	return loans, nil
}

type UpdateLoanOption func(*goqu.UpdateDataset) *goqu.UpdateDataset

func WithLoanIDFilter(loanID uint64) UpdateLoanOption {
	return func(query *goqu.UpdateDataset) *goqu.UpdateDataset {
		return query.Where(goqu.Ex{"id": loanID})
	}
}

func (r *LoanSQLGateway) UpdateLoan(
	ctx context.Context,
	in sqlentity.UpdateEntity,
	opts ...UpdateLoanOption,
) error {
	query := r.queryBuilder.Update(r.loanTableName).Set(in.MappedValues())

	for _, opt := range opts {
		query = opt(query)
	}

	sql, _, err := query.ToSQL()
	if err != nil {
		r.logger.Errorw("failed to build query", "error", err)

		return err
	}

	res, err := r.db.ExecContext(ctx, sql)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		r.logger.Errorw("failed to execute query", "error", err)

		return err
	}

	row, err := res.RowsAffected()
	if err != nil {
		r.logger.Errorw("failed to get last insert id", "error", err)

		return err
	}

	if row == 0 {
		return fmt.Errorf("loan not found")
	}

	return nil
}
