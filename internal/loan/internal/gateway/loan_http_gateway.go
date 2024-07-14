package gateway

import (
	"context"
	"net/http"

	"github.com/farislr/daneizo/internal/loan/internal/usecase"
	"github.com/farislr/daneizo/internal/pkg/pkghttp/v1"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

func NewLoanHTTPGateway(
	httpRouter *httprouter.Router,
	logger *zap.SugaredLogger,
	loanHTTPEndpoint *LoanHTTPEndpoint,
) {
	server := pkghttp.NewServer()

	httpRouter.Handler(http.MethodPost, "/loan", server.Serve(loanHTTPEndpoint.CreateNewLoan))
}

type LoanHTTPEndpoint struct {
	createProposedLoanUsecase usecase.CreateProposedLoan

	logger *zap.SugaredLogger
}

func NewLoanHTTPEndpoint(
	logger *zap.SugaredLogger,
	createNewLoanUsecase usecase.CreateProposedLoan,
) *LoanHTTPEndpoint {
	return &LoanHTTPEndpoint{
		createProposedLoanUsecase: createNewLoanUsecase,
		logger:                    logger,
	}
}

func (l *LoanHTTPEndpoint) CreateNewLoan(
	ctx context.Context,
	request pkghttp.Request,
) (any, error) {
	var input usecase.CreateProposedLoanInput
	if err := request.Decode(&input); err != nil {
		l.logger.Errorw("failed to decode request", "error", err)

		return nil, err
	}

	if err := l.createProposedLoanUsecase.Execute(ctx, input); err != nil {
		l.logger.Errorw("failed to create new loan", "error", err)

		return nil, err
	}

	return nil, nil
}
