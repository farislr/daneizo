package gateway

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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
	opts := []pkghttp.EndpointOption{
		pkghttp.WithEndpointResponseEncoder(pkghttp.CodeMessageResponseEncoder),
		pkghttp.WithEndpointErrorResponseEncoder(pkghttp.CodeMessageErrorEncoder),
	}

	httpRouter.Handler(
		http.MethodPost,
		"/loan/:loan_id/approve",
		pkghttp.NewHandler(loanHTTPEndpoint.ApproveLoan, opts...),
	)

	httpRouter.Handler(http.MethodPost, "/loan", pkghttp.NewHandler(loanHTTPEndpoint.CreateNewLoan, opts...))

	httpRouter.Handler(
		http.MethodPost,
		"/loan/:loan_id/invest",
		pkghttp.NewHandler(loanHTTPEndpoint.InvestLoan, opts...),
	)

	httpRouter.Handler(
		http.MethodPost,
		"/loan/:loan_id/disburse",
		pkghttp.NewHandler(loanHTTPEndpoint.DisburseLoan, opts...),
	)

	httpRouter.Handler(
		http.MethodPost,
		"/loan/:loan_id/upload-agreement-letter",
		pkghttp.NewHandler(loanHTTPEndpoint.UploadAgreementLetter, opts...),
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
	request *pkghttp.Request[usecase.CreateProposedLoanInput],
) (any, error) {
	input := request.Body()

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
	request *pkghttp.Request[usecase.ApprovedLoanInput],
) (resp any, err error) {
	input := request.Body()
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
	request *pkghttp.Request[usecase.InvestLoanInput],
) (resp any, err error) {
	input := request.Body()

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
	request *pkghttp.Request[usecase.DisburseLoanInput],
) (resp any, err error) {
	input := request.Body()

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

	if err := l.disburseLoanUsecase.Execute(ctx, input); err != nil {
		l.logger.Errorw("failed to disburse loan", "error", err)

		return nil, err
	}

	return nil, nil
}

func (l *LoanHTTPEndpoint) UploadAgreementLetter(
	ctx context.Context,
	request *pkghttp.Request[struct{}],
) (resp any, err error) {
	rawRequest := request.Raw()
	params := httprouter.ParamsFromContext(ctx)

	// 5MB max file size
	if err := rawRequest.ParseMultipartForm(5 << 20); err != nil {
		l.logger.Errorw("failed to parse multipart form", "error", err)

		return nil, pkgerror.ServerErrorFrom(err)
	}

	multipartFile, multipartHeader, err := rawRequest.FormFile("agreement_letter")
	if err != nil {
		l.logger.Errorw("failed to get multipart file", "error", err)

		return nil, pkgerror.ServerErrorFrom(err)
	}
	defer multipartFile.Close()

	workDir, err := os.Getwd()
	if err != nil {
		l.logger.Errorw("failed to get current working directory", "error", err)

		return nil, pkgerror.ServerErrorFrom(err)
	}

	loanID := params.ByName("loan_id")
	uploadPath := fmt.Sprintf("%s/%s/%s", workDir, "files/agreement-letter/", loanID)

	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		l.logger.Errorw("failed to create upload directory", "error", err)

		return nil, pkgerror.ServerErrorFrom(err)

	}

	file, err := os.Create(
		fmt.Sprintf(
			"%s/%s%s",
			uploadPath,
			l.fileNameWithoutExtension(multipartHeader.Filename),
			filepath.Ext(multipartHeader.Filename),
		),
	)
	if err != nil {
		l.logger.Errorw("failed to create temp file", "error", err)

		return nil, pkgerror.ServerErrorFrom(err)
	}
	defer file.Close()

	if _, err := file.ReadFrom(multipartFile); err != nil {
		l.logger.Errorw("failed to read file", "error", err)

		return nil, pkgerror.ServerErrorFrom(err)
	}

	return nil, nil
}

func (l *LoanHTTPEndpoint) fileNameWithoutExtension(filename string) string {
	return filename[:len(filename)-len(filepath.Ext(filename))]
}
