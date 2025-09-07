package loanv2

import (
	"context"
	"database/sql"
	
	"github.com/farislr/daneizo/internal/loanv2/application/commands"
	"github.com/farislr/daneizo/internal/loanv2/domain/repositories"
	sqlrepo "github.com/farislr/daneizo/internal/loanv2/infrastructure/persistence/sql"
	
	"github.com/farislr/daneizo/internal/pkg/pkguid"
	"go.uber.org/zap"
)

// Module represents the loan domain module with all its dependencies
type Module struct {
	// Repositories
	loanRepository repositories.LoanRepository
	
	// Infrastructure (stubbed for now)
	eventPublisher      commands.DomainEventPublisher
	notificationService commands.NotificationService
}

// ModuleConfig contains configuration for the loan module
type ModuleConfig struct {
	Database     *sql.DB
	Logger       *zap.Logger
	SnowflakeGen pkguid.Snowflake
}

// NewModule creates and configures a new loan domain module
func NewModule(config ModuleConfig) (*Module, error) {
	// Repositories
	loanRepo := sqlrepo.NewSQLLoanRepository(config.Database, nil, config.Logger)
	
	return &Module{
		loanRepository:      loanRepo,
		eventPublisher:      nil, // TODO: implement event publisher
		notificationService: nil, // TODO: implement notification service
	}, nil
}

// GetLoanRepository returns the loan repository
func (m *Module) GetLoanRepository() repositories.LoanRepository {
	return m.loanRepository
}

// Close gracefully closes all resources
func (m *Module) Close(ctx context.Context) error {
	// Close any resources that need cleanup
	return nil
}

// Health checks the health of the module
func (m *Module) Health(ctx context.Context) error {
	// For now, just return nil as we don't have a fully functional repository yet
	// In the future, this would check database connectivity
	return nil
}