package commands

import (
	"context"
	"fmt"
	"github.com/farislr/daneizo/internal/loanv2/domain/valueobjects"
	"github.com/farislr/daneizo/internal/loanv2/domain/repositories"
)

// ApproveLoanCommand represents the command to approve a loan
type ApproveLoanCommand struct {
	LoanID        valueobjects.LoanID
	ApproverID    valueobjects.UserID
	Conditions    []string
	Comments      string
	CorrelationID string
}

// ApproveLoanCommandHandler handles loan approval commands
type ApproveLoanCommandHandler struct {
	loanRepository       repositories.LoanRepository
	userRepository       repositories.UserRepository
	domainEventPublisher DomainEventPublisher
	notificationService  NotificationService
}

// NewApproveLoanCommandHandler creates a new command handler
func NewApproveLoanCommandHandler(
	loanRepository repositories.LoanRepository,
	userRepository repositories.UserRepository,
	domainEventPublisher DomainEventPublisher,
	notificationService NotificationService,
) *ApproveLoanCommandHandler {
	return &ApproveLoanCommandHandler{
		loanRepository:       loanRepository,
		userRepository:       userRepository,
		domainEventPublisher: domainEventPublisher,
		notificationService:  notificationService,
	}
}

// Handle processes the approve loan command
func (h *ApproveLoanCommandHandler) Handle(ctx context.Context, cmd ApproveLoanCommand) error {
	// 1. Validate the command
	if err := cmd.Validate(); err != nil {
		return fmt.Errorf("invalid command: %w", err)
	}

	// 2. Load loan aggregate
	loan, err := h.loanRepository.FindByID(ctx, cmd.LoanID)
	if err != nil {
		return fmt.Errorf("loan not found: %w", err)
	}

	// 3. Verify approver exists and has permission
	approver, err := h.userRepository.FindByID(ctx, cmd.ApproverID)
	if err != nil {
		return fmt.Errorf("approver not found: %w", err)
	}

	if !approver.IsActive {
		return fmt.Errorf("approver is not active")
	}

	// 4. Apply business logic through aggregate
	if err := loan.Approve(cmd.ApproverID, cmd.Conditions); err != nil {
		return fmt.Errorf("failed to approve loan: %w", err)
	}

	// 5. Save the updated aggregate
	if err := h.loanRepository.Save(ctx, loan); err != nil {
		return fmt.Errorf("failed to save approved loan: %w", err)
	}

	// 6. Publish domain events
	if len(loan.DomainEvents()) > 0 {
		if err := h.domainEventPublisher.Publish(ctx, loan.DomainEvents()); err != nil {
			fmt.Printf("Failed to publish domain events: %v\n", err)
		}
		loan.ClearDomainEvents()
	}

	// 7. Send notification to borrower
	borrower, err := h.userRepository.FindByID(ctx, loan.BorrowerID())
	if err == nil { // Don't fail if borrower lookup fails
		notificationErr := h.notificationService.SendEmail(
			ctx,
			borrower.Email,
			"Loan Application Approved",
			fmt.Sprintf(
				"Great news! Your loan application for %s has been approved. "+
					"Your loan is now available for investment. "+
					"Loan ID: %s",
				loan.PrincipalAmount().String(),
				loan.ID().String(),
			),
		)
		if notificationErr != nil {
			fmt.Printf("Failed to send approval notification: %v\n", notificationErr)
		}
	}

	return nil
}

// Validate validates the command
func (cmd ApproveLoanCommand) Validate() error {
	if cmd.LoanID.IsZero() {
		return fmt.Errorf("loan ID is required")
	}

	if cmd.ApproverID.IsZero() {
		return fmt.Errorf("approver ID is required")
	}

	return nil
}