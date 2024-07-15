package loan

import (
	"database/sql"

	"github.com/farislr/daneizo/internal/loan/internal/gateway"
	"github.com/farislr/daneizo/internal/loan/internal/interactor"
	"github.com/farislr/daneizo/internal/pkg/pkgsql"
	"github.com/farislr/daneizo/internal/pkg/pkguid"
	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type Exposed struct {
}

type Dependencies struct {
	DB           *sql.DB
	Logger       *zap.SugaredLogger
	QueryBuilder pkgsql.GoquBuilder
	SnowflakeGen pkguid.Snowflake
	HttpRouter   *httprouter.Router
	Validator    *validator.Validate
}

func New(deps Dependencies) *Exposed {

	loanSQLstore := gateway.NewLoanSQLGateway(deps.DB, deps.Logger, deps.QueryBuilder)

	createProposedLoanUsecase := interactor.NewCreateProposedLoan(
		loanSQLstore,
		deps.Logger,
		deps.SnowflakeGen,
	)

	loanHTTPEndpoint := gateway.NewLoanHTTPEndpoint(
		deps.Logger,
		createProposedLoanUsecase,
		deps.Validator,
	)

	gateway.NewLoanHTTPGateway(deps.HttpRouter, deps.Logger, loanHTTPEndpoint, deps.Validator)

	return &Exposed{}
}
