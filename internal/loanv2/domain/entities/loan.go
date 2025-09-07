package entities

import (
	"time"
	"github.com/shopspring/decimal"
	"github.com/farislr/daneizo/internal/loanv2/domain/valueobjects"
	"github.com/farislr/daneizo/internal/loanv2/domain/events"
	"github.com/farislr/daneizo/internal/loanv2/domain/errors"
)

// Loan represents the loan aggregate root with rich business behavior
type Loan struct {
	// Core identifiers and amounts
	id              valueobjects.LoanID
	borrowerID      valueobjects.UserID
	principalAmount valueobjects.Money
	investedAmount  valueobjects.Money
	interestRate    valueobjects.InterestRate
	termMonths      int
	loanPurpose     string

	// Status and lifecycle
	status       valueobjects.LoanStatus
	proposedAt   time.Time
	approvedAt   *time.Time
	approverID   *valueobjects.UserID
	disbursedAt  *time.Time
	firstPaymentDue *time.Time

	// Related entities
	investments []Investment

	// Documents and agreements
	agreementDocumentURL *string

	// Domain events for decoupled communication
	domainEvents []events.DomainEvent

	// Audit fields
	createdAt time.Time
	updatedAt time.Time
}

// NewLoanProposal creates a new loan proposal with business validation
func NewLoanProposal(
	id valueobjects.LoanID,
	borrowerID valueobjects.UserID,
	principalAmount valueobjects.Money,
	interestRate valueobjects.InterestRate,
	termMonths int,
	purpose string,
) (*Loan, error) {
	// Business validation at creation
	if err := validateLoanProposal(principalAmount, interestRate, termMonths); err != nil {
		return nil, err
	}

	now := time.Now()
	loan := &Loan{
		id:              id,
		borrowerID:      borrowerID,
		principalAmount: principalAmount,
		investedAmount:  valueobjects.NewZeroMoney(principalAmount.Currency()),
		interestRate:    interestRate,
		termMonths:      termMonths,
		loanPurpose:     purpose,
		status:          valueobjects.Proposed,
		proposedAt:      now,
		investments:     make([]Investment, 0),
		domainEvents:    make([]events.DomainEvent, 0),
		createdAt:       now,
		updatedAt:       now,
	}

	// Generate domain event for system integration
	event := events.NewLoanProposedEvent(
		id, borrowerID, principalAmount, interestRate, termMonths, purpose, id.String(),
	)
	loan.AddDomainEvent(event)

	return loan, nil
}

// validateLoanProposal performs business validation for loan creation
func validateLoanProposal(amount valueobjects.Money, rate valueobjects.InterestRate, termMonths int) error {
	// Minimum loan amount validation
	minAmount, _ := valueobjects.NewMoney(decimal.NewFromFloat(1000), amount.Currency())
	if isLess, _ := amount.LessThan(minAmount); isLess {
		return errors.ErrMinimumLoanAmount
	}

	// Maximum loan amount validation
	maxAmount, _ := valueobjects.NewMoney(decimal.NewFromFloat(1000000), amount.Currency())
	if isGreater, _ := amount.GreaterThan(maxAmount); isGreater {
		return errors.ErrMaximumLoanAmount
	}

	// Term validation
	if termMonths < 6 || termMonths > 60 {
		return errors.ErrInvalidLoanTerm
	}

	// Interest rate validation (already validated in InterestRate value object, but additional business rules)
	if rate.Value().GreaterThan(decimal.NewFromFloat(30.0)) {
		return errors.ErrInvalidInterestRate
	}

	return nil
}

// Approve approves the loan with business rule validation
func (l *Loan) Approve(approverID valueobjects.UserID, conditions []string) error {
	// Business rule validation
	if l.status != valueobjects.Proposed && l.status != valueobjects.UnderReview {
		return errors.ErrLoanNotProposed
	}

	if l.borrowerID.Equal(approverID) {
		return errors.NewLoanBusinessRuleError("borrower cannot approve their own loan")
	}

	// State transition with business logic
	l.status = valueobjects.Approved
	l.approverID = &approverID
	now := time.Now()
	l.approvedAt = &now
	l.updatedAt = now

	// Domain event for system notification
	event := events.NewLoanApprovedEvent(
		l.id, approverID, l.principalAmount, l.interestRate, conditions, l.id.String(),
	)
	l.AddDomainEvent(event)

	return nil
}

// Reject rejects the loan application
func (l *Loan) Reject(reviewerID valueobjects.UserID, reason string, details []string) error {
	// Business rule validation
	if l.status != valueobjects.Proposed && l.status != valueobjects.UnderReview {
		return errors.ErrLoanNotProposed
	}

	if l.status.IsTerminal() {
		return errors.NewLoanBusinessRuleError("cannot reject loan in terminal status")
	}

	// State transition
	l.status = valueobjects.Rejected
	l.updatedAt = time.Now()

	// Domain event
	event := events.NewLoanRejectedEvent(l.id, reviewerID, reason, details, l.id.String())
	l.AddDomainEvent(event)

	return nil
}

// AddInvestment adds an investment to the loan with comprehensive business rules
func (l *Loan) AddInvestment(
	investorID valueobjects.UserID,
	amount valueobjects.Money,
	expectedReturn valueobjects.Money,
) error {
	// Business rule validation
	if !l.status.CanAcceptInvestments() {
		return errors.ErrLoanNotInvestable
	}

	// Prevent self-investment
	if l.borrowerID.Equal(investorID) {
		return errors.ErrSelfInvestment
	}

	// Validate investment amount
	if amount.IsZero() {
		return errors.ErrMinimumInvestmentAmount
	}

	// Business rule: Prevent over-investment
	newTotal, err := l.investedAmount.Add(amount)
	if err != nil {
		return err
	}

	if isGreater, _ := newTotal.GreaterThan(l.principalAmount); isGreater {
		return errors.ErrInvestmentExceedsLoan
	}

	// Create investment entity
	investment := NewInvestment(investorID, amount, expectedReturn, time.Now())
	l.investments = append(l.investments, investment)
	l.investedAmount = newTotal
	l.updatedAt = time.Now()

	// Determine remaining amount needed
	remainingAmount, _ := l.principalAmount.Subtract(l.investedAmount)

	// State management based on funding level
	var isFullyFunded bool
	if l.investedAmount.Equal(l.principalAmount) {
		l.status = valueobjects.FullyFunded
		isFullyFunded = true
		
		// Emit fully funded event
		fundedEvent := events.NewLoanFullyFundedEvent(
			l.id, l.investedAmount, len(l.investments), true, l.id.String(),
		)
		l.AddDomainEvent(fundedEvent)
	} else {
		l.status = valueobjects.PartiallyFunded
		isFullyFunded = false
	}

	// Emit investment made event
	event := events.NewInvestmentMadeEvent(
		l.id, investorID, investment.id, amount, expectedReturn, 
		remainingAmount, isFullyFunded, l.id.String(),
	)
	l.AddDomainEvent(event)

	return nil
}

// Disburse disburses the loan funds with strict business rules
func (l *Loan) Disburse(agreementURL string, firstPaymentDue time.Time) error {
	// Business rule: Only fully funded loans can be disbursed
	if l.status != valueobjects.FullyFunded {
		return errors.ErrLoanNotFullyFunded
	}

	// Additional business validation
	if agreementURL == "" {
		return errors.ErrAgreementDocumentRequired
	}

	if firstPaymentDue.Before(time.Now().AddDate(0, 1, 0)) {
		return errors.NewLoanBusinessRuleError("first payment must be at least 1 month from disbursement")
	}

	// State transition
	l.status = valueobjects.Disbursed
	now := time.Now()
	l.disbursedAt = &now
	l.agreementDocumentURL = &agreementURL
	l.firstPaymentDue = &firstPaymentDue
	l.updatedAt = now

	// Domain event
	event := events.NewLoanDisbursedEvent(
		l.id, l.borrowerID, l.principalAmount, agreementURL, firstPaymentDue, l.id.String(),
	)
	l.AddDomainEvent(event)

	return nil
}

// Cancel cancels the loan before disbursement
func (l *Loan) Cancel(userID valueobjects.UserID, reason string) error {
	// Can only cancel loans that haven't been disbursed
	if l.status == valueobjects.Disbursed || l.status == valueobjects.Active || l.status.IsTerminal() {
		return errors.NewLoanBusinessRuleError("cannot cancel loan in current status")
	}

	l.status = valueobjects.Cancelled
	l.updatedAt = time.Now()

	// Domain event
	event := events.NewLoanStatusChangedEvent(
		l.id, l.status, valueobjects.Cancelled, reason, &userID, l.id.String(),
	)
	l.AddDomainEvent(event)

	return nil
}

// MarkAsActive transitions the loan to active status after disbursement
func (l *Loan) MarkAsActive() error {
	if l.status != valueobjects.Disbursed {
		return errors.NewLoanBusinessRuleError("loan must be disbursed before becoming active")
	}

	l.status = valueobjects.Active
	l.updatedAt = time.Now()

	return nil
}

// Aggregate invariants to ensure business consistency
func (l *Loan) IsValid() error {
	if isGreater, _ := l.investedAmount.GreaterThan(l.principalAmount); isGreater {
		return errors.NewLoanBusinessRuleError("invested amount cannot exceed principal amount")
	}

	if l.status == valueobjects.Approved && len(l.investments) > 0 {
		return errors.NewLoanBusinessRuleError("approved loans should not have investments yet")
	}

	if l.status == valueobjects.Disbursed {
		if !l.investedAmount.Equal(l.principalAmount) {
			return errors.NewLoanBusinessRuleError("disbursed loans must be fully funded")
		}
		if l.agreementDocumentURL == nil {
			return errors.ErrAgreementDocumentRequired
		}
	}

	return nil
}

// Domain event management
func (l *Loan) AddDomainEvent(event events.DomainEvent) {
	l.domainEvents = append(l.domainEvents, event)
}

func (l *Loan) DomainEvents() []events.DomainEvent {
	return l.domainEvents
}

func (l *Loan) ClearDomainEvents() {
	l.domainEvents = make([]events.DomainEvent, 0)
}

// Getters for read access (encapsulation)
func (l *Loan) ID() valueobjects.LoanID { return l.id }
func (l *Loan) BorrowerID() valueobjects.UserID { return l.borrowerID }
func (l *Loan) PrincipalAmount() valueobjects.Money { return l.principalAmount }
func (l *Loan) InvestedAmount() valueobjects.Money { return l.investedAmount }
func (l *Loan) InterestRate() valueobjects.InterestRate { return l.interestRate }
func (l *Loan) TermMonths() int { return l.termMonths }
func (l *Loan) LoanPurpose() string { return l.loanPurpose }
func (l *Loan) Status() valueobjects.LoanStatus { return l.status }
func (l *Loan) ProposedAt() time.Time { return l.proposedAt }
func (l *Loan) ApprovedAt() *time.Time { return l.approvedAt }
func (l *Loan) ApproverID() *valueobjects.UserID { return l.approverID }
func (l *Loan) DisbursedAt() *time.Time { return l.disbursedAt }
func (l *Loan) FirstPaymentDue() *time.Time { return l.firstPaymentDue }
func (l *Loan) Investments() []Investment { return l.investments }
func (l *Loan) AgreementDocumentURL() *string { return l.agreementDocumentURL }
func (l *Loan) CreatedAt() time.Time { return l.createdAt }
func (l *Loan) UpdatedAt() time.Time { return l.updatedAt }

// Business query methods
func (l *Loan) RemainingAmount() valueobjects.Money {
	remaining, _ := l.principalAmount.Subtract(l.investedAmount)
	return remaining
}

func (l *Loan) IsFullyFunded() bool {
	return l.investedAmount.Equal(l.principalAmount)
}

func (l *Loan) FundingPercentage() float64 {
	if l.principalAmount.IsZero() {
		return 0
	}
	return l.investedAmount.Amount().Div(l.principalAmount.Amount()).InexactFloat64() * 100
}

func (l *Loan) NumberOfInvestors() int {
	return len(l.investments)
}

func (l *Loan) CalculateMonthlyPayment() valueobjects.Money {
	return l.interestRate.CalculateMonthlyPayment(l.principalAmount, l.termMonths)
}

func (l *Loan) CalculateTotalInterest() valueobjects.Money {
	return l.interestRate.CalculateTotalInterest(l.principalAmount, l.termMonths)
}