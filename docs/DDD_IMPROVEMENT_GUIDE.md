# Domain-Driven Design Improvement Guide for Daneizo

## Executive Summary

This guide provides a comprehensive roadmap for transforming Daneizo's current **anemic domain model** into a **rich, behavior-driven domain** that truly captures the complexity and nuances of P2P lending business logic.

**Current State**: Procedural architecture with business logic in interactors  
**Target State**: Rich domain model with business logic encapsulated in entities and domain services  
**Improvement Impact**: Enhanced maintainability, testability, and business rule consistency

---

## 🔍 Current DDD Assessment

### Strengths
- ✅ Clear domain boundaries (`internal/loan/`)
- ✅ Ubiquitous language (Proposed, Approved, Invested, Disbursed)
- ✅ State machine implementation
- ✅ Clean interface segregation

### Critical Issues
- ❌ **Anemic Domain Model**: Entities are data containers without behavior
- ❌ **Scattered Business Logic**: Rules spread across interactors
- ❌ **Primitive Obsession**: Overuse of `uint64`, `decimal.Decimal`  
- ❌ **Missing Domain Services**: Complex calculations in wrong layer
- ❌ **No Domain Events**: State changes don't trigger business events
- ❌ **Weak Aggregates**: No consistency boundaries or invariants

---

## 🎯 Domain Model Transformation

### Phase 1: Value Objects (Foundation)

Replace primitive types with business-meaningful value objects:

#### Current Problem
```go
// Primitive obsession - no business meaning
type CreateProposedLoanInput struct {
    UserID       uint64          `json:"user_id"`
    InterestRate decimal.Decimal `json:"interest_rate"`  
    Amount       decimal.Decimal `json:"amount"`
}
```

#### Improved Solution
```go
// Rich value objects with business validation
package domain

import "errors"

// Money value object with currency awareness
type Money struct {
    amount   decimal.Decimal
    currency Currency
}

func NewMoney(amount decimal.Decimal, currency Currency) (Money, error) {
    if amount.IsNegative() {
        return Money{}, errors.New("money amount cannot be negative")
    }
    return Money{amount: amount, currency: currency}, nil
}

func (m Money) Add(other Money) (Money, error) {
    if m.currency != other.currency {
        return Money{}, errors.New("cannot add different currencies")
    }
    return Money{amount: m.amount.Add(other.amount), currency: m.currency}, nil
}

func (m Money) IsZero() bool { return m.amount.IsZero() }
func (m Money) Amount() decimal.Decimal { return m.amount }

// Interest Rate with business validation
type InterestRate struct {
    rate decimal.Decimal // as percentage
}

func NewInterestRate(rate decimal.Decimal) (InterestRate, error) {
    if rate.IsNegative() || rate.GreaterThan(decimal.NewFromInt(30)) {
        return InterestRate{}, errors.New("interest rate must be between 0-30%")
    }
    return InterestRate{rate: rate}, nil
}

func (ir InterestRate) CalculateInterest(principal Money, termMonths int) Money {
    monthlyRate := ir.rate.Div(decimal.NewFromInt(100)).Div(decimal.NewFromInt(12))
    interest := principal.amount.Mul(monthlyRate).Mul(decimal.NewFromInt(int64(termMonths)))
    result, _ := NewMoney(interest, principal.currency)
    return result
}

// Strongly-typed IDs
type LoanID struct {
    id uint64
}

func NewLoanID(id uint64) LoanID {
    return LoanID{id: id}
}

func (l LoanID) Value() uint64 { return l.id }

type UserID struct {
    id uint64  
}

func NewUserID(id uint64) UserID {
    return UserID{id: id}
}
```

---

### Phase 2: Rich Domain Entities

Transform anemic entities into behavior-rich aggregates:

#### Current Problem
```go
// Anemic entity - just data container
func (c *CreateProposedLoan) Execute(ctx context.Context, in usecase.CreateProposedLoanInput) error {
    // Business logic in interactor ❌
    if err := c.store.InsertLoan(ctx, sqlentity.Loan{
        ID:              c.snowflakeGen.Generate(),
        BorrowerID:      in.UserID,
        PrincipalAmount: in.Amount,
        InterestRate:    in.InterestRate,
        InvestedAmount:  decimal.Zero,
        Status:          sqlentity.Proposed,
    }); err != nil {
        return err
    }
    return nil
}
```

#### Improved Solution
```go
// Rich domain aggregate with business behavior
package domain

type Loan struct {
    id              LoanID
    borrowerID      UserID
    principalAmount Money
    investedAmount  Money
    interestRate    InterestRate
    status          LoanStatus
    approvalDate    *time.Time
    approverID      *UserID
    investments     []Investment
    
    // Domain events (to be published)
    events []DomainEvent
}

// Factory method with business validation
func NewLoanProposal(
    id LoanID, 
    borrowerID UserID, 
    amount Money, 
    rate InterestRate,
) (*Loan, error) {
    
    // Business validation
    if amount.Amount().LessThan(decimal.NewFromInt(1000)) {
        return nil, errors.New("minimum loan amount is $1000")
    }
    
    loan := &Loan{
        id:              id,
        borrowerID:      borrowerID,
        principalAmount: amount,
        investedAmount:  Money{amount: decimal.Zero, currency: amount.currency},
        interestRate:    rate,
        status:          Proposed,
        investments:     make([]Investment, 0),
        events:          make([]DomainEvent, 0),
    }
    
    // Generate domain event
    loan.AddEvent(LoanProposedEvent{
        LoanID:    id,
        BorrowerID: borrowerID,
        Amount:    amount,
        Rate:      rate,
        Timestamp: time.Now(),
    })
    
    return loan, nil
}

// Business methods encapsulate domain logic
func (l *Loan) Approve(approverID UserID) error {
    // Business rule: Can only approve proposed loans
    if l.status != Proposed {
        return NewDomainError("can only approve proposed loans")
    }
    
    l.status = Approved
    l.approverID = &approverID
    now := time.Now()
    l.approvalDate = &now
    
    // Generate domain event
    l.AddEvent(LoanApprovedEvent{
        LoanID:     l.id,
        ApproverID: approverID,
        Timestamp:  now,
    })
    
    return nil
}

func (l *Loan) AddInvestment(investorID UserID, amount Money) error {
    // Business rule: Can only invest in approved loans
    if l.status != Approved {
        return NewDomainError("can only invest in approved loans")
    }
    
    // Business rule: Cannot over-invest
    newTotal, err := l.investedAmount.Add(amount)
    if err != nil {
        return err
    }
    
    if newTotal.Amount().GreaterThan(l.principalAmount.Amount()) {
        return NewDomainError("investment would exceed loan amount")
    }
    
    investment := Investment{
        investorID: investorID,
        amount:     amount,
        timestamp:  time.Now(),
    }
    
    l.investments = append(l.investments, investment)
    l.investedAmount = newTotal
    
    // Check if fully funded
    if l.investedAmount.Amount().Equal(l.principalAmount.Amount()) {
        l.status = FullyFunded
        l.AddEvent(LoanFullyFundedEvent{
            LoanID:    l.id,
            Timestamp: time.Now(),
        })
    }
    
    l.AddEvent(InvestmentMadeEvent{
        LoanID:     l.id,
        InvestorID: investorID,
        Amount:     amount,
        Timestamp:  time.Now(),
    })
    
    return nil
}

func (l *Loan) Disburse(agreementURL string) error {
    if l.status != FullyFunded {
        return NewDomainError("can only disburse fully funded loans")
    }
    
    l.status = Disbursed
    l.AddEvent(LoanDisbursedEvent{
        LoanID:       l.id,
        AgreementURL: agreementURL,
        Timestamp:    time.Now(),
    })
    
    return nil
}

// Aggregate invariants
func (l *Loan) IsValid() error {
    if l.investedAmount.Amount().GreaterThan(l.principalAmount.Amount()) {
        return errors.New("invested amount exceeds principal")
    }
    return nil
}

// Domain events support
func (l *Loan) AddEvent(event DomainEvent) {
    l.events = append(l.events, event)
}

func (l *Loan) Events() []DomainEvent {
    return l.events
}

func (l *Loan) ClearEvents() {
    l.events = make([]DomainEvent, 0)
}
```

---

### Phase 3: Domain Services

Extract complex business logic into dedicated domain services:

```go
package domain

// Risk Assessment Domain Service
type RiskAssessmentService struct {
    creditBureauAPI CreditBureauPort
    fraudDetection  FraudDetectionPort
}

type RiskScore int

const (
    LowRisk RiskScore = iota + 1
    MediumRisk
    HighRisk
)

func (ras *RiskAssessmentService) AssessLoanRisk(
    borrower Borrower,
    loanAmount Money,
) (RiskScore, error) {
    
    // Complex business logic for risk assessment
    creditScore, err := ras.creditBureauAPI.GetCreditScore(borrower.ID())
    if err != nil {
        return 0, err
    }
    
    fraudRisk, err := ras.fraudDetection.CheckFraud(borrower.ID())
    if err != nil {
        return 0, err
    }
    
    // Business rules for risk calculation
    if creditScore >= 750 && !fraudRisk && 
       loanAmount.Amount().LessThan(decimal.NewFromInt(50000)) {
        return LowRisk, nil
    }
    
    if creditScore >= 650 && !fraudRisk {
        return MediumRisk, nil
    }
    
    return HighRisk, nil
}

// Interest Rate Calculation Domain Service
type InterestRateCalculationService struct{}

func (ircs *InterestRateCalculationService) CalculateRate(
    riskScore RiskScore,
    loanAmount Money,
    termMonths int,
) (InterestRate, error) {
    
    baseRate := decimal.NewFromFloat(5.0) // 5% base rate
    
    // Risk adjustment
    switch riskScore {
    case LowRisk:
        // No adjustment
    case MediumRisk:
        baseRate = baseRate.Add(decimal.NewFromFloat(2.0)) // +2%
    case HighRisk:
        baseRate = baseRate.Add(decimal.NewFromFloat(5.0)) // +5%
    }
    
    // Amount adjustment (smaller loans = higher rates)
    if loanAmount.Amount().LessThan(decimal.NewFromInt(10000)) {
        baseRate = baseRate.Add(decimal.NewFromFloat(1.0)) // +1%
    }
    
    // Term adjustment
    if termMonths > 36 {
        baseRate = baseRate.Add(decimal.NewFromFloat(0.5)) // +0.5%
    }
    
    return NewInterestRate(baseRate)
}
```

---

### Phase 4: Domain Events

Implement domain events for decoupled communication:

```go
package domain

// Domain Event interface
type DomainEvent interface {
    EventName() string
    EventID() string
    Timestamp() time.Time
    AggregateID() string
}

// Specific domain events
type LoanProposedEvent struct {
    eventID   string
    LoanID    LoanID
    BorrowerID UserID
    Amount    Money
    Rate      InterestRate
    timestamp time.Time
}

func (e LoanProposedEvent) EventName() string { return "loan.proposed" }
func (e LoanProposedEvent) EventID() string { return e.eventID }
func (e LoanProposedEvent) Timestamp() time.Time { return e.timestamp }
func (e LoanProposedEvent) AggregateID() string { return fmt.Sprintf("%d", e.LoanID.Value()) }

type LoanApprovedEvent struct {
    eventID    string
    LoanID     LoanID
    ApproverID UserID
    timestamp  time.Time
}

type InvestmentMadeEvent struct {
    eventID    string
    LoanID     LoanID
    InvestorID UserID
    Amount     Money
    timestamp  time.Time
}

// Event Publisher
type DomainEventPublisher interface {
    Publish(events []DomainEvent) error
}
```

---

### Phase 5: Repository Pattern

Proper repository implementation for aggregates:

```go
package domain

// Domain repository interface
type LoanRepository interface {
    Save(loan *Loan) error
    FindByID(id LoanID) (*Loan, error)
    FindByBorrower(borrowerID UserID) ([]*Loan, error)
    FindByStatus(status LoanStatus) ([]*Loan, error)
    FindRequiringApproval() ([]*Loan, error)
}

// Implementation in infrastructure layer
package infrastructure

type SQLLoanRepository struct {
    db           *sql.DB
    eventPublisher domain.DomainEventPublisher
}

func (r *SQLLoanRepository) Save(loan *domain.Loan) error {
    tx, err := r.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // Save aggregate root
    if err := r.saveLoan(tx, loan); err != nil {
        return err
    }
    
    // Save child entities (investments)
    if err := r.saveInvestments(tx, loan); err != nil {
        return err
    }
    
    // Publish domain events
    if err := r.eventPublisher.Publish(loan.Events()); err != nil {
        return err
    }
    
    loan.ClearEvents()
    
    return tx.Commit()
}
```

---

## 🚀 Implementation Strategy

### Migration Phases

#### Phase 1: Foundation (Weeks 1-2)
1. **Value Objects Introduction**
   - Create `Money`, `InterestRate`, `LoanID`, `UserID`
   - Update DTOs to use value objects
   - Add validation logic

2. **Entity Enhancement**  
   - Add business methods to `Loan` entity
   - Move validation from interactors to entities
   - Create aggregate root pattern

#### Phase 2: Business Logic (Weeks 3-4)
3. **Domain Services**
   - Extract risk assessment logic
   - Create interest rate calculation service
   - Implement loan matching service

4. **Interactor Refactoring**
   - Slim down interactors to orchestration
   - Move business logic to domain layer
   - Update tests

#### Phase 3: Events & Advanced (Weeks 5-6)
5. **Domain Events**
   - Implement event infrastructure
   - Add event publishing
   - Create event handlers

6. **Repository Pattern**
   - Implement proper aggregate repositories
   - Add optimistic locking
   - Update persistence layer

### Backward Compatibility

Maintain existing APIs during migration:

```go
// Adapter pattern for gradual migration
type LegacyInteractorAdapter struct {
    domainService *domain.LoanService
}

func (a *LegacyInteractorAdapter) Execute(
    ctx context.Context,
    in usecase.CreateProposedLoanInput,
) error {
    // Convert legacy input to domain objects
    loanID := domain.NewLoanID(a.idGenerator.Generate())
    borrowerID := domain.NewUserID(in.UserID)
    amount, _ := domain.NewMoney(in.Amount, domain.USD)
    rate, _ := domain.NewInterestRate(in.InterestRate)
    
    // Use domain service
    return a.domainService.ProposeNewLoan(ctx, loanID, borrowerID, amount, rate)
}
```

---

## 📊 Success Metrics

### Code Quality Improvements
- **Cyclomatic Complexity**: Reduce by 40% through proper domain modeling
- **Business Logic Centralization**: 90% of business rules in domain layer
- **Test Coverage**: Increase domain layer coverage to 95%

### Business Value
- **Rule Consistency**: Centralized business rules prevent inconsistencies
- **Feature Velocity**: Faster feature development through rich domain model
- **Maintainability**: Easier to modify business rules without touching infrastructure

### Technical Benefits
- **Type Safety**: Eliminate primitive obsession bugs
- **Event-Driven**: Enable loosely coupled integrations
- **Testing**: Isolated business logic testing

---

## 🔗 Next Steps

1. **Start with Value Objects** - Low risk, high impact foundation
2. **Gradually Move Business Logic** - One use case at a time
3. **Add Domain Events** - Enable future scalability patterns
4. **Implement CQRS** - Separate read/write models as system grows
5. **Consider Event Sourcing** - For full audit trail and temporal queries

This transformation will position Daneizo as a maintainable, scalable P2P lending platform with business logic that truly reflects the domain complexity.