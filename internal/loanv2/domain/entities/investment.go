package entities

import (
	"time"
	"github.com/google/uuid"
	"github.com/farislr/daneizo/internal/loanv2/domain/valueobjects"
	"github.com/farislr/daneizo/internal/loanv2/domain/errors"
)

// InvestmentStatus represents the status of an investment
type InvestmentStatus string

const (
	InvestmentActive     InvestmentStatus = "ACTIVE"
	InvestmentWithdrawn  InvestmentStatus = "WITHDRAWN"
	InvestmentMatured    InvestmentStatus = "MATURED"
	InvestmentDefaulted  InvestmentStatus = "DEFAULTED"
)

// Investment represents an investment made in a loan
type Investment struct {
	id             string
	investorID     valueobjects.UserID
	amount         valueobjects.Money
	expectedReturn valueobjects.Money
	actualReturn   valueobjects.Money
	status         InvestmentStatus
	investmentDate time.Time
	maturityDate   *time.Time
	withdrawnDate  *time.Time
	createdAt      time.Time
	updatedAt      time.Time
}

// NewInvestment creates a new investment
func NewInvestment(
	investorID valueobjects.UserID,
	amount valueobjects.Money,
	expectedReturn valueobjects.Money,
	investmentDate time.Time,
) Investment {
	now := time.Now()
	return Investment{
		id:             uuid.New().String(),
		investorID:     investorID,
		amount:         amount,
		expectedReturn: expectedReturn,
		actualReturn:   valueobjects.NewZeroMoney(amount.Currency()),
		status:         InvestmentActive,
		investmentDate: investmentDate,
		createdAt:      now,
		updatedAt:      now,
	}
}

// Withdraw marks the investment as withdrawn
func (i *Investment) Withdraw(withdrawnDate time.Time) error {
	if i.status != InvestmentActive {
		return errors.NewInvestmentBusinessRuleError("can only withdraw active investments")
	}
	
	i.status = InvestmentWithdrawn
	i.withdrawnDate = &withdrawnDate
	i.updatedAt = time.Now()
	
	return nil
}

// AddReturn adds to the actual return received
func (i *Investment) AddReturn(returnAmount valueobjects.Money) error {
	newReturn, err := i.actualReturn.Add(returnAmount)
	if err != nil {
		return err
	}
	
	i.actualReturn = newReturn
	i.updatedAt = time.Now()
	
	return nil
}

// MarkAsMatured marks the investment as matured
func (i *Investment) MarkAsMatured(maturityDate time.Time) {
	i.status = InvestmentMatured
	i.maturityDate = &maturityDate
	i.updatedAt = time.Now()
}

// MarkAsDefaulted marks the investment as defaulted
func (i *Investment) MarkAsDefaulted() {
	i.status = InvestmentDefaulted
	i.updatedAt = time.Now()
}

// Getters
func (i Investment) ID() string { return i.id }
func (i Investment) InvestorID() valueobjects.UserID { return i.investorID }
func (i Investment) Amount() valueobjects.Money { return i.amount }
func (i Investment) ExpectedReturn() valueobjects.Money { return i.expectedReturn }
func (i Investment) ActualReturn() valueobjects.Money { return i.actualReturn }
func (i Investment) Status() InvestmentStatus { return i.status }
func (i Investment) InvestmentDate() time.Time { return i.investmentDate }
func (i Investment) MaturityDate() *time.Time { return i.maturityDate }
func (i Investment) WithdrawnDate() *time.Time { return i.withdrawnDate }
func (i Investment) CreatedAt() time.Time { return i.createdAt }
func (i Investment) UpdatedAt() time.Time { return i.updatedAt }

// Business methods
func (i Investment) IsActive() bool {
	return i.status == InvestmentActive
}

func (i Investment) IsWithdrawn() bool {
	return i.status == InvestmentWithdrawn
}

func (i Investment) IsMatured() bool {
	return i.status == InvestmentMatured
}

func (i Investment) IsDefaulted() bool {
	return i.status == InvestmentDefaulted
}

func (i Investment) ReturnRate() float64 {
	if i.amount.IsZero() {
		return 0
	}
	return i.actualReturn.Amount().Div(i.amount.Amount()).InexactFloat64()
}