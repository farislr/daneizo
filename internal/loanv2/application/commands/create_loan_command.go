package commands

import (
	"context"
	"fmt"
	"github.com/farislr/daneizo/internal/loanv2/domain/valueobjects"
	"github.com/farislr/daneizo/internal/loanv2/domain/entities"
	"github.com/farislr/daneizo/internal/loanv2/domain/repositories"
	"github.com/farislr/daneizo/internal/loanv2/domain/services"
	"github.com/farislr/daneizo/internal/loanv2/domain/events"
)

// CreateLoanCommand represents the command to create a new loan
type CreateLoanCommand struct {
	LoanID      valueobjects.LoanID
	BorrowerID  valueobjects.UserID
	Amount      valueobjects.Money
	TermMonths  int
	Purpose     string
	UserID      *valueobjects.UserID // User making the request (for auditing)
	CorrelationID string
}

// CreateLoanCommandHandler handles loan creation commands
type CreateLoanCommandHandler struct {
	loanRepository         repositories.LoanRepository
	userRepository         repositories.UserRepository
	riskAssessmentService  *services.RiskAssessmentService
	domainEventPublisher   DomainEventPublisher
	notificationService    NotificationService
}

// DomainEventPublisher interface for publishing domain events
type DomainEventPublisher interface {
	Publish(ctx context.Context, events []events.DomainEvent) error
	PublishAsync(ctx context.Context, events []events.DomainEvent) error
}

// NotificationService interface for sending notifications
type NotificationService interface {
	SendEmail(ctx context.Context, to valueobjects.EmailAddress, subject string, body string) error
	SendSMS(ctx context.Context, to valueobjects.PhoneNumber, message string) error
}

// NewCreateLoanCommandHandler creates a new command handler
func NewCreateLoanCommandHandler(
	loanRepository repositories.LoanRepository,
	userRepository repositories.UserRepository,
	riskAssessmentService *services.RiskAssessmentService,
	domainEventPublisher DomainEventPublisher,
	notificationService NotificationService,
) *CreateLoanCommandHandler {
	return &CreateLoanCommandHandler{
		loanRepository:        loanRepository,
		userRepository:        userRepository,
		riskAssessmentService: riskAssessmentService,
		domainEventPublisher:  domainEventPublisher,
		notificationService:   notificationService,
	}
}

// Handle processes the create loan command
func (h *CreateLoanCommandHandler) Handle(ctx context.Context, cmd CreateLoanCommand) error {
	// 1. Validate borrower exists and is eligible
	borrower, err := h.userRepository.FindByID(ctx, cmd.BorrowerID)
	if err != nil {
		return fmt.Errorf("borrower not found: %w", err)
	}

	if !borrower.IsActive || !borrower.IsVerified {
		return fmt.Errorf("borrower is not active or verified")
	}

	// 2. Check if loan ID is already taken
	exists, err := h.loanRepository.Exists(ctx, cmd.LoanID)
	if err != nil {
		return fmt.Errorf("failed to check loan existence: %w", err)
	}
	if exists {
		return fmt.Errorf("loan with ID %s already exists", cmd.LoanID.String())
	}

	// 3. Perform initial risk assessment to get appropriate interest rate
	riskResult, err := h.riskAssessmentService.AssessLoanRisk(
		ctx, cmd.BorrowerID, cmd.Amount, cmd.TermMonths,
	)
	if err != nil {
		return fmt.Errorf("risk assessment failed: %w", err)
	}

	// 4. Check if loan is approvable based on risk
	if riskResult.OverallRisk == valueobjects.UnacceptableRisk {
		return fmt.Errorf("loan application rejected due to unacceptable risk: %v", riskResult.Reasoning)
	}

	// 5. Create loan aggregate using domain logic
	loan, err := entities.NewLoanProposal(
		cmd.LoanID, cmd.BorrowerID, cmd.Amount, 
		riskResult.RecommendedRate, cmd.TermMonths, cmd.Purpose,
	)
	if err != nil {
		return fmt.Errorf("failed to create loan proposal: %w", err)
	}

	// 6. Save aggregate (this will also handle the business invariants)
	if err := h.loanRepository.Save(ctx, loan); err != nil {
		return fmt.Errorf("failed to save loan: %w", err)
	}

	// 7. Publish domain events
	if len(loan.DomainEvents()) > 0 {
		if err := h.domainEventPublisher.Publish(ctx, loan.DomainEvents()); err != nil {
			// Don't fail the operation, but log the error
			// In production, you might want to use a saga or retry mechanism
			fmt.Printf("Failed to publish domain events: %v\n", err)
		}
		loan.ClearDomainEvents()
	}

	// 8. Send notification to borrower
	if err := h.notificationService.SendEmail(
		ctx,
		borrower.Email,
		"Loan Application Received",
		fmt.Sprintf(
			"Your loan application for %s has been received and is under review. "+
				"Application ID: %s. You will receive updates as your application progresses.",
			cmd.Amount.String(),
			cmd.LoanID.String(),
		),
	); err != nil {
		// Don't fail the operation for notification failures
		fmt.Printf("Failed to send notification: %v\n", err)
	}

	return nil
}

// Validate validates the command
func (cmd CreateLoanCommand) Validate() error {
	if cmd.LoanID.IsZero() {
		return fmt.Errorf("loan ID is required")
	}

	if cmd.BorrowerID.IsZero() {
		return fmt.Errorf("borrower ID is required")
	}

	if cmd.Amount.IsZero() {
		return fmt.Errorf("loan amount must be greater than zero")
	}

	if cmd.TermMonths < 6 || cmd.TermMonths > 60 {
		return fmt.Errorf("loan term must be between 6 and 60 months")
	}

	if cmd.Purpose == "" {
		return fmt.Errorf("loan purpose is required")
	}

	return nil
}