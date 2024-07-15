package gateway

import (
	"context"
	"net/http"
	"strconv"

	"github.com/farislr/daneizo/internal/loan/internal/usecase"
	"github.com/farislr/daneizo/internal/pkg/pkgerror"
	"github.com/farislr/daneizo/internal/pkg/pkghttp/v1"
	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

func NewLoanHTTPGateway(
	httpRouter *httprouter.Router,
	logger *zap.SugaredLogger,
	loanHTTPEndpoint *LoanHTTPEndpoint,
	validator *validator.Validate,
) {
	server := pkghttp.NewServer(
		pkghttp.WithResponseEncoder(pkghttp.CodeMessageResponseEncoder),
		pkghttp.WithErrorResponseEncoder(pkghttp.CodeMessageErrorEncoder),
	)

	httpRouter.Handler(
		http.MethodPost,
		"/loan/:loan_id/approve",
		server.Serve(loanHTTPEndpoint.ApproveLoan),
	)

	httpRouter.Handler(http.MethodPost, "/loan", server.Serve(loanHTTPEndpoint.CreateNewLoan))

	httpRouter.Handler(
		http.MethodPost,
		"/loan/:loan_id/invest",
		server.Serve(loanHTTPEndpoint.InvestLoan),
	)

	httpRouter.Handler(
		http.MethodPost,
		"/loan/:loan_id/disburse",
		server.Serve(loanHTTPEndpoint.DisburseLoan),
	)
}

type LoanHTTPEndpoint struct {
	createProposedLoanUsecase usecase.CreateProposedLoan
	approveLoanUsecase        usecase.ApprovedLoan
	investLoanUsecase         usecase.InvestLoan
	disburseLoanUsecase       usecase.DisburseLoan

	validator *validator.Validate
	logger    *zap.SugaredLogger
}

func NewLoanHTTPEndpoint(
	createNewLoanUsecase usecase.CreateProposedLoan,
	approveLoanUsecase usecase.ApprovedLoan,
	investLoanUsecase usecase.InvestLoan,
	disburseLoanUsecase usecase.DisburseLoan,

	logger *zap.SugaredLogger,
	validator *validator.Validate,

) *LoanHTTPEndpoint {
	return &LoanHTTPEndpoint{
		createProposedLoanUsecase: createNewLoanUsecase,
		approveLoanUsecase:        approveLoanUsecase,
		investLoanUsecase:         investLoanUsecase,
		disburseLoanUsecase:       disburseLoanUsecase,

		logger:    logger,
		validator: validator,
	}
}

func (l *LoanHTTPEndpoint) CreateNewLoan(
	ctx context.Context,
	request pkghttp.Request,
) (any, error) {
	var input usecase.CreateProposedLoanInput
	if err := request.Decode(&input); err != nil {
		l.logger.Errorw("failed to decode request", "error", err)

		return nil, pkgerror.ServerErrorFrom(err)
	}

	if err := l.validator.Struct(input); err != nil {
		l.logger.Errorw("failed to validate request", "error", err)

		return nil, pkgerror.ValidationErrorFrom(err)
	}

	if err := l.createProposedLoanUsecase.Execute(ctx, input); err != nil {
		l.logger.Errorw("failed to create new loan", "error", err)

		return nil, err
	}

	return nil, nil
}

func (l *LoanHTTPEndpoint) ApproveLoan(
	ctx context.Context,
	request pkghttp.Request,
) (resp any, err error) {
	var input usecase.ApprovedLoanInput
	if err := request.Decode(&input); err != nil {
		l.logger.Errorw("failed to decode request", "error", err)

		return nil, pkgerror.ServerErrorFrom(err)
	}

	params := httprouter.ParamsFromContext(ctx)

	loanID := params.ByName("loan_id")

	input.LoanID, err = strconv.ParseUint(loanID, 10, 64)
	if err != nil {
		l.logger.Errorw("failed to parse loan id", "error", err)

		return nil, pkgerror.ValidationErrorFrom(err)
	}

	if err := l.validator.Struct(input); err != nil {
		l.logger.Errorw("failed to validate request", "error", err)

		return nil, pkgerror.ValidationErrorFrom(err)
	}

	if err := l.approveLoanUsecase.Execute(ctx, input); err != nil {
		l.logger.Errorw("failed to approve loan", "error", err)

		return nil, err
	}

	return nil, nil
}

func (l *LoanHTTPEndpoint) InvestLoan(
	ctx context.Context,
	request pkghttp.Request,
) (resp any, err error) {
	var input usecase.InvestLoanInput
	if err := request.Decode(&input); err != nil {
		l.logger.Errorw("failed to decode request", "error", err)

		return nil, pkgerror.ServerErrorFrom(err)
	}

	if err := l.validator.Struct(input); err != nil {
		l.logger.Errorw("failed to validate request", "error", err)

		return nil, pkgerror.ValidationErrorFrom(err)
	}

	params := httprouter.ParamsFromContext(ctx)

	loanID := params.ByName("loan_id")

	input.LoanID, err = strconv.ParseUint(loanID, 10, 64)
	if err != nil {
		l.logger.Errorw("failed to parse loan id", "error", err)

		return nil, pkgerror.ValidationErrorFrom(err)
	}

	if err := l.investLoanUsecase.Execute(ctx, input); err != nil {
		l.logger.Errorw("failed to invest loan", "error", err)

		return nil, err
	}

	return nil, nil
}

func (l *LoanHTTPEndpoint) DisburseLoan(
	ctx context.Context,
	request pkghttp.Request,
) (resp any, err error) {
	var input usecase.DisburseLoanInput
	if err := request.Decode(&input); err != nil {
		l.logger.Errorw("failed to decode request", "error", err)

		return nil, pkgerror.ServerErrorFrom(err)
	}

	if err := l.validator.Struct(input); err != nil {
		l.logger.Errorw("failed to validate request", "error", err)

		return nil, pkgerror.ValidationErrorFrom(err)
	}

	params := httprouter.ParamsFromContext(ctx)

	loanID := params.ByName("loan_id")

	input.LoanID, err = strconv.ParseUint(loanID, 10, 64)
	if err != nil {
		l.logger.Errorw("failed to parse loan id", "error", err)

		return nil, pkgerror.ValidationErrorFrom(err)
	}

	return nil, nil
}
