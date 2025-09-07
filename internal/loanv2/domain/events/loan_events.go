package events

import (
	"time"
	"github.com/farislr/daneizo/internal/loanv2/domain/valueobjects"
)

// LoanProposedEvent is emitted when a new loan is proposed
type LoanProposedEvent struct {
	BaseEvent
	LoanID          valueobjects.LoanID
	BorrowerID      valueobjects.UserID
	PrincipalAmount valueobjects.Money
	InterestRate    valueobjects.InterestRate
	TermMonths      int
	LoanPurpose     string
}

// NewLoanProposedEvent creates a new loan proposed event
func NewLoanProposedEvent(
	loanID valueobjects.LoanID,
	borrowerID valueobjects.UserID,
	principalAmount valueobjects.Money,
	interestRate valueobjects.InterestRate,
	termMonths int,
	purpose string,
	correlationID string,
) LoanProposedEvent {
	return LoanProposedEvent{
		BaseEvent: NewBaseEvent(
			loanID.String(),
			"Loan",
			&borrowerID,
			correlationID,
			nil,
		),
		LoanID:          loanID,
		BorrowerID:      borrowerID,
		PrincipalAmount: principalAmount,
		InterestRate:    interestRate,
		TermMonths:      termMonths,
		LoanPurpose:     purpose,
	}
}

// EventName returns the event name
func (e LoanProposedEvent) EventName() string {
	return "loan.proposed"
}

// Payload returns the event payload
func (e LoanProposedEvent) Payload() interface{} {
	return e
}

// LoanApprovedEvent is emitted when a loan is approved
type LoanApprovedEvent struct {
	BaseEvent
	LoanID         valueobjects.LoanID
	ApproverID     valueobjects.UserID
	ApprovalDate   time.Time
	ApprovedAmount valueobjects.Money
	ApprovedRate   valueobjects.InterestRate
	Conditions     []string
}

// NewLoanApprovedEvent creates a new loan approved event
func NewLoanApprovedEvent(
	loanID valueobjects.LoanID,
	approverID valueobjects.UserID,
	approvedAmount valueobjects.Money,
	approvedRate valueobjects.InterestRate,
	conditions []string,
	correlationID string,
) LoanApprovedEvent {
	return LoanApprovedEvent{
		BaseEvent: NewBaseEvent(
			loanID.String(),
			"Loan",
			&approverID,
			correlationID,
			nil,
		),
		LoanID:         loanID,
		ApproverID:     approverID,
		ApprovalDate:   time.Now(),
		ApprovedAmount: approvedAmount,
		ApprovedRate:   approvedRate,
		Conditions:     conditions,
	}
}

// EventName returns the event name
func (e LoanApprovedEvent) EventName() string {
	return "loan.approved"
}

// Payload returns the event payload
func (e LoanApprovedEvent) Payload() interface{} {
	return e
}

// LoanRejectedEvent is emitted when a loan is rejected
type LoanRejectedEvent struct {
	BaseEvent
	LoanID    valueobjects.LoanID
	ReviewerID valueobjects.UserID
	Reason    string
	Details   []string
}

// NewLoanRejectedEvent creates a new loan rejected event
func NewLoanRejectedEvent(
	loanID valueobjects.LoanID,
	reviewerID valueobjects.UserID,
	reason string,
	details []string,
	correlationID string,
) LoanRejectedEvent {
	return LoanRejectedEvent{
		BaseEvent: NewBaseEvent(
			loanID.String(),
			"Loan",
			&reviewerID,
			correlationID,
			nil,
		),
		LoanID:     loanID,
		ReviewerID: reviewerID,
		Reason:     reason,
		Details:    details,
	}
}

// EventName returns the event name
func (e LoanRejectedEvent) EventName() string {
	return "loan.rejected"
}

// Payload returns the event payload
func (e LoanRejectedEvent) Payload() interface{} {
	return e
}

// LoanFullyFundedEvent is emitted when a loan becomes fully funded
type LoanFullyFundedEvent struct {
	BaseEvent
	LoanID              valueobjects.LoanID
	TotalAmount         valueobjects.Money
	NumberOfInvestors   int
	FundingCompletedAt  time.Time
	DisbursementEligible bool
}

// NewLoanFullyFundedEvent creates a new loan fully funded event
func NewLoanFullyFundedEvent(
	loanID valueobjects.LoanID,
	totalAmount valueobjects.Money,
	numInvestors int,
	disbursementEligible bool,
	correlationID string,
) LoanFullyFundedEvent {
	return LoanFullyFundedEvent{
		BaseEvent: NewBaseEvent(
			loanID.String(),
			"Loan",
			nil,
			correlationID,
			nil,
		),
		LoanID:               loanID,
		TotalAmount:          totalAmount,
		NumberOfInvestors:    numInvestors,
		FundingCompletedAt:   time.Now(),
		DisbursementEligible: disbursementEligible,
	}
}

// EventName returns the event name
func (e LoanFullyFundedEvent) EventName() string {
	return "loan.fully_funded"
}

// Payload returns the event payload
func (e LoanFullyFundedEvent) Payload() interface{} {
	return e
}

// LoanDisbursedEvent is emitted when loan funds are disbursed
type LoanDisbursedEvent struct {
	BaseEvent
	LoanID              valueobjects.LoanID
	BorrowerID          valueobjects.UserID
	DisbursedAmount     valueobjects.Money
	DisbursementDate    time.Time
	AgreementDocumentURL string
	FirstPaymentDue     time.Time
}

// NewLoanDisbursedEvent creates a new loan disbursed event
func NewLoanDisbursedEvent(
	loanID valueobjects.LoanID,
	borrowerID valueobjects.UserID,
	amount valueobjects.Money,
	agreementURL string,
	firstPaymentDue time.Time,
	correlationID string,
) LoanDisbursedEvent {
	return LoanDisbursedEvent{
		BaseEvent: NewBaseEvent(
			loanID.String(),
			"Loan",
			&borrowerID,
			correlationID,
			nil,
		),
		LoanID:               loanID,
		BorrowerID:           borrowerID,
		DisbursedAmount:      amount,
		DisbursementDate:     time.Now(),
		AgreementDocumentURL: agreementURL,
		FirstPaymentDue:      firstPaymentDue,
	}
}

// EventName returns the event name
func (e LoanDisbursedEvent) EventName() string {
	return "loan.disbursed"
}

// Payload returns the event payload
func (e LoanDisbursedEvent) Payload() interface{} {
	return e
}

// LoanStatusChangedEvent is emitted when a loan status changes
type LoanStatusChangedEvent struct {
	BaseEvent
	LoanID    valueobjects.LoanID
	FromStatus valueobjects.LoanStatus
	ToStatus   valueobjects.LoanStatus
	Reason     string
}

// NewLoanStatusChangedEvent creates a new loan status changed event
func NewLoanStatusChangedEvent(
	loanID valueobjects.LoanID,
	fromStatus, toStatus valueobjects.LoanStatus,
	reason string,
	userID *valueobjects.UserID,
	correlationID string,
) LoanStatusChangedEvent {
	return LoanStatusChangedEvent{
		BaseEvent: NewBaseEvent(
			loanID.String(),
			"Loan",
			userID,
			correlationID,
			nil,
		),
		LoanID:     loanID,
		FromStatus: fromStatus,
		ToStatus:   toStatus,
		Reason:     reason,
	}
}

// EventName returns the event name
func (e LoanStatusChangedEvent) EventName() string {
	return "loan.status_changed"
}

// Payload returns the event payload
func (e LoanStatusChangedEvent) Payload() interface{} {
	return e
}