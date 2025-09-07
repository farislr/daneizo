package repositories

import (
	"context"
	"github.com/farislr/daneizo/internal/loanv2/domain/entities"
	"github.com/farislr/daneizo/internal/loanv2/domain/valueobjects"
)

// LoanRepository defines the interface for loan aggregate persistence
type LoanRepository interface {
	// Save persists a loan aggregate (insert or update)
	Save(ctx context.Context, loan *entities.Loan) error
	
	// FindByID retrieves a loan by its ID
	FindByID(ctx context.Context, id valueobjects.LoanID) (*entities.Loan, error)
	
	// FindByBorrower retrieves all loans for a specific borrower
	FindByBorrower(ctx context.Context, borrowerID valueobjects.UserID) ([]*entities.Loan, error)
	
	// FindByStatus retrieves loans with a specific status
	FindByStatus(ctx context.Context, status valueobjects.LoanStatus) ([]*entities.Loan, error)
	
	// FindByStatusWithPagination retrieves loans with pagination
	FindByStatusWithPagination(ctx context.Context, status valueobjects.LoanStatus, offset, limit int) ([]*entities.Loan, error)
	
	// FindRequiringApproval retrieves loans that require approval
	FindRequiringApproval(ctx context.Context) ([]*entities.Loan, error)
	
	// FindAvailableForInvestment retrieves loans available for investment
	FindAvailableForInvestment(ctx context.Context) ([]*entities.Loan, error)
	
	// FindFullyFundedLoans retrieves loans that are fully funded and ready for disbursement
	FindFullyFundedLoans(ctx context.Context) ([]*entities.Loan, error)
	
	// FindActiveLoans retrieves loans that are currently active (accepting payments)
	FindActiveLoans(ctx context.Context) ([]*entities.Loan, error)
	
	// FindOverdueLoans retrieves loans with overdue payments
	FindOverdueLoans(ctx context.Context) ([]*entities.Loan, error)
	
	// CountByStatus returns the count of loans by status
	CountByStatus(ctx context.Context, status valueobjects.LoanStatus) (int64, error)
	
	// CountByBorrower returns the count of loans for a borrower
	CountByBorrower(ctx context.Context, borrowerID valueobjects.UserID) (int64, error)
	
	// Delete removes a loan (should rarely be used - prefer status changes)
	Delete(ctx context.Context, id valueobjects.LoanID) error
	
	// Exists checks if a loan exists
	Exists(ctx context.Context, id valueobjects.LoanID) (bool, error)
}