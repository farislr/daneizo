package repositories

import (
	"context"
	"github.com/farislr/daneizo/internal/loanv2/domain/valueobjects"
)

// User represents a user entity (simplified for this domain context)
type User struct {
	ID                 valueobjects.UserID
	Email              valueobjects.EmailAddress
	PhoneNumber        valueobjects.PhoneNumber
	CreditScore        valueobjects.CreditScore
	MonthlyIncome      valueobjects.Money
	MonthlyDebtPayment valueobjects.Money
	IsActive           bool
	IsVerified         bool
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	// FindByID retrieves a user by their ID
	FindByID(ctx context.Context, id valueobjects.UserID) (*User, error)
	
	// FindByEmail retrieves a user by their email address
	FindByEmail(ctx context.Context, email valueobjects.EmailAddress) (*User, error)
	
	// Save persists user information
	Save(ctx context.Context, user *User) error
	
	// UpdateCreditScore updates a user's credit score
	UpdateCreditScore(ctx context.Context, userID valueobjects.UserID, score valueobjects.CreditScore) error
	
	// UpdateIncome updates a user's monthly income
	UpdateIncome(ctx context.Context, userID valueobjects.UserID, income valueobjects.Money) error
	
	// UpdateDebtPayment updates a user's monthly debt payment
	UpdateDebtPayment(ctx context.Context, userID valueobjects.UserID, debt valueobjects.Money) error
	
	// Exists checks if a user exists
	Exists(ctx context.Context, id valueobjects.UserID) (bool, error)
	
	// IsEligibleForLoan checks if a user meets basic loan eligibility criteria
	IsEligibleForLoan(ctx context.Context, userID valueobjects.UserID, requestedAmount valueobjects.Money) (bool, error)
}