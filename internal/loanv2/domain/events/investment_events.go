package events

import (
	"time"
	"github.com/farislr/daneizo/internal/loanv2/domain/valueobjects"
)

// InvestmentMadeEvent is emitted when an investment is made in a loan
type InvestmentMadeEvent struct {
	BaseEvent
	LoanID            valueobjects.LoanID
	InvestorID        valueobjects.UserID
	InvestmentID      string
	Amount            valueobjects.Money
	InvestmentDate    time.Time
	ExpectedReturn    valueobjects.Money
	RemainingAmount   valueobjects.Money
	IsFullyFunded     bool
}

// NewInvestmentMadeEvent creates a new investment made event
func NewInvestmentMadeEvent(
	loanID valueobjects.LoanID,
	investorID valueobjects.UserID,
	investmentID string,
	amount valueobjects.Money,
	expectedReturn valueobjects.Money,
	remainingAmount valueobjects.Money,
	isFullyFunded bool,
	correlationID string,
) InvestmentMadeEvent {
	return InvestmentMadeEvent{
		BaseEvent: NewBaseEvent(
			loanID.String(),
			"Loan",
			&investorID,
			correlationID,
			nil,
		),
		LoanID:            loanID,
		InvestorID:        investorID,
		InvestmentID:      investmentID,
		Amount:            amount,
		InvestmentDate:    time.Now(),
		ExpectedReturn:    expectedReturn,
		RemainingAmount:   remainingAmount,
		IsFullyFunded:     isFullyFunded,
	}
}

// EventName returns the event name
func (e InvestmentMadeEvent) EventName() string {
	return "investment.made"
}

// Payload returns the event payload
func (e InvestmentMadeEvent) Payload() interface{} {
	return e
}

// InvestmentWithdrawnEvent is emitted when an investment is withdrawn (if allowed)
type InvestmentWithdrawnEvent struct {
	BaseEvent
	LoanID            valueobjects.LoanID
	InvestorID        valueobjects.UserID
	InvestmentID      string
	WithdrawnAmount   valueobjects.Money
	WithdrawalReason  string
	WithdrawalFee     valueobjects.Money
	RemainingInvested valueobjects.Money
}

// NewInvestmentWithdrawnEvent creates a new investment withdrawn event
func NewInvestmentWithdrawnEvent(
	loanID valueobjects.LoanID,
	investorID valueobjects.UserID,
	investmentID string,
	withdrawnAmount valueobjects.Money,
	reason string,
	fee valueobjects.Money,
	remainingInvested valueobjects.Money,
	correlationID string,
) InvestmentWithdrawnEvent {
	return InvestmentWithdrawnEvent{
		BaseEvent: NewBaseEvent(
			loanID.String(),
			"Investment",
			&investorID,
			correlationID,
			nil,
		),
		LoanID:            loanID,
		InvestorID:        investorID,
		InvestmentID:      investmentID,
		WithdrawnAmount:   withdrawnAmount,
		WithdrawalReason:  reason,
		WithdrawalFee:     fee,
		RemainingInvested: remainingInvested,
	}
}

// EventName returns the event name
func (e InvestmentWithdrawnEvent) EventName() string {
	return "investment.withdrawn"
}

// Payload returns the event payload
func (e InvestmentWithdrawnEvent) Payload() interface{} {
	return e
}

// InvestmentReturnPaidEvent is emitted when returns are paid to investors
type InvestmentReturnPaidEvent struct {
	BaseEvent
	LoanID        valueobjects.LoanID
	InvestorID    valueobjects.UserID
	InvestmentID  string
	PaidAmount    valueobjects.Money
	PrincipalPart valueobjects.Money
	InterestPart  valueobjects.Money
	PaymentDate   time.Time
	IsPartial     bool
	IsFinal       bool
}

// NewInvestmentReturnPaidEvent creates a new investment return paid event
func NewInvestmentReturnPaidEvent(
	loanID valueobjects.LoanID,
	investorID valueobjects.UserID,
	investmentID string,
	paidAmount valueobjects.Money,
	principalPart valueobjects.Money,
	interestPart valueobjects.Money,
	isPartial bool,
	isFinal bool,
	correlationID string,
) InvestmentReturnPaidEvent {
	return InvestmentReturnPaidEvent{
		BaseEvent: NewBaseEvent(
			investmentID,
			"Investment",
			nil,
			correlationID,
			nil,
		),
		LoanID:        loanID,
		InvestorID:    investorID,
		InvestmentID:  investmentID,
		PaidAmount:    paidAmount,
		PrincipalPart: principalPart,
		InterestPart:  interestPart,
		PaymentDate:   time.Now(),
		IsPartial:     isPartial,
		IsFinal:       isFinal,
	}
}

// EventName returns the event name
func (e InvestmentReturnPaidEvent) EventName() string {
	return "investment.return_paid"
}

// Payload returns the event payload
func (e InvestmentReturnPaidEvent) Payload() interface{} {
	return e
}