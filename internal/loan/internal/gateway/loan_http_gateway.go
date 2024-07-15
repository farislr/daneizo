package gateway

import (
	"context"
	"net/http"

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

	httpRouter.Handler(http.MethodPost, "/loan", server.Serve(loanHTTPEndpoint.CreateNewLoan))
}

type LoanHTTPEndpoint struct {
	createProposedLoanUsecase usecase.CreateProposedLoan

	validator *validator.Validate
	logger    *zap.SugaredLogger
}

func NewLoanHTTPEndpoint(
	logger *zap.SugaredLogger,
	createNewLoanUsecase usecase.CreateProposedLoan,
	validator *validator.Validate,

) *LoanHTTPEndpoint {
	return &LoanHTTPEndpoint{
		createProposedLoanUsecase: createNewLoanUsecase,
		logger:                    logger,
		validator:                 validator,
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
