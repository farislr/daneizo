# Comprehensive Domain-Driven Design Architecture Guide for Daneizo

## Table of Contents
1. [Executive Summary](#executive-summary)
2. [Current Architecture Assessment](#current-architecture-assessment)
3. [Domain-Driven Design Analysis](#domain-driven-design-analysis)
4. [Clean Architecture & Hexagonal Patterns](#clean-architecture--hexagonal-patterns)
5. [Target Architecture Design](#target-architecture-design)
6. [Detailed Implementation Strategy](#detailed-implementation-strategy)
7. [Project Structure Transformation](#project-structure-transformation)
8. [Migration Roadmap](#migration-roadmap)
9. [Success Metrics & Validation](#success-metrics--validation)
10. [Future Architecture Considerations](#future-architecture-considerations)

---

## Executive Summary

This comprehensive guide provides a complete roadmap for transforming Daneizo, a P2P lending state engine, from its current **anemic domain model** into a **rich, behavior-driven architecture** that fully embraces Domain-Driven Design principles while maintaining Clean Architecture and Hexagonal patterns.

### Current State Analysis
**Architecture Grade**: A- (8.3/10)
- **Strengths**: Clear layering, dependency inversion, excellent testability
- **Critical Gap**: Anemic domain model with business logic scattered across interactors

### Transformation Objectives
**Target State**: Rich domain model with 90% business logic centralized in domain layer
- **Expected Impact**: 40% reduction in cyclomatic complexity, 60% faster bug fixes
- **Business Value**: Consistent business rules, improved maintainability, enhanced collaboration with domain experts

### Key Transformation Areas
1. **Domain Layer Enrichment**: From data containers to behavior-rich aggregates
2. **Value Objects**: Eliminate primitive obsession with business-meaningful types
3. **Domain Services**: Extract complex business logic into dedicated services
4. **Event-Driven Architecture**: Implement domain events for loose coupling
5. **Project Structure**: Reorganize for proper DDD layer separation

---

## Current Architecture Assessment

### 🏗️ Architecture Overview

Daneizo currently implements a **hybrid Clean Architecture** with strong **Hexagonal Architecture** patterns, demonstrating excellent structural discipline while suffering from domain model anemia.

#### Current Layer Implementation
```
┌─────────────────────┐
│   HTTP Endpoints    │ ← Frameworks & Drivers (8/10)
├─────────────────────┤
│   HTTP/SQL Gateway  │ ← Interface Adapters (8/10)
├─────────────────────┤
│    Interactors      │ ← Use Cases (7/10) ⚠️ Business logic scattered
├─────────────────────┤
│  Entities/UseCase   │ ← Enterprise Rules (6/10) ⚠️ Anemic model
└─────────────────────┘
```

#### Current Project Structure (Problems Identified)
```
internal/loan/
├── internal/
│   ├── entity/sqlentity/     ❌ Data containers only
│   ├── gateway/              ❌ Mixed concerns (HTTP + SQL)
│   ├── interactor/           ❌ Business logic scattered
│   ├── usecase/              ❌ Input/Output structs
│   └── mocks/                ✅ Generated mocks
└── module.go                 ✅ Dependency injection
```

### 🎯 Detailed Assessment Scores

#### Domain-Driven Design: 7/10
**Strengths:**
- ✅ Clear bounded context (`internal/loan/`)
- ✅ Ubiquitous language (Proposed → Approved → Invested → Disbursed)
- ✅ State machine implementation
- ✅ Clean interface segregation

**Critical Issues:**
- ❌ **Anemic Domain Model**: Entities are data containers without behavior
- ❌ **Scattered Business Logic**: Rules spread across interactors
- ❌ **Primitive Obsession**: Overuse of `uint64`, `decimal.Decimal`
- ❌ **Missing Domain Services**: Complex calculations in wrong layer
- ❌ **No Domain Events**: State changes don't trigger business events

#### Clean Architecture: 8/10
**Strengths:**
- ✅ Perfect dependency rule compliance
- ✅ Framework-independent business logic
- ✅ Clear use case definitions
- ✅ Interface segregation principle

**Minor Issues:**
- ⚠️ Cross-cutting concerns scattered (logging, validation)
- ⚠️ No centralized error handling strategy

#### Hexagonal Architecture: 8/10
**Strengths:**
- ✅ Clear port and adapter pattern implementation
- ✅ Mockable interfaces for all external dependencies
- ✅ Easy to swap implementations
- ✅ Type safety through Go interfaces

**Enhancement Opportunities:**
- ⚠️ Adapter standardization needed
- ⚠️ Missing health check interfaces
- ⚠️ Configuration management scattered

---

## Domain-Driven Design Analysis

### 🔍 Current Domain Model Issues

#### 1. Anemic Entity Example (Current Problem)
```go
// Current: Anemic entity - just data container
type Loan struct {
    ID                         uint64
    BorrowerID                 uint64
    PrincipalAmount            decimal.Decimal
    InvestedAmount             decimal.Decimal
    InterestRate               decimal.Decimal
    Status                     LoanStatus
    ApprovalDate               sql.NullTime
    ApprovalEmployeeID         sql.NullInt64
    DisbursementDate           sql.NullTime
    AgreementLetterDocumentURL sql.NullString
    CreatedAt                  time.Time
    UpdatedAt                  time.Time
}

// Business logic scattered in interactors ❌
func (c *CreateProposedLoan) Execute(ctx context.Context, in usecase.CreateProposedLoanInput) error {
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

#### 2. Primitive Obsession Issues
```go
// Current: Primitive obsession - no business meaning
type CreateProposedLoanInput struct {
    UserID       uint64          `json:"user_id"`       // Could be any ID
    InterestRate decimal.Decimal `json:"interest_rate"` // No validation
    Amount       decimal.Decimal `json:"amount"`        // No currency info
}

// Problems:
// - No validation at type level
// - No business meaning
// - Easy to mix up different types of IDs
// - Currency confusion possible
// - Interest rate ranges not enforced
```

### 🎯 Target Domain Model Design

#### 1. Rich Aggregate Root (Target Solution)
```go
// Rich domain aggregate with business behavior
package domain

type Loan struct {
    // Value objects provide type safety and business validation
    id              LoanID
    borrowerID      UserID
    principalAmount Money
    investedAmount  Money
    interestRate    InterestRate
    status          LoanStatus
    approvalDate    *time.Time
    approverID      *UserID
    investments     []Investment
    
    // Domain events for decoupled communication
    events []DomainEvent
}

// Factory method with business validation
func NewLoanProposal(
    id LoanID, 
    borrowerID UserID, 
    amount Money, 
    rate InterestRate,
) (*Loan, error) {
    
    // Business validation at creation
    if amount.Amount().LessThan(decimal.NewFromInt(1000)) {
        return nil, NewDomainError("minimum loan amount is $1000")
    }
    
    if rate.Value().GreaterThan(decimal.NewFromFloat(30.0)) {
        return nil, NewDomainError("maximum interest rate is 30%")
    }
    
    loan := &Loan{
        id:              id,
        borrowerID:      borrowerID,
        principalAmount: amount,
        investedAmount:  NewZeroMoney(amount.Currency()),
        interestRate:    rate,
        status:          Proposed,
        investments:     make([]Investment, 0),
        events:          make([]DomainEvent, 0),
    }
    
    // Generate domain event for system integration
    loan.AddEvent(NewLoanProposedEvent(id, borrowerID, amount, rate))
    
    return loan, nil
}

// Business methods encapsulate domain logic and rules
func (l *Loan) Approve(approverID UserID) error {
    // Business rule validation
    if l.status != Proposed {
        return NewDomainError("can only approve proposed loans")
    }
    
    // State transition with business logic
    l.status = Approved
    l.approverID = &approverID
    now := time.Now()
    l.approvalDate = &now
    
    // Domain event for system notification
    l.AddEvent(NewLoanApprovedEvent(l.id, approverID, now))
    
    return nil
}

func (l *Loan) AddInvestment(investorID UserID, amount Money) error {
    // Complex business rule: Investment validation
    if l.status != Approved && l.status != PartiallyFunded {
        return NewDomainError("can only invest in approved or partially funded loans")
    }
    
    // Business rule: Prevent over-investment
    newTotal, err := l.investedAmount.Add(amount)
    if err != nil {
        return err
    }
    
    if newTotal.Amount().GreaterThan(l.principalAmount.Amount()) {
        return NewDomainError("investment would exceed loan amount")
    }
    
    // Create investment entity
    investment := NewInvestment(investorID, amount, time.Now())
    l.investments = append(l.investments, investment)
    l.investedAmount = newTotal
    
    // State management based on funding level
    if l.investedAmount.Amount().Equal(l.principalAmount.Amount()) {
        l.status = FullyFunded
        l.AddEvent(NewLoanFullyFundedEvent(l.id, time.Now()))
    } else {
        l.status = PartiallyFunded
    }
    
    l.AddEvent(NewInvestmentMadeEvent(l.id, investorID, amount, time.Now()))
    
    return nil
}

func (l *Loan) Disburse(agreementURL string) error {
    // Business rule: Only fully funded loans can be disbursed
    if l.status != FullyFunded {
        return NewDomainError("can only disburse fully funded loans")
    }
    
    // Additional business validation
    if agreementURL == "" {
        return NewDomainError("agreement document URL is required for disbursement")
    }
    
    l.status = Disbursed
    l.AddEvent(NewLoanDisbursedEvent(l.id, agreementURL, time.Now()))
    
    return nil
}

// Aggregate invariants ensure business consistency
func (l *Loan) IsValid() error {
    if l.investedAmount.Amount().GreaterThan(l.principalAmount.Amount()) {
        return NewDomainError("invested amount cannot exceed principal amount")
    }
    
    if l.status == Approved && len(l.investments) > 0 {
        return NewDomainError("approved loans should not have investments yet")
    }
    
    if l.status == Disbursed && l.investedAmount.Amount().LessThan(l.principalAmount.Amount()) {
        return NewDomainError("disbursed loans must be fully funded")
    }
    
    return nil
}

// Domain event management
func (l *Loan) AddEvent(event DomainEvent) {
    l.events = append(l.events, event)
}

func (l *Loan) Events() []DomainEvent {
    return l.events
}

func (l *Loan) ClearEvents() {
    l.events = make([]DomainEvent, 0)
}

// Getters for read access (encapsulation)
func (l *Loan) ID() LoanID { return l.id }
func (l *Loan) BorrowerID() UserID { return l.borrowerID }
func (l *Loan) PrincipalAmount() Money { return l.principalAmount }
func (l *Loan) InvestedAmount() Money { return l.investedAmount }
func (l *Loan) Status() LoanStatus { return l.status }
func (l *Loan) Investments() []Investment { return l.investments }
```

#### 2. Value Objects for Type Safety
```go
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

func NewZeroMoney(currency Currency) Money {
    return Money{amount: decimal.Zero, currency: currency}
}

func (m Money) Add(other Money) (Money, error) {
    if m.currency != other.currency {
        return Money{}, errors.New("cannot add different currencies")
    }
    return Money{amount: m.amount.Add(other.amount), currency: m.currency}, nil
}

func (m Money) Subtract(other Money) (Money, error) {
    if m.currency != other.currency {
        return Money{}, errors.New("cannot subtract different currencies")
    }
    result := m.amount.Sub(other.amount)
    if result.IsNegative() {
        return Money{}, errors.New("result would be negative")
    }
    return Money{amount: result, currency: m.currency}, nil
}

func (m Money) MultiplyByRate(rate decimal.Decimal) Money {
    return Money{amount: m.amount.Mul(rate), currency: m.currency}
}

func (m Money) IsZero() bool { return m.amount.IsZero() }
func (m Money) Amount() decimal.Decimal { return m.amount }
func (m Money) Currency() Currency { return m.currency }

// Interest Rate with business validation and calculations
type InterestRate struct {
    rate decimal.Decimal // as percentage (0-100)
}

func NewInterestRate(rate decimal.Decimal) (InterestRate, error) {
    if rate.IsNegative() {
        return InterestRate{}, errors.New("interest rate cannot be negative")
    }
    if rate.GreaterThan(decimal.NewFromFloat(50.0)) {
        return InterestRate{}, errors.New("interest rate cannot exceed 50%")
    }
    return InterestRate{rate: rate}, nil
}

func (ir InterestRate) CalculateMonthlyInterest(principal Money, termMonths int) Money {
    monthlyRate := ir.rate.Div(decimal.NewFromInt(100)).Div(decimal.NewFromInt(12))
    interest := principal.amount.Mul(monthlyRate).Mul(decimal.NewFromInt(int64(termMonths)))
    result, _ := NewMoney(interest, principal.currency)
    return result
}

func (ir InterestRate) Value() decimal.Decimal { return ir.rate }
func (ir InterestRate) AsMonthlyRate() decimal.Decimal {
    return ir.rate.Div(decimal.NewFromInt(100)).Div(decimal.NewFromInt(12))
}

// Strongly-typed IDs prevent mixing different entity types
type LoanID struct {
    id uint64
}

func NewLoanID(id uint64) LoanID {
    if id == 0 {
        panic("LoanID cannot be zero")
    }
    return LoanID{id: id}
}

func (l LoanID) Value() uint64 { return l.id }
func (l LoanID) String() string { return fmt.Sprintf("loan-%d", l.id) }

type UserID struct {
    id uint64
}

func NewUserID(id uint64) UserID {
    if id == 0 {
        panic("UserID cannot be zero")
    }
    return UserID{id: id}
}

func (u UserID) Value() uint64 { return u.id }
func (u UserID) String() string { return fmt.Sprintf("user-%d", u.id) }

// Credit Score value object with business validation
type CreditScore struct {
    score int
}

func NewCreditScore(score int) (CreditScore, error) {
    if score < 300 || score > 850 {
        return CreditScore{}, errors.New("credit score must be between 300 and 850")
    }
    return CreditScore{score: score}, nil
}

func (cs CreditScore) Value() int { return cs.score }
func (cs CreditScore) IsExcellent() bool { return cs.score >= 750 }
func (cs CreditScore) IsGood() bool { return cs.score >= 650 }
func (cs CreditScore) IsFair() bool { return cs.score >= 600 }
func (cs CreditScore) IsPoor() bool { return cs.score < 600 }

func (cs CreditScore) RiskCategory() RiskCategory {
    if cs.IsExcellent() {
        return LowRisk
    } else if cs.IsGood() {
        return MediumRisk
    }
    return HighRisk
}
```

---

## Clean Architecture & Hexagonal Patterns

### 🏗️ Enhanced Layer Architecture

#### Target Layer Design
```
┌─────────────────────────────────────────────────────────┐
│                    Frameworks & Drivers                 │
│  HTTP Handlers • Database • External APIs • File System │
│                        (Infrastructure)                 │
├─────────────────────────────────────────────────────────┤
│                   Interface Adapters                    │
│     Controllers • Gateways • Presenters • Repositories  │
│                      (Infrastructure)                   │
├─────────────────────────────────────────────────────────┤
│                      Application                        │
│   Use Cases • Commands • Queries • Application Services │
│                      (Application)                      │
├─────────────────────────────────────────────────────────┤
│                       Domain                            │
│  Entities • Value Objects • Domain Services • Events    │
│                       (Domain)                          │
└─────────────────────────────────────────────────────────┘
```

#### Enhanced Hexagonal Design
```
                    ┌─ HTTP Adapter
                    │
                    ▼
            ┌───────────────┐
            │               │
   Database │               │ External APIs
   Adapter  │   Application │ Adapter
       ────▶│     Core      │◀────
            │               │
            │               │ Event Bus
            └───────────────┘ Adapter
                    │
                    ▼
               File System
                Adapter
```

### 🔧 Port and Adapter Implementation

#### Domain Ports (Interfaces)
```go
package domain

// Repository ports (secondary ports - driven)
type LoanRepository interface {
    Save(loan *Loan) error
    FindByID(id LoanID) (*Loan, error)
    FindByBorrower(borrowerID UserID) ([]*Loan, error)
    FindByStatus(status LoanStatus) ([]*Loan, error)
    FindRequiringApproval() ([]*Loan, error)
}

type UserRepository interface {
    FindByID(id UserID) (*User, error)
    Save(user *User) error
}

// External service ports (secondary ports - driven)
type CreditBureauPort interface {
    GetCreditScore(userID UserID) (CreditScore, error)
    GetCreditHistory(userID UserID) (CreditHistory, error)
}

type FraudDetectionPort interface {
    CheckFraud(userID UserID) (FraudRisk, error)
    ReportSuspiciousActivity(userID UserID, activity SuspiciousActivity) error
}

type PaymentProcessorPort interface {
    ProcessPayment(payment Payment) error
    RefundPayment(paymentID PaymentID) error
}

// Event publishing port (secondary port - driven)
type DomainEventPublisher interface {
    Publish(events []DomainEvent) error
    PublishAsync(events []DomainEvent) error
}

// Notification port (secondary port - driven)  
type NotificationPort interface {
    SendEmail(to EmailAddress, subject string, body string) error
    SendSMS(to PhoneNumber, message string) error
    SendPushNotification(userID UserID, message string) error
}
```

#### Application Layer (Use Cases)
```go
package application

// Command handlers (primary ports - driving)
type CreateLoanCommandHandler struct {
    loanRepo           domain.LoanRepository
    userRepo           domain.UserRepository
    riskService        *domain.RiskAssessmentService
    interestService    *domain.InterestCalculationService
    eventPublisher     domain.DomainEventPublisher
    notificationService domain.NotificationPort
    logger             Logger
}

func NewCreateLoanCommandHandler(
    loanRepo domain.LoanRepository,
    userRepo domain.UserRepository,
    riskService *domain.RiskAssessmentService,
    interestService *domain.InterestCalculationService,
    eventPublisher domain.DomainEventPublisher,
    notificationService domain.NotificationPort,
    logger Logger,
) *CreateLoanCommandHandler {
    return &CreateLoanCommandHandler{
        loanRepo:            loanRepo,
        userRepo:            userRepo,
        riskService:         riskService,
        interestService:     interestService,
        eventPublisher:      eventPublisher,
        notificationService: notificationService,
        logger:              logger,
    }
}

func (h *CreateLoanCommandHandler) Handle(ctx context.Context, cmd CreateLoanCommand) error {
    h.logger.Info("Processing create loan command", 
        "borrowerId", cmd.BorrowerID,
        "amount", cmd.Amount.String(),
    )
    
    // 1. Load domain objects
    borrower, err := h.userRepo.FindByID(cmd.BorrowerID)
    if err != nil {
        h.logger.Error("Failed to find borrower", "error", err)
        return fmt.Errorf("borrower not found: %w", err)
    }
    
    // 2. Risk assessment using domain service
    riskScore, err := h.riskService.AssessLoanRisk(borrower, cmd.Amount)
    if err != nil {
        h.logger.Error("Risk assessment failed", "error", err)
        return fmt.Errorf("risk assessment failed: %w", err)
    }
    
    if riskScore == domain.HighRisk {
        h.logger.Warn("Loan rejected due to high risk", 
            "borrowerId", cmd.BorrowerID,
            "riskScore", riskScore,
        )
        return NewDomainError("loan application rejected due to high risk")
    }
    
    // 3. Calculate interest rate using domain service
    interestRate, err := h.interestService.CalculateRate(riskScore, cmd.Amount, cmd.TermMonths)
    if err != nil {
        h.logger.Error("Interest rate calculation failed", "error", err)
        return fmt.Errorf("interest rate calculation failed: %w", err)
    }
    
    // 4. Create loan aggregate using domain logic
    loan, err := domain.NewLoanProposal(cmd.LoanID, cmd.BorrowerID, cmd.Amount, interestRate)
    if err != nil {
        h.logger.Error("Failed to create loan proposal", "error", err)
        return fmt.Errorf("failed to create loan proposal: %w", err)
    }
    
    // 5. Save aggregate (includes event publishing)
    if err := h.loanRepo.Save(loan); err != nil {
        h.logger.Error("Failed to save loan", "error", err)
        return fmt.Errorf("failed to save loan: %w", err)
    }
    
    // 6. Publish domain events
    if err := h.eventPublisher.Publish(loan.Events()); err != nil {
        h.logger.Error("Failed to publish events", "error", err)
        // Don't fail the operation, but log the error
    }
    
    // 7. Send notification
    if err := h.notificationService.SendEmail(
        borrower.Email(),
        "Loan Application Received",
        fmt.Sprintf("Your loan application for %s has been received and is under review.", cmd.Amount.String()),
    ); err != nil {
        h.logger.Error("Failed to send notification", "error", err)
        // Don't fail the operation
    }
    
    h.logger.Info("Loan created successfully", 
        "loanId", loan.ID(),
        "interestRate", interestRate.Value(),
    )
    
    return nil
}

// Query handlers (primary ports - driving)
type LoanQueriesHandler struct {
    loanRepo domain.LoanRepository
    logger   Logger
}

func (h *LoanQueriesHandler) GetLoanByID(ctx context.Context, loanID domain.LoanID) (*LoanDTO, error) {
    loan, err := h.loanRepo.FindByID(loanID)
    if err != nil {
        h.logger.Error("Failed to find loan", "loanId", loanID, "error", err)
        return nil, fmt.Errorf("loan not found: %w", err)
    }
    
    return h.mapLoanToDTO(loan), nil
}

func (h *LoanQueriesHandler) GetLoansByBorrower(ctx context.Context, borrowerID domain.UserID) ([]*LoanDTO, error) {
    loans, err := h.loanRepo.FindByBorrower(borrowerID)
    if err != nil {
        h.logger.Error("Failed to find loans by borrower", "borrowerId", borrowerID, "error", err)
        return nil, fmt.Errorf("loans not found: %w", err)
    }
    
    dtos := make([]*LoanDTO, len(loans))
    for i, loan := range loans {
        dtos[i] = h.mapLoanToDTO(loan)
    }
    
    return dtos, nil
}

func (h *LoanQueriesHandler) mapLoanToDTO(loan *domain.Loan) *LoanDTO {
    return &LoanDTO{
        ID:              loan.ID().Value(),
        BorrowerID:      loan.BorrowerID().Value(),
        PrincipalAmount: loan.PrincipalAmount().Amount(),
        Currency:        string(loan.PrincipalAmount().Currency()),
        InvestedAmount:  loan.InvestedAmount().Amount(),
        InterestRate:    loan.InterestRate().Value(),
        Status:          string(loan.Status()),
        CreatedAt:       loan.CreatedAt(),
    }
}
```

#### Infrastructure Adapters
```go
package infrastructure

// SQL Repository Adapter (secondary adapter)
type SQLLoanRepository struct {
    db             *sql.DB
    eventPublisher domain.DomainEventPublisher
    logger         Logger
}

func NewSQLLoanRepository(db *sql.DB, eventPublisher domain.DomainEventPublisher, logger Logger) domain.LoanRepository {
    return &SQLLoanRepository{
        db:             db,
        eventPublisher: eventPublisher,
        logger:         logger,
    }
}

func (r *SQLLoanRepository) Save(loan *domain.Loan) error {
    tx, err := r.db.Begin()
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback()
    
    // Map domain aggregate to SQL entities
    loanEntity := r.mapDomainToSQLEntity(loan)
    
    // Save aggregate root
    if err := r.saveLoanEntity(tx, loanEntity); err != nil {
        return fmt.Errorf("failed to save loan entity: %w", err)
    }
    
    // Save child entities (investments)
    for _, investment := range loan.Investments() {
        investmentEntity := r.mapInvestmentToSQLEntity(investment, loan.ID())
        if err := r.saveInvestmentEntity(tx, investmentEntity); err != nil {
            return fmt.Errorf("failed to save investment entity: %w", err)
        }
    }
    
    // Commit transaction
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    // Publish domain events (after successful persistence)
    if len(loan.Events()) > 0 {
        if err := r.eventPublisher.Publish(loan.Events()); err != nil {
            r.logger.Error("Failed to publish domain events", "error", err)
            // Don't fail the save operation
        }
        loan.ClearEvents()
    }
    
    return nil
}

func (r *SQLLoanRepository) FindByID(id domain.LoanID) (*domain.Loan, error) {
    // Query aggregate root
    loanEntity, err := r.findLoanEntityByID(id)
    if err != nil {
        return nil, fmt.Errorf("failed to find loan: %w", err)
    }
    
    // Query child entities
    investmentEntities, err := r.findInvestmentEntitiesByLoanID(id)
    if err != nil {
        return nil, fmt.Errorf("failed to find investments: %w", err)
    }
    
    // Reconstruct domain aggregate
    loan, err := r.mapSQLToDomainAggregate(loanEntity, investmentEntities)
    if err != nil {
        return nil, fmt.Errorf("failed to reconstruct loan aggregate: %w", err)
    }
    
    return loan, nil
}

func (r *SQLLoanRepository) mapDomainToSQLEntity(loan *domain.Loan) *LoanSQLEntity {
    return &LoanSQLEntity{
        ID:                         loan.ID().Value(),
        BorrowerID:                 loan.BorrowerID().Value(),
        PrincipalAmount:            loan.PrincipalAmount().Amount(),
        Currency:                   string(loan.PrincipalAmount().Currency()),
        InvestedAmount:             loan.InvestedAmount().Amount(),
        InterestRate:               loan.InterestRate().Value(),
        Status:                     string(loan.Status()),
        ApprovalDate:               r.timeToNullTime(loan.ApprovalDate()),
        ApprovalEmployeeID:         r.userIDToNullInt64(loan.ApproverID()),
        DisbursementDate:           r.timeToNullTime(loan.DisbursementDate()),
        AgreementLetterDocumentURL: r.stringToNullString(loan.AgreementURL()),
        CreatedAt:                  loan.CreatedAt(),
        UpdatedAt:                  time.Now(),
    }
}

// HTTP Adapter (primary adapter)
type LoanHTTPHandler struct {
    createLoanHandler *application.CreateLoanCommandHandler
    loanQueries      *application.LoanQueriesHandler
    logger           Logger
}

func NewLoanHTTPHandler(
    createLoanHandler *application.CreateLoanCommandHandler,
    loanQueries *application.LoanQueriesHandler,
    logger Logger,
) *LoanHTTPHandler {
    return &LoanHTTPHandler{
        createLoanHandler: createLoanHandler,
        loanQueries:      loanQueries,
        logger:           logger,
    }
}

func (h *LoanHTTPHandler) CreateLoan(w http.ResponseWriter, r *http.Request) {
    // 1. Parse HTTP request
    var request CreateLoanRequest
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        h.logger.Error("Failed to decode request", "error", err)
        http.Error(w, "Invalid request format", http.StatusBadRequest)
        return
    }
    
    // 2. Validate request
    if err := h.validateCreateLoanRequest(request); err != nil {
        h.logger.Error("Request validation failed", "error", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // 3. Convert to domain command
    command, err := h.mapRequestToCommand(request)
    if err != nil {
        h.logger.Error("Failed to map request to command", "error", err)
        http.Error(w, "Invalid request data", http.StatusBadRequest)
        return
    }
    
    // 4. Execute use case
    ctx := r.Context()
    if err := h.createLoanHandler.Handle(ctx, command); err != nil {
        h.logger.Error("Failed to create loan", "error", err)
        
        // Map domain errors to HTTP status codes
        if isDomainError(err) {
            http.Error(w, err.Error(), http.StatusBadRequest)
        } else {
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }
    
    // 5. Return HTTP response
    response := CreateLoanResponse{
        LoanID: command.LoanID.Value(),
        Status: "created",
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(response)
}

func (h *LoanHTTPHandler) GetLoan(w http.ResponseWriter, r *http.Request) {
    // Extract loan ID from URL
    vars := mux.Vars(r)
    loanIDStr := vars["loanId"]
    
    loanIDValue, err := strconv.ParseUint(loanIDStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid loan ID", http.StatusBadRequest)
        return
    }
    
    loanID := domain.NewLoanID(loanIDValue)
    
    // Execute query
    ctx := r.Context()
    loanDTO, err := h.loanQueries.GetLoanByID(ctx, loanID)
    if err != nil {
        if isNotFoundError(err) {
            http.Error(w, "Loan not found", http.StatusNotFound)
        } else {
            h.logger.Error("Failed to get loan", "error", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }
    
    // Return response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(loanDTO)
}

func (h *LoanHTTPHandler) mapRequestToCommand(request CreateLoanRequest) (application.CreateLoanCommand, error) {
    // Convert primitive types to value objects
    borrowerID := domain.NewUserID(request.BorrowerID)
    
    amount, err := domain.NewMoney(request.Amount, domain.Currency(request.Currency))
    if err != nil {
        return application.CreateLoanCommand{}, fmt.Errorf("invalid amount: %w", err)
    }
    
    loanID := domain.NewLoanID(generateSnowflakeID()) // Assume ID generation
    
    return application.CreateLoanCommand{
        LoanID:     loanID,
        BorrowerID: borrowerID,
        Amount:     amount,
        TermMonths: request.TermMonths,
    }, nil
}
```

---

## Target Architecture Design

### 🎯 Complete Target Structure

#### Enhanced Project Structure (Full DDD Implementation)
```
internal/
├── loan/                           # 📋 Bounded Context
│   ├── domain/                    # 🎯 Domain Layer (Pure Business Logic)
│   │   ├── entities/              # Aggregates & Rich Entities
│   │   │   ├── loan.go                   # Loan aggregate root with business methods
│   │   │   ├── investment.go             # Investment entity within loan aggregate
│   │   │   ├── borrower.go               # Borrower entity with credit profile
│   │   │   └── loan_term.go              # Loan term entity with payment schedule
│   │   ├── valueobjects/          # Value Objects with Business Rules
│   │   │   ├── money.go                  # Money with currency and arithmetic
│   │   │   ├── interest_rate.go          # InterestRate with validation & calculations
│   │   │   ├── credit_score.go           # CreditScore with risk categorization
│   │   │   ├── loan_id.go                # Strongly-typed LoanID
│   │   │   ├── user_id.go                # Strongly-typed UserID  
│   │   │   ├── email_address.go          # EmailAddress with validation
│   │   │   ├── phone_number.go           # PhoneNumber with formatting
│   │   │   └── currency.go               # Currency enumeration
│   │   ├── services/              # Domain Services for Complex Business Logic
│   │   │   ├── risk_assessment_service.go      # Credit risk evaluation
│   │   │   ├── interest_calculation_service.go # Dynamic interest rate calculation
│   │   │   ├── loan_matching_service.go        # Investor-loan matching algorithm
│   │   │   ├── payment_schedule_service.go     # Payment schedule generation
│   │   │   └── loan_validation_service.go      # Complex loan validation rules
│   │   ├── events/                # Domain Events for System Integration
│   │   │   ├── loan_events.go            # Loan lifecycle events
│   │   │   ├── investment_events.go      # Investment-related events
│   │   │   ├── payment_events.go         # Payment processing events
│   │   │   └── risk_events.go            # Risk assessment events
│   │   ├── specifications/        # Business Rules as Specifications Pattern
│   │   │   ├── loan_approval_spec.go     # Loan approval criteria
│   │   │   ├── investment_eligibility_spec.go # Investment rules
│   │   │   └── disbursement_requirements_spec.go # Disbursement validation
│   │   ├── repositories/          # Repository Interfaces (Ports)
│   │   │   ├── loan_repository.go        # Loan aggregate repository
│   │   │   ├── user_repository.go        # User repository
│   │   │   └── payment_repository.go     # Payment repository
│   │   └── errors/                # Domain-Specific Errors
│   │       ├── domain_errors.go          # Base domain error types
│   │       └── business_rule_errors.go   # Business rule violation errors
│   ├── application/               # 🔧 Application Layer (Use Case Orchestration)
│   │   ├── commands/              # Command Handlers (Write Operations - CQRS)
│   │   │   ├── create_loan_command.go           # Create new loan proposal
│   │   │   ├── approve_loan_command.go          # Approve loan application
│   │   │   ├── invest_loan_command.go           # Make investment in loan
│   │   │   ├── disburse_loan_command.go         # Disburse loan funds
│   │   │   ├── process_payment_command.go       # Process loan payment
│   │   │   └── update_credit_score_command.go   # Update borrower credit
│   │   ├── queries/               # Query Handlers (Read Operations - CQRS)
│   │   │   ├── loan_queries.go              # Loan query operations
│   │   │   ├── investment_queries.go        # Investment query operations
│   │   │   ├── borrower_queries.go          # Borrower query operations
│   │   │   ├── payment_queries.go           # Payment history queries
│   │   │   └── risk_assessment_queries.go   # Risk assessment queries
│   │   ├── services/              # Application Services (Cross-Aggregate Operations)
│   │   │   ├── loan_application_service.go      # Loan application orchestration
│   │   │   ├── investment_application_service.go # Investment orchestration
│   │   │   ├── payment_processing_service.go    # Payment processing workflows
│   │   │   └── notification_service.go          # User notification workflows
│   │   ├── dto/                   # Data Transfer Objects
│   │   │   ├── loan_dto.go              # Loan data transfer objects
│   │   │   ├── investment_dto.go        # Investment DTOs
│   │   │   ├── payment_dto.go           # Payment DTOs
│   │   │   └── user_dto.go              # User DTOs
│   │   └── workflows/             # Complex Business Workflows
│   │       ├── loan_approval_workflow.go    # Multi-step loan approval
│   │       ├── disbursement_workflow.go     # Fund disbursement process
│   │       └── collection_workflow.go       # Payment collection process
│   ├── infrastructure/            # 🔌 Infrastructure Layer (External Concerns)
│   │   ├── persistence/           # Repository Implementations
│   │   │   ├── sql/               # SQL Database Implementations
│   │   │   │   ├── loan_sql_repository.go      # SQL loan repository
│   │   │   │   ├── user_sql_repository.go      # SQL user repository
│   │   │   │   ├── payment_sql_repository.go   # SQL payment repository
│   │   │   │   ├── migrations/                 # Database migrations
│   │   │   │   │   ├── 001_create_loans.sql
│   │   │   │   │   ├── 002_create_investments.sql
│   │   │   │   │   └── 003_create_payments.sql
│   │   │   │   └── entities/               # SQL Mapping Entities
│   │   │   │       ├── loan_sql_entity.go     # SQL-specific loan entity
│   │   │   │       ├── investment_sql_entity.go # SQL investment entity
│   │   │   │       └── payment_sql_entity.go  # SQL payment entity
│   │   │   ├── redis/             # Redis Cache Implementations
│   │   │   │   ├── loan_cache_repository.go
│   │   │   │   └── user_session_repository.go
│   │   │   └── memory/            # In-Memory Implementations (Testing)
│   │   │       ├── loan_memory_repository.go
│   │   │       └── user_memory_repository.go
│   │   ├── http/                  # HTTP Adapters (Primary Adapters)
│   │   │   ├── handlers/          # HTTP Request Handlers
│   │   │   │   ├── loan_handler.go            # Loan API endpoints
│   │   │   │   ├── investment_handler.go      # Investment API endpoints
│   │   │   │   ├── payment_handler.go         # Payment API endpoints
│   │   │   │   └── health_handler.go          # Health check endpoints
│   │   │   ├── middleware/        # HTTP Middleware
│   │   │   │   ├── authentication_middleware.go # Auth middleware
│   │   │   │   ├── validation_middleware.go     # Request validation
│   │   │   │   ├── logging_middleware.go        # Request logging
│   │   │   │   ├── rate_limiting_middleware.go  # Rate limiting
│   │   │   │   └── cors_middleware.go           # CORS handling
│   │   │   ├── serializers/       # Request/Response Serialization
│   │   │   │   ├── loan_serializer.go          # Loan JSON serialization
│   │   │   │   └── error_serializer.go         # Error response serialization
│   │   │   └── routes.go          # Route definitions
│   │   ├── events/                # Event Infrastructure
│   │   │   ├── publishers/        # Event Publishers
│   │   │   │   ├── domain_event_publisher.go   # Domain event publisher
│   │   │   │   ├── kafka_event_publisher.go    # Kafka event publisher
│   │   │   │   └── rabbitmq_event_publisher.go # RabbitMQ publisher
│   │   │   ├── subscribers/       # Event Subscribers
│   │   │   │   ├── loan_event_subscriber.go    # Loan event handlers
│   │   │   │   ├── payment_event_subscriber.go # Payment event handlers
│   │   │   │   └── notification_event_subscriber.go # Notification handlers
│   │   │   └── serializers/       # Event Serialization
│   │   │       └── event_serializer.go
│   │   ├── external/              # External Service Adapters (Secondary Adapters)
│   │   │   ├── credit_bureau_adapter.go        # Credit bureau integration
│   │   │   ├── payment_processor_adapter.go    # Payment processor (Stripe, PayPal)
│   │   │   ├── fraud_detection_adapter.go      # Fraud detection service
│   │   │   ├── email_service_adapter.go        # Email service integration
│   │   │   ├── sms_service_adapter.go          # SMS service integration
│   │   │   └── document_storage_adapter.go     # Document storage (S3, etc.)
│   │   ├── monitoring/            # Monitoring & Observability
│   │   │   ├── metrics_collector.go           # Application metrics
│   │   │   ├── health_checker.go              # Health check implementation
│   │   │   └── tracing_middleware.go          # Distributed tracing
│   │   └── configuration/         # Configuration Management
│   │       ├── config_loader.go               # Configuration loading
│   │       └── environment_validator.go       # Environment validation
│   ├── test/                      # Integration & Acceptance Tests
│   │   ├── integration/           # Integration Tests
│   │   │   ├── loan_integration_test.go       # Loan flow integration tests
│   │   │   ├── payment_integration_test.go    # Payment flow tests
│   │   │   └── api_integration_test.go        # API endpoint tests
│   │   ├── acceptance/            # Acceptance Tests (BDD Style)
│   │   │   ├── loan_lifecycle_test.go         # End-to-end loan lifecycle
│   │   │   ├── investment_flow_test.go        # Investment flow testing
│   │   │   └── payment_processing_test.go     # Payment processing E2E
│   │   ├── fixtures/              # Test Data Fixtures
│   │   │   ├── loan_fixtures.go              # Loan test data
│   │   │   └── user_fixtures.go              # User test data
│   │   └── helpers/               # Test Helper Functions
│   │       ├── database_helper.go            # Database test utilities
│   │       └── http_helper.go                # HTTP test utilities
│   ├── docs/                      # Domain-Specific Documentation
│   │   ├── README.md              # Loan domain overview
│   │   ├── DOMAIN_MODEL.md        # Domain model documentation
│   │   ├── UBIQUITOUS_LANGUAGE.md # Business terminology glossary
│   │   ├── API.md                 # API documentation
│   │   └── adr/                   # Architecture Decision Records
│   │       ├── 001-aggregate-boundaries.md
│   │       ├── 002-event-sourcing-decision.md
│   │       ├── 003-cqrs-implementation.md
│   │       └── 004-payment-processing-strategy.md
│   └── module.go                  # Dependency Injection & Module Setup
├── shared/                        # 🌐 Shared Kernel (Cross-Domain)
│   ├── domain/                   # Common Domain Concepts
│   │   ├── events/               # Base Event Infrastructure
│   │   │   ├── domain_event.go          # Base domain event interface
│   │   │   ├── event_store.go           # Event store interface
│   │   │   └── event_dispatcher.go      # Event dispatcher interface
│   │   ├── errors/               # Common Error Types
│   │   │   ├── domain_error.go          # Base domain error
│   │   │   ├── business_error.go        # Business rule violations
│   │   │   ├── validation_error.go      # Validation errors
│   │   │   └── not_found_error.go       # Entity not found errors
│   │   ├── valueobjects/         # Shared Value Objects
│   │   │   ├── currency.go              # Currency enumeration
│   │   │   ├── timestamp.go             # Business timestamp
│   │   │   ├── percentage.go            # Percentage calculations
│   │   │   └── address.go               # Physical address
│   │   ├── specifications/       # Common Specification Pattern
│   │   │   ├── specification.go         # Base specification interface
│   │   │   └── composite_specification.go # Composite specifications
│   │   └── policies/             # Cross-Domain Business Policies
│   │       ├── rate_limiting_policy.go  # API rate limiting rules
│   │       └── data_retention_policy.go # Data retention rules
│   └── infrastructure/          # Shared Infrastructure Components
│       ├── persistence/         # Common Persistence Utilities
│       │   ├── transaction_manager.go   # Database transaction management
│       │   ├── connection_manager.go    # Database connection pooling
│       │   ├── migration_runner.go      # Database migration runner
│       │   └── query_builder.go         # Common query building utilities
│       ├── events/              # Shared Event Infrastructure
│       │   ├── event_bus.go            # In-process event bus
│       │   ├── event_store.go          # Event store implementation
│       │   ├── event_dispatcher.go     # Event routing and dispatch
│       │   └── saga_orchestrator.go    # Saga pattern implementation
│       ├── messaging/           # Message Queue Infrastructure
│       │   ├── message_publisher.go    # Message publishing interface
│       │   ├── message_consumer.go     # Message consumption interface
│       │   ├── kafka_client.go         # Kafka client implementation
│       │   └── rabbitmq_client.go      # RabbitMQ client implementation
│       ├── logging/             # Structured Logging
│       │   ├── structured_logger.go    # Logger interface and implementation
│       │   ├── log_formatter.go        # Log formatting utilities
│       │   └── correlation_id.go       # Request correlation tracking
│       ├── monitoring/          # Monitoring & Metrics
│       │   ├── metrics_collector.go    # Application metrics collection
│       │   ├── health_monitor.go       # Health monitoring
│       │   └── performance_monitor.go  # Performance monitoring
│       ├── security/            # Security Utilities
│       │   ├── encryption_service.go   # Data encryption/decryption
│       │   ├── token_validator.go      # JWT token validation
│       │   └── audit_logger.go         # Security audit logging
│       └── validation/          # Common Validation
│           ├── validator.go            # Validation utilities
│           ├── business_rules.go       # Cross-domain business rules
│           └── constraint_validator.go # Constraint validation
├── user/                         # 🏠 User Management Bounded Context (Future)
│   ├── domain/                  # User domain logic
│   ├── application/             # User use cases
│   ├── infrastructure/          # User infrastructure
│   └── module.go                # User module setup
├── payment/                      # 💳 Payment Processing Bounded Context (Future)
│   ├── domain/                  # Payment domain logic
│   ├── application/             # Payment use cases  
│   ├── infrastructure/          # Payment infrastructure
│   └── module.go                # Payment module setup
├── notification/                 # 📧 Notification Bounded Context (Future)
│   ├── domain/                  # Notification domain logic
│   ├── application/             # Notification use cases
│   ├── infrastructure/          # Notification infrastructure
│   └── module.go                # Notification module setup
└── app/                         # 🚀 Application Bootstrap & Configuration
    ├── bootstrap/               # Application Startup
    │   ├── app.go                      # Main application bootstrap
    │   ├── dependency_injection.go     # DI container setup
    │   ├── module_registry.go          # Module registration
    │   └── graceful_shutdown.go        # Graceful shutdown handling
    ├── config/                  # Configuration Management
    │   ├── config.go                   # Configuration structures
    │   ├── environment.go              # Environment-specific configs
    │   ├── database_config.go          # Database configuration
    │   └── service_config.go           # External service configuration
    ├── server/                  # Server Implementation
    │   ├── http_server.go              # HTTP server setup
    │   ├── grpc_server.go              # gRPC server setup (future)
    │   └── websocket_server.go         # WebSocket server (future)
    └── scripts/                 # Deployment & Utility Scripts
        ├── migrate.go                  # Database migration script
        ├── seed.go                     # Database seeding script
        └── health_check.go             # Health check script
```

### 🎯 Domain Services Architecture

#### Risk Assessment Service (Complex Business Logic)
```go
package domain

import (
    "context"
    "fmt"
)

// Risk Assessment Service - Domain Service for Complex Business Logic
type RiskAssessmentService struct {
    creditBureauPort   CreditBureauPort
    fraudDetectionPort FraudDetectionPort
    logger            Logger
}

func NewRiskAssessmentService(
    creditBureauPort CreditBureauPort,
    fraudDetectionPort FraudDetectionPort,
    logger Logger,
) *RiskAssessmentService {
    return &RiskAssessmentService{
        creditBureauPort:   creditBureauPort,
        fraudDetectionPort: fraudDetectionPort,
        logger:            logger,
    }
}

type RiskAssessmentResult struct {
    OverallRisk          RiskCategory
    CreditScore          CreditScore
    FraudRisk           FraudRiskLevel
    DebtToIncomeRatio   decimal.Decimal
    RecommendedRate     InterestRate
    MaxLoanAmount       Money
    Reasoning          []string
}

type RiskCategory int
const (
    LowRisk RiskCategory = iota + 1
    MediumRisk
    HighRisk
    UnacceptableRisk
)

func (rs RiskCategory) String() string {
    return [...]string{"", "Low", "Medium", "High", "Unacceptable"}[rs]
}

type FraudRiskLevel int
const (
    NoFraudRisk FraudRiskLevel = iota
    LowFraudRisk
    MediumFraudRisk
    HighFraudRisk
)

func (ras *RiskAssessmentService) AssessLoanRisk(
    ctx context.Context,
    borrower *User,
    loanAmount Money,
    termMonths int,
) (*RiskAssessmentResult, error) {
    
    ras.logger.Info("Starting risk assessment", 
        "userId", borrower.ID(),
        "amount", loanAmount.String(),
        "term", termMonths,
    )
    
    result := &RiskAssessmentResult{
        Reasoning: make([]string, 0),
    }
    
    // 1. Get credit score from external credit bureau
    creditScore, err := ras.creditBureauPort.GetCreditScore(ctx, borrower.ID())
    if err != nil {
        ras.logger.Error("Failed to get credit score", "error", err)
        return nil, fmt.Errorf("credit assessment failed: %w", err)
    }
    result.CreditScore = creditScore
    
    // 2. Perform fraud detection check
    fraudRisk, err := ras.fraudDetectionPort.AssessFraudRisk(ctx, borrower.ID())
    if err != nil {
        ras.logger.Error("Failed to assess fraud risk", "error", err)
        // Don't fail the assessment, but log the issue
        fraudRisk = MediumFraudRisk
        result.Reasoning = append(result.Reasoning, "Could not complete fraud assessment")
    }
    result.FraudRisk = fraudRisk
    
    // 3. Calculate debt-to-income ratio
    debtToIncomeRatio := ras.calculateDebtToIncomeRatio(borrower, loanAmount, termMonths)
    result.DebtToIncomeRatio = debtToIncomeRatio
    
    // 4. Determine overall risk category using business rules
    overallRisk := ras.determineOverallRisk(creditScore, fraudRisk, debtToIncomeRatio, loanAmount)
    result.OverallRisk = overallRisk
    
    // 5. Calculate recommended interest rate based on risk
    recommendedRate, err := ras.calculateRecommendedRate(overallRisk, creditScore, loanAmount, termMonths)
    if err != nil {
        return nil, fmt.Errorf("rate calculation failed: %w", err)
    }
    result.RecommendedRate = recommendedRate
    
    // 6. Determine maximum loan amount for this borrower
    maxLoanAmount := ras.calculateMaxLoanAmount(creditScore, borrower.Income(), debtToIncomeRatio)
    result.MaxLoanAmount = maxLoanAmount
    
    // 7. Add business reasoning
    result.Reasoning = append(result.Reasoning, ras.generateReasoning(result)...)
    
    ras.logger.Info("Risk assessment completed",
        "overallRisk", overallRisk.String(),
        "recommendedRate", recommendedRate.Value(),
        "maxAmount", maxLoanAmount.String(),
    )
    
    return result, nil
}

func (ras *RiskAssessmentService) determineOverallRisk(
    creditScore CreditScore,
    fraudRisk FraudRiskLevel,
    debtToIncomeRatio decimal.Decimal,
    loanAmount Money,
) RiskCategory {
    
    // Business rules for risk determination
    
    // Immediate disqualifiers
    if fraudRisk == HighFraudRisk {
        return UnacceptableRisk
    }
    
    if debtToIncomeRatio.GreaterThan(decimal.NewFromFloat(0.5)) { // >50% DTI
        return UnacceptableRisk
    }
    
    // Excellent credit with low risk factors
    if creditScore.IsExcellent() && fraudRisk == NoFraudRisk && 
       debtToIncomeRatio.LessThan(decimal.NewFromFloat(0.3)) &&
       loanAmount.Amount().LessThan(decimal.NewFromInt(100000)) {
        return LowRisk
    }
    
    // Good credit with reasonable risk factors
    if creditScore.IsGood() && fraudRisk <= LowFraudRisk &&
       debtToIncomeRatio.LessThan(decimal.NewFromFloat(0.4)) {
        return MediumRisk
    }
    
    // Fair credit or higher risk factors
    if creditScore.IsFair() || fraudRisk == MediumFraudRisk ||
       debtToIncomeRatio.GreaterThan(decimal.NewFromFloat(0.4)) {
        return HighRisk
    }
    
    // Default to high risk for poor credit
    return HighRisk
}

func (ras *RiskAssessmentService) calculateRecommendedRate(
    riskCategory RiskCategory,
    creditScore CreditScore,
    loanAmount Money,
    termMonths int,
) (InterestRate, error) {
    
    baseRate := decimal.NewFromFloat(8.0) // 8% base rate
    
    // Risk-based adjustments
    switch riskCategory {
    case LowRisk:
        baseRate = baseRate.Sub(decimal.NewFromFloat(2.0)) // -2% for low risk
    case MediumRisk:
        // No adjustment
    case HighRisk:
        baseRate = baseRate.Add(decimal.NewFromFloat(3.0)) // +3% for high risk
    case UnacceptableRisk:
        return InterestRate{}, NewDomainError("loan cannot be approved due to unacceptable risk")
    }
    
    // Credit score fine-tuning
    if creditScore.IsExcellent() {
        baseRate = baseRate.Sub(decimal.NewFromFloat(0.5)) // Additional -0.5% for excellent credit
    } else if creditScore.IsPoor() {
        baseRate = baseRate.Add(decimal.NewFromFloat(1.0)) // Additional +1% for poor credit
    }
    
    // Loan amount adjustments (smaller loans = higher rates due to fixed costs)
    if loanAmount.Amount().LessThan(decimal.NewFromInt(10000)) {
        baseRate = baseRate.Add(decimal.NewFromFloat(1.5)) // +1.5% for small loans
    } else if loanAmount.Amount().GreaterThan(decimal.NewFromInt(100000)) {
        baseRate = baseRate.Sub(decimal.NewFromFloat(0.5)) // -0.5% for large loans
    }
    
    // Term adjustments (longer terms = higher rates)
    if termMonths > 60 {
        baseRate = baseRate.Add(decimal.NewFromFloat(1.0)) // +1% for long terms
    } else if termMonths <= 12 {
        baseRate = baseRate.Sub(decimal.NewFromFloat(0.5)) // -0.5% for short terms
    }
    
    // Ensure rate stays within acceptable bounds
    minRate := decimal.NewFromFloat(5.0)
    maxRate := decimal.NewFromFloat(35.0)
    
    if baseRate.LessThan(minRate) {
        baseRate = minRate
    }
    if baseRate.GreaterThan(maxRate) {
        baseRate = maxRate
    }
    
    return NewInterestRate(baseRate)
}

func (ras *RiskAssessmentService) calculateDebtToIncomeRatio(
    borrower *User,
    loanAmount Money,
    termMonths int,
) decimal.Decimal {
    
    monthlyIncome := borrower.MonthlyIncome()
    if monthlyIncome.IsZero() {
        return decimal.NewFromFloat(1.0) // 100% DTI if no income reported
    }
    
    // Calculate estimated monthly payment (simplified calculation)
    interestRate := decimal.NewFromFloat(10.0) // Assume 10% for DTI calculation
    monthlyRate := interestRate.Div(decimal.NewFromFloat(100)).Div(decimal.NewFromFloat(12))
    
    // Monthly payment formula: P * [r(1+r)^n] / [(1+r)^n - 1]
    onePlusR := decimal.NewFromFloat(1).Add(monthlyRate)
    onePlusRPowN := onePlusR.Pow(decimal.NewFromInt(int64(termMonths)))
    
    monthlyPayment := loanAmount.Amount().Mul(
        monthlyRate.Mul(onePlusRPowN).Div(onePlusRPowN.Sub(decimal.NewFromFloat(1))),
    )
    
    existingDebt := borrower.MonthlyDebtPayments()
    totalMonthlyDebt := existingDebt.Add(monthlyPayment)
    
    return totalMonthlyDebt.Div(monthlyIncome.Amount())
}

func (ras *RiskAssessmentService) calculateMaxLoanAmount(
    creditScore CreditScore,
    income Money,
    currentDTI decimal.Decimal,
) Money {
    
    maxDTI := decimal.NewFromFloat(0.43) // 43% maximum DTI
    
    if currentDTI.GreaterThanOrEqual(maxDTI) {
        return NewZeroMoney(income.Currency()) // No additional borrowing capacity
    }
    
    availableDTI := maxDTI.Sub(currentDTI)
    availablePayment := income.Amount().Mul(availableDTI)
    
    // Convert available payment to loan amount (simplified)
    // Assume 10% rate, 36 month term for calculation
    estimatedLoanAmount := availablePayment.Mul(decimal.NewFromInt(30)) // Rough multiplier
    
    // Apply credit score limits
    var maxAmount decimal.Decimal
    switch {
    case creditScore.IsExcellent():
        maxAmount = decimal.NewFromInt(250000)
    case creditScore.IsGood():
        maxAmount = decimal.NewFromInt(150000)
    case creditScore.IsFair():
        maxAmount = decimal.NewFromInt(75000)
    default:
        maxAmount = decimal.NewFromInt(25000)
    }
    
    if estimatedLoanAmount.GreaterThan(maxAmount) {
        estimatedLoanAmount = maxAmount
    }
    
    result, _ := NewMoney(estimatedLoanAmount, income.Currency())
    return result
}

func (ras *RiskAssessmentService) generateReasoning(result *RiskAssessmentResult) []string {
    reasoning := make([]string, 0)
    
    // Credit score reasoning
    if result.CreditScore.IsExcellent() {
        reasoning = append(reasoning, "Excellent credit score indicates strong payment history")
    } else if result.CreditScore.IsPoor() {
        reasoning = append(reasoning, "Poor credit score indicates payment difficulties")
    }
    
    // Fraud risk reasoning
    if result.FraudRisk == HighFraudRisk {
        reasoning = append(reasoning, "High fraud risk detected - loan denied")
    } else if result.FraudRisk == NoFraudRisk {
        reasoning = append(reasoning, "No fraud indicators detected")
    }
    
    // DTI reasoning
    if result.DebtToIncomeRatio.GreaterThan(decimal.NewFromFloat(0.4)) {
        reasoning = append(reasoning, "High debt-to-income ratio indicates potential payment stress")
    } else if result.DebtToIncomeRatio.LessThan(decimal.NewFromFloat(0.3)) {
        reasoning = append(reasoning, "Low debt-to-income ratio indicates good payment capacity")
    }
    
    // Overall risk reasoning
    switch result.OverallRisk {
    case LowRisk:
        reasoning = append(reasoning, "Low risk profile qualifies for preferential rates")
    case HighRisk:
        reasoning = append(reasoning, "High risk profile requires elevated interest rate")
    case UnacceptableRisk:
        reasoning = append(reasoning, "Risk profile exceeds lending criteria")
    }
    
    return reasoning
}
```

#### Interest Rate Calculation Service
```go
package domain

// Interest Rate Calculation Service - Sophisticated Rate Calculation
type InterestCalculationService struct {
    marketDataPort MarketDataPort
    logger        Logger
}

func NewInterestCalculationService(
    marketDataPort MarketDataPort,
    logger Logger,
) *InterestCalculationService {
    return &InterestCalculationService{
        marketDataPort: marketDataPort,
        logger:        logger,
    }
}

type RateCalculationInput struct {
    RiskCategory    RiskCategory
    CreditScore     CreditScore
    LoanAmount      Money
    TermMonths      int
    LoanPurpose     LoanPurpose
    CollateralType  CollateralType
    MarketCondition MarketCondition
}

type LoanPurpose int
const (
    PersonalLoan LoanPurpose = iota + 1
    HomeImprovement
    DebtConsolidation
    BusinessLoan
    Education
)

type CollateralType int
const (
    Unsecured CollateralType = iota
    SecuredByProperty
    SecuredByVehicle
    SecuredByInvestment
)

func (ics *InterestCalculationService) CalculateOptimalRate(
    ctx context.Context,
    input RateCalculationInput,
) (InterestRate, error) {
    
    ics.logger.Info("Calculating optimal interest rate",
        "riskCategory", input.RiskCategory,
        "creditScore", input.CreditScore.Value(),
        "amount", input.LoanAmount.String(),
    )
    
    // 1. Get market base rate
    marketRate, err := ics.marketDataPort.GetCurrentBaseRate(ctx)
    if err != nil {
        ics.logger.Error("Failed to get market rate", "error", err)
        // Fallback to default rate
        marketRate, _ = NewInterestRate(decimal.NewFromFloat(6.0))
    }
    
    baseRate := marketRate.Value()
    
    // 2. Apply risk-based margin
    riskMargin := ics.calculateRiskMargin(input.RiskCategory, input.CreditScore)
    adjustedRate := baseRate.Add(riskMargin)
    
    // 3. Apply loan amount adjustments
    amountAdjustment := ics.calculateAmountAdjustment(input.LoanAmount)
    adjustedRate = adjustedRate.Add(amountAdjustment)
    
    // 4. Apply term-based adjustments
    termAdjustment := ics.calculateTermAdjustment(input.TermMonths)
    adjustedRate = adjustedRate.Add(termAdjustment)
    
    // 5. Apply purpose-based adjustments
    purposeAdjustment := ics.calculatePurposeAdjustment(input.LoanPurpose)
    adjustedRate = adjustedRate.Add(purposeAdjustment)
    
    // 6. Apply collateral adjustments
    collateralAdjustment := ics.calculateCollateralAdjustment(input.CollateralType)
    adjustedRate = adjustedRate.Add(collateralAdjustment)
    
    // 7. Apply market condition adjustments
    marketAdjustment := ics.calculateMarketAdjustment(input.MarketCondition)
    adjustedRate = adjustedRate.Add(marketAdjustment)
    
    // 8. Ensure rate is within acceptable bounds
    finalRate := ics.applyRateBounds(adjustedRate)
    
    result, err := NewInterestRate(finalRate)
    if err != nil {
        return InterestRate{}, fmt.Errorf("invalid calculated rate: %w", err)
    }
    
    ics.logger.Info("Interest rate calculated",
        "baseRate", baseRate,
        "finalRate", finalRate,
        "riskMargin", riskMargin,
    )
    
    return result, nil
}

func (ics *InterestCalculationService) calculateRiskMargin(
    riskCategory RiskCategory,
    creditScore CreditScore,
) decimal.Decimal {
    
    var baseMargin decimal.Decimal
    
    // Base margin by risk category
    switch riskCategory {
    case LowRisk:
        baseMargin = decimal.NewFromFloat(-1.0) // -1% for low risk
    case MediumRisk:
        baseMargin = decimal.NewFromFloat(1.0) // +1% for medium risk
    case HighRisk:
        baseMargin = decimal.NewFromFloat(4.0) // +4% for high risk
    default:
        baseMargin = decimal.NewFromFloat(8.0) // +8% for unacceptable risk
    }
    
    // Fine-tune based on credit score
    if creditScore.IsExcellent() {
        baseMargin = baseMargin.Sub(decimal.NewFromFloat(0.5)) // Additional -0.5%
    } else if creditScore.IsPoor() {
        baseMargin = baseMargin.Add(decimal.NewFromFloat(1.5)) // Additional +1.5%
    }
    
    return baseMargin
}

func (ics *InterestCalculationService) calculateAmountAdjustment(amount Money) decimal.Decimal {
    amountValue := amount.Amount()
    
    if amountValue.LessThan(decimal.NewFromInt(5000)) {
        return decimal.NewFromFloat(2.0) // +2% for very small loans
    } else if amountValue.LessThan(decimal.NewFromInt(25000)) {
        return decimal.NewFromFloat(0.5) // +0.5% for small loans
    } else if amountValue.GreaterThan(decimal.NewFromInt(100000)) {
        return decimal.NewFromFloat(-0.5) // -0.5% for large loans
    }
    
    return decimal.Zero // No adjustment for medium amounts
}

func (ics *InterestCalculationService) calculateTermAdjustment(termMonths int) decimal.Decimal {
    if termMonths <= 12 {
        return decimal.NewFromFloat(-0.25) // -0.25% for short terms
    } else if termMonths <= 36 {
        return decimal.Zero // No adjustment for standard terms
    } else if termMonths <= 60 {
        return decimal.NewFromFloat(0.5) // +0.5% for medium terms
    } else {
        return decimal.NewFromFloat(1.5) // +1.5% for long terms
    }
}

func (ics *InterestCalculationService) calculatePurposeAdjustment(purpose LoanPurpose) decimal.Decimal {
    switch purpose {
    case DebtConsolidation:
        return decimal.NewFromFloat(-0.5) // -0.5% for debt consolidation (lower risk)
    case HomeImprovement:
        return decimal.NewFromFloat(-0.25) // -0.25% for home improvement (adds value)
    case Education:
        return decimal.NewFromFloat(-0.25) // -0.25% for education (investment in earning potential)
    case BusinessLoan:
        return decimal.NewFromFloat(1.0) // +1% for business loans (higher risk)
    default: // Personal loan
        return decimal.Zero
    }
}

func (ics *InterestCalculationService) calculateCollateralAdjustment(collateral CollateralType) decimal.Decimal {
    switch collateral {
    case SecuredByProperty:
        return decimal.NewFromFloat(-2.0) // -2% for property-secured loans
    case SecuredByVehicle:
        return decimal.NewFromFloat(-1.0) // -1% for vehicle-secured loans
    case SecuredByInvestment:
        return decimal.NewFromFloat(-1.5) // -1.5% for investment-secured loans
    default: // Unsecured
        return decimal.Zero
    }
}

func (ics *InterestCalculationService) calculateMarketAdjustment(condition MarketCondition) decimal.Decimal {
    switch condition {
    case BullMarket:
        return decimal.NewFromFloat(0.25) // +0.25% in bull markets
    case BearMarket:
        return decimal.NewFromFloat(-0.25) // -0.25% in bear markets (encourage lending)
    case VolatileMarket:
        return decimal.NewFromFloat(0.5) // +0.5% in volatile markets (risk premium)
    default: // Stable market
        return decimal.Zero
    }
}

func (ics *InterestCalculationService) applyRateBounds(rate decimal.Decimal) decimal.Decimal {
    minRate := decimal.NewFromFloat(4.0)  // 4% minimum
    maxRate := decimal.NewFromFloat(35.0) // 35% maximum
    
    if rate.LessThan(minRate) {
        return minRate
    }
    if rate.GreaterThan(maxRate) {
        return maxRate
    }
    
    return rate
}
```

### 🎯 Domain Events Architecture

#### Comprehensive Event System
```go
package domain

import (
    "encoding/json"
    "fmt"
    "time"
)

// Base Domain Event Interface
type DomainEvent interface {
    EventName() string
    EventID() string
    AggregateID() string
    AggregateType() string
    EventVersion() int
    Timestamp() time.Time
    UserID() *UserID
    CorrelationID() string
    CausationID() *string
    Payload() interface{}
    Metadata() map[string]interface{}
}

// Base Event Implementation
type BaseEvent struct {
    eventID       string
    aggregateID   string
    aggregateType string
    eventVersion  int
    timestamp     time.Time
    userID        *UserID
    correlationID string
    causationID   *string
    metadata      map[string]interface{}
}

func NewBaseEvent(
    eventID string,
    aggregateID string,
    aggregateType string,
    userID *UserID,
    correlationID string,
    causationID *string,
) BaseEvent {
    return BaseEvent{
        eventID:       eventID,
        aggregateID:   aggregateID,
        aggregateType: aggregateType,
        eventVersion:  1,
        timestamp:     time.Now(),
        userID:        userID,
        correlationID: correlationID,
        causationID:   causationID,
        metadata:      make(map[string]interface{}),
    }
}

func (e BaseEvent) EventID() string { return e.eventID }
func (e BaseEvent) AggregateID() string { return e.aggregateID }
func (e BaseEvent) AggregateType() string { return e.aggregateType }
func (e BaseEvent) EventVersion() int { return e.eventVersion }
func (e BaseEvent) Timestamp() time.Time { return e.timestamp }
func (e BaseEvent) UserID() *UserID { return e.userID }
func (e BaseEvent) CorrelationID() string { return e.correlationID }
func (e BaseEvent) CausationID() *string { return e.causationID }
func (e BaseEvent) Metadata() map[string]interface{} { return e.metadata }

// Loan Domain Events
type LoanProposedEvent struct {
    BaseEvent
    LoanID          LoanID
    BorrowerID      UserID
    PrincipalAmount Money
    InterestRate    InterestRate
    TermMonths      int
    LoanPurpose     LoanPurpose
}

func NewLoanProposedEvent(
    loanID LoanID,
    borrowerID UserID,
    principalAmount Money,
    interestRate InterestRate,
    termMonths int,
    purpose LoanPurpose,
    correlationID string,
) LoanProposedEvent {
    return LoanProposedEvent{
        BaseEvent: NewBaseEvent(
            generateEventID(),
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

func (e LoanProposedEvent) EventName() string { return "loan.proposed" }
func (e LoanProposedEvent) Payload() interface{} { return e }

type LoanApprovedEvent struct {
    BaseEvent
    LoanID             LoanID
    ApproverID         UserID
    ApprovalDate       time.Time
    ApprovedAmount     Money
    ApprovedRate       InterestRate
    RiskAssessment     RiskAssessmentResult
    ApprovalConditions []string
}

func NewLoanApprovedEvent(
    loanID LoanID,
    approverID UserID,
    approvedAmount Money,
    approvedRate InterestRate,
    riskAssessment RiskAssessmentResult,
    conditions []string,
    correlationID string,
) LoanApprovedEvent {
    return LoanApprovedEvent{
        BaseEvent: NewBaseEvent(
            generateEventID(),
            loanID.String(),
            "Loan",
            &approverID,
            correlationID,
            nil,
        ),
        LoanID:             loanID,
        ApproverID:         approverID,
        ApprovalDate:       time.Now(),
        ApprovedAmount:     approvedAmount,
        ApprovedRate:       approvedRate,
        RiskAssessment:     riskAssessment,
        ApprovalConditions: conditions,
    }
}

func (e LoanApprovedEvent) EventName() string { return "loan.approved" }
func (e LoanApprovedEvent) Payload() interface{} { return e }

type InvestmentMadeEvent struct {
    BaseEvent
    LoanID            LoanID
    InvestorID        UserID
    InvestmentID      InvestmentID
    Amount            Money
    InvestmentDate    time.Time
    ExpectedReturn    Money
    RemainingAmount   Money
    IsFullyFunded     bool
}

func NewInvestmentMadeEvent(
    loanID LoanID,
    investorID UserID,
    investmentID InvestmentID,
    amount Money,
    expectedReturn Money,
    remainingAmount Money,
    isFullyFunded bool,
    correlationID string,
) InvestmentMadeEvent {
    return InvestmentMadeEvent{
        BaseEvent: NewBaseEvent(
            generateEventID(),
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

func (e InvestmentMadeEvent) EventName() string { return "investment.made" }
func (e InvestmentMadeEvent) Payload() interface{} { return e }

type LoanFullyFundedEvent struct {
    BaseEvent
    LoanID              LoanID
    TotalAmount         Money
    NumberOfInvestors   int
    FundingCompletedAt  time.Time
    DisbursementEligible bool
}

func NewLoanFullyFundedEvent(
    loanID LoanID,
    totalAmount Money,
    numInvestors int,
    disbursementEligible bool,
    correlationID string,
) LoanFullyFundedEvent {
    return LoanFullyFundedEvent{
        BaseEvent: NewBaseEvent(
            generateEventID(),
            loanID.String(),
            "Loan",
            nil,
            correlationID,
            nil,
        ),
        LoanID:              loanID,
        TotalAmount:         totalAmount,
        NumberOfInvestors:   numInvestors,
        FundingCompletedAt:  time.Now(),
        DisbursementEligible: disbursementEligible,
    }
}

func (e LoanFullyFundedEvent) EventName() string { return "loan.fully_funded" }
func (e LoanFullyFundedEvent) Payload() interface{} { return e }

type LoanDisbursedEvent struct {
    BaseEvent
    LoanID              LoanID
    BorrowerID          UserID
    DisbursedAmount     Money
    DisbursementDate    time.Time
    PaymentSchedule     []PaymentScheduleItem
    AgreementDocumentURL string
    FirstPaymentDue     time.Time
}

func NewLoanDisbursedEvent(
    loanID LoanID,
    borrowerID UserID,
    amount Money,
    paymentSchedule []PaymentScheduleItem,
    agreementURL string,
    firstPaymentDue time.Time,
    correlationID string,
) LoanDisbursedEvent {
    return LoanDisbursedEvent{
        BaseEvent: NewBaseEvent(
            generateEventID(),
            loanID.String(),
            "Loan",
            &borrowerID,
            correlationID,
            nil,
        ),
        LoanID:              loanID,
        BorrowerID:          borrowerID,
        DisbursedAmount:     amount,
        DisbursementDate:    time.Now(),
        PaymentSchedule:     paymentSchedule,
        AgreementDocumentURL: agreementURL,
        FirstPaymentDue:     firstPaymentDue,
    }
}

func (e LoanDisbursedEvent) EventName() string { return "loan.disbursed" }
func (e LoanDisbursedEvent) Payload() interface{} { return e }

// Payment Events
type PaymentReceivedEvent struct {
    BaseEvent
    LoanID          LoanID
    PaymentID       PaymentID
    BorrowerID      UserID
    Amount          Money
    PaymentDate     time.Time
    PaymentMethod   PaymentMethod
    PrincipalPortion Money
    InterestPortion Money
    RemainingBalance Money
    IsLatePayment   bool
    LateFee         Money
}

func NewPaymentReceivedEvent(
    loanID LoanID,
    paymentID PaymentID,
    borrowerID UserID,
    amount Money,
    paymentMethod PaymentMethod,
    principalPortion Money,
    interestPortion Money,
    remainingBalance Money,
    isLate bool,
    lateFee Money,
    correlationID string,
) PaymentReceivedEvent {
    return PaymentReceivedEvent{
        BaseEvent: NewBaseEvent(
            generateEventID(),
            loanID.String(),
            "Loan",
            &borrowerID,
            correlationID,
            nil,
        ),
        LoanID:          loanID,
        PaymentID:       paymentID,
        BorrowerID:      borrowerID,
        Amount:          amount,
        PaymentDate:     time.Now(),
        PaymentMethod:   paymentMethod,
        PrincipalPortion: principalPortion,
        InterestPortion: interestPortion,
        RemainingBalance: remainingBalance,
        IsLatePayment:   isLate,
        LateFee:         lateFee,
    }
}

func (e PaymentReceivedEvent) EventName() string { return "payment.received" }
func (e PaymentReceivedEvent) Payload() interface{} { return e }

// Event Publishing Infrastructure
type DomainEventPublisher interface {
    Publish(ctx context.Context, events []DomainEvent) error
    PublishAsync(ctx context.Context, events []DomainEvent) error
}

type EventSubscriber interface {
    EventName() string
    Handle(ctx context.Context, event DomainEvent) error
}

type EventBus interface {
    Subscribe(subscriber EventSubscriber)
    Publish(ctx context.Context, event DomainEvent) error
    PublishBatch(ctx context.Context, events []DomainEvent) error
}

// Event Serialization
type SerializedEvent struct {
    EventID       string                 `json:"eventId"`
    EventName     string                 `json:"eventName"`
    AggregateID   string                 `json:"aggregateId"`
    AggregateType string                 `json:"aggregateType"`
    EventVersion  int                    `json:"eventVersion"`
    Timestamp     time.Time              `json:"timestamp"`
    UserID        *string                `json:"userId,omitempty"`
    CorrelationID string                 `json:"correlationId"`
    CausationID   *string                `json:"causationId,omitempty"`
    Payload       json.RawMessage        `json:"payload"`
    Metadata      map[string]interface{} `json:"metadata"`
}

func SerializeEvent(event DomainEvent) (*SerializedEvent, error) {
    payload, err := json.Marshal(event.Payload())
    if err != nil {
        return nil, fmt.Errorf("failed to serialize event payload: %w", err)
    }
    
    var userIDStr *string
    if event.UserID() != nil {
        str := event.UserID().String()
        userIDStr = &str
    }
    
    return &SerializedEvent{
        EventID:       event.EventID(),
        EventName:     event.EventName(),
        AggregateID:   event.AggregateID(),
        AggregateType: event.AggregateType(),
        EventVersion:  event.EventVersion(),
        Timestamp:     event.Timestamp(),
        UserID:        userIDStr,
        CorrelationID: event.CorrelationID(),
        CausationID:   event.CausationID(),
        Payload:       payload,
        Metadata:      event.Metadata(),
    }, nil
}

func generateEventID() string {
    // Implementation would generate unique event ID
    return fmt.Sprintf("event-%d", time.Now().UnixNano())
}
```

---

## Detailed Implementation Strategy

### 🚀 8-Phase Migration Roadmap

#### Phase 1: Foundation Setup (Weeks 1-2)
**Objective**: Establish value objects and basic domain structure

**Tasks:**
1. **Value Objects Creation** (Week 1)
   ```go
   // Create new files in internal/loan/domain/valueobjects/
   internal/loan/domain/valueobjects/
   ├── money.go              ✅ Currency-aware money calculations
   ├── interest_rate.go      ✅ Rate validation and calculations
   ├── credit_score.go       ✅ Credit score with risk categories
   ├── loan_id.go           ✅ Strongly-typed loan identifiers
   ├── user_id.go           ✅ Strongly-typed user identifiers
   └── currency.go          ✅ Currency enumeration
   ```

2. **Domain Errors** (Week 1)
   ```go
   // Create internal/loan/domain/errors/
   ├── domain_errors.go           ✅ Base domain error types
   └── business_rule_errors.go    ✅ Business rule violations
   ```

3. **Initial Entity Enhancement** (Week 2)
   ```go
   // Enhance existing entities
   internal/loan/domain/entities/
   ├── loan.go              ✅ Add basic business methods
   └── investment.go        ✅ Investment entity within aggregate
   ```

**Success Criteria:**
- All primitive types replaced with value objects
- Basic business validation at type level
- Domain errors properly structured
- Existing tests pass with new types

**Risk Mitigation:**
- Maintain backward compatibility with adapters
- Create type conversion utilities
- Gradual rollout with feature flags

---

#### Phase 2: Rich Domain Model (Weeks 3-4)
**Objective**: Move business logic from interactors to domain entities

**Tasks:**
1. **Aggregate Root Enhancement** (Week 3)
   ```go
   // Enhanced Loan aggregate with business methods
   func (l *Loan) Approve(approverID UserID) error
   func (l *Loan) AddInvestment(investorID UserID, amount Money) error
   func (l *Loan) Disburse(agreementURL string) error
   func (l *Loan) IsValid() error
   ```

2. **Business Logic Migration** (Week 4)
   - Move validation from interactors to entities
   - Implement state transition methods
   - Add aggregate invariant checking
   - Create domain event placeholders

**Migration Strategy:**
```go
// Dual implementation approach
type CreateProposedLoanInteractor struct {
    // Old implementation
    legacyStore    LegacyStore
    
    // New implementation  
    loanRepo       domain.LoanRepository
    domainService  *domain.LoanService
    
    // Feature flag
    useDomainModel bool
}

func (c *CreateProposedLoanInteractor) Execute(ctx context.Context, in Input) error {
    if c.useDomainModel {
        return c.executeWithDomainModel(ctx, in)
    }
    return c.executeLegacy(ctx, in)
}
```

**Success Criteria:**
- 60% of business logic moved to domain layer
- All state transitions encapsulated in entities
- Interactors reduced to orchestration
- Domain invariants properly enforced

---

#### Phase 3: Domain Services (Weeks 5-6)  
**Objective**: Extract complex business logic into domain services

**Tasks:**
1. **Risk Assessment Service** (Week 5)
   ```go
   type RiskAssessmentService struct {
       creditBureauPort   CreditBureauPort
       fraudDetectionPort FraudDetectionPort
   }
   
   func (ras *RiskAssessmentService) AssessLoanRisk(
       borrower *User, 
       loanAmount Money,
   ) (*RiskAssessmentResult, error)
   ```

2. **Interest Calculation Service** (Week 5)
   ```go
   type InterestCalculationService struct {
       marketDataPort MarketDataPort
   }
   
   func (ics *InterestCalculationService) CalculateOptimalRate(
       input RateCalculationInput,
   ) (InterestRate, error)
   ```

3. **Loan Matching Service** (Week 6)
   ```go
   type LoanMatchingService struct {
       investorRepo domain.InvestorRepository
   }
   
   func (lms *LoanMatchingService) FindMatchingInvestors(
       loan *Loan,
   ) ([]*Investor, error)
   ```

**Integration Pattern:**
```go
// Application service uses domain services
type LoanApplicationService struct {
    loanRepo       domain.LoanRepository
    riskService    *domain.RiskAssessmentService
    rateService    *domain.InterestCalculationService
    matchingService *domain.LoanMatchingService
}
```

**Success Criteria:**
- Complex calculations extracted from interactors
- Domain services properly integrated
- External service dependencies abstracted
- Business logic fully centralized

---

#### Phase 4: Event-Driven Architecture (Weeks 7-8)
**Objective**: Implement domain events for system integration

**Tasks:**
1. **Domain Event Infrastructure** (Week 7)
   ```go
   // Event interfaces and base implementation
   type DomainEvent interface {
       EventName() string
       EventID() string
       Timestamp() time.Time
       AggregateID() string
   }
   
   // Event publisher port
   type DomainEventPublisher interface {
       Publish(events []DomainEvent) error
   }
   ```

2. **Loan Domain Events** (Week 7)
   ```go
   // Key events for loan lifecycle
   type LoanProposedEvent struct { ... }
   type LoanApprovedEvent struct { ... }
   type InvestmentMadeEvent struct { ... }
   type LoanDisbursedEvent struct { ... }
   ```

3. **Event Integration** (Week 8)
   ```go
   // Repository publishes events after persistence
   func (r *SQLLoanRepository) Save(loan *domain.Loan) error {
       // ... save aggregate
       if err := r.eventPublisher.Publish(loan.Events()); err != nil {
           r.logger.Error("Failed to publish events", "error", err)
       }
       loan.ClearEvents()
       return nil
   }
   ```

**Success Criteria:**
- All state changes generate domain events
- Events properly published after persistence
- Event subscribers handle system integration
- Audit trail and debugging capabilities

---

#### Phase 5: Repository Pattern (Weeks 9-10)
**Objective**: Implement proper aggregate repositories

**Tasks:**
1. **Repository Interface Design** (Week 9)
   ```go
   type LoanRepository interface {
       Save(loan *Loan) error
       FindByID(id LoanID) (*Loan, error)
       FindByBorrower(borrowerID UserID) ([]*Loan, error)
       FindByStatus(status LoanStatus) ([]*Loan, error)
   }
   ```

2. **SQL Repository Implementation** (Week 10)
   ```go
   type SQLLoanRepository struct {
       db             *sql.DB
       eventPublisher domain.DomainEventPublisher
   }
   
   // Aggregate reconstruction from SQL entities
   func (r *SQLLoanRepository) mapSQLToDomainAggregate(
       loanEntity *LoanSQLEntity, 
       investmentEntities []*InvestmentSQLEntity,
   ) (*domain.Loan, error)
   ```

**Key Patterns:**
- Aggregate boundary respect
- Optimistic locking
- Event publishing integration
- Transaction management

**Success Criteria:**
- Proper aggregate boundaries maintained
- No N+1 query problems
- Consistent event publishing
- Performance benchmarks met

---

#### Phase 6: Application Layer Refactoring (Weeks 11-12)
**Objective**: Implement CQRS and clean application services

**Tasks:**
1. **CQRS Implementation** (Week 11)
   ```go
   // Separate command and query handlers
   internal/loan/application/
   ├── commands/
   │   ├── create_loan_command.go
   │   ├── approve_loan_command.go
   │   └── invest_loan_command.go
   └── queries/
       ├── loan_queries.go
       └── investment_queries.go
   ```

2. **Application Service Redesign** (Week 12)
   ```go
   type CreateLoanCommandHandler struct {
       loanRepo        domain.LoanRepository
       riskService     *domain.RiskAssessmentService
       eventPublisher  domain.DomainEventPublisher
   }
   
   func (h *CreateLoanCommandHandler) Handle(
       ctx context.Context, 
       cmd CreateLoanCommand,
   ) error {
       // Orchestration only, business logic in domain
   }
   ```

**Success Criteria:**
- Clear separation of read/write operations
- Thin application services (orchestration only)
- Proper use case boundaries
- Clean dependency injection

---

#### Phase 7: Infrastructure Enhancement (Weeks 13-14)
**Objective**: Enhance infrastructure adapters and monitoring

**Tasks:**
1. **Enhanced Adapters** (Week 13)
   ```go
   // Standardized adapter interfaces
   type HealthCheckable interface {
       HealthCheck() error
   }
   
   type Configurable interface {
       Configure(config interface{}) error
   }
   
   // All adapters implement standard interfaces
   type CreditBureauAdapter struct {
       client HTTPClient
       config CreditBureauConfig
   }
   ```

2. **Monitoring and Observability** (Week 14)
   ```go
   // Domain metrics
   type DomainMetrics interface {
       IncrementLoanCreated()
       IncrementLoanApproved()
       RecordRiskAssessmentTime(duration time.Duration)
   }
   ```

**Success Criteria:**
- All adapters follow standard patterns
- Comprehensive health checks
- Domain-specific metrics
- Proper logging and tracing

---

#### Phase 8: Project Structure Transformation (Weeks 15-16)
**Objective**: Reorganize project structure for proper DDD

**Tasks:**
1. **Directory Restructuring** (Week 15)
   ```bash
   # Migration script for directory restructure
   ./scripts/migrate_project_structure.sh
   
   # New structure validation
   ./scripts/validate_structure.sh
   ```

2. **Final Integration and Testing** (Week 16)
   ```go
   // Comprehensive integration tests
   internal/loan/test/integration/
   ├── loan_lifecycle_test.go
   ├── payment_flow_test.go
   └── event_propagation_test.go
   ```

**Migration Script Example:**
```bash
#!/bin/bash
# migrate_project_structure.sh

# Create new directory structure
mkdir -p internal/loan/domain/{entities,valueobjects,services,events,repositories,errors}
mkdir -p internal/loan/application/{commands,queries,services,workflows}
mkdir -p internal/loan/infrastructure/{persistence,http,events,external}

# Move existing files to new locations
mv internal/loan/internal/entity/loan.go internal/loan/domain/entities/
mv internal/loan/internal/gateway/ internal/loan/infrastructure/persistence/
mv internal/loan/internal/interactor/ internal/loan/application/commands/

# Update import paths
find . -name "*.go" -exec sed -i 's|internal/loan/internal/entity|internal/loan/domain/entities|g' {} \;
find . -name "*.go" -exec sed -i 's|internal/loan/internal/gateway|internal/loan/infrastructure/persistence|g' {} \;

echo "Structure migration completed. Run tests to verify."
```

**Success Criteria:**
- New structure fully implemented
- All tests pass
- Import paths updated
- Documentation reflects new structure

---

### 🎯 Validation and Testing Strategy

#### Domain Layer Testing
```go
// Rich unit tests for domain logic
func TestLoan_Approve_ValidProposedLoan_ShouldSucceed(t *testing.T) {
    // Given
    loan := createTestLoanProposal()
    approverID := domain.NewUserID(123)
    
    // When
    err := loan.Approve(approverID)
    
    // Then
    assert.NoError(t, err)
    assert.Equal(t, domain.Approved, loan.Status())
    assert.Equal(t, approverID, *loan.ApproverID())
    assert.Len(t, loan.Events(), 1)
    assert.Equal(t, "loan.approved", loan.Events()[0].EventName())
}

func TestLoan_Approve_AlreadyApprovedLoan_ShouldFail(t *testing.T) {
    // Given
    loan := createTestLoanProposal()
    loan.Approve(domain.NewUserID(123))
    
    // When
    err := loan.Approve(domain.NewUserID(456))
    
    // Then
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "can only approve proposed loans")
}
```

#### Integration Testing Strategy
```go
// End-to-end loan lifecycle test
func TestLoanLifecycle_CompleteFlow(t *testing.T) {
    // Setup
    ctx := context.Background()
    testDB := setupTestDatabase()
    
    // Create loan
    createCmd := application.CreateLoanCommand{
        LoanID:     generateTestLoanID(),
        BorrowerID: generateTestUserID(),
        Amount:     mustNewMoney(decimal.NewFromInt(10000), domain.USD),
        TermMonths: 36,
    }
    
    err := commandHandler.Handle(ctx, createCmd)
    assert.NoError(t, err)
    
    // Approve loan  
    approveCmd := application.ApproveLoanCommand{
        LoanID:     createCmd.LoanID,
        ApproverID: generateTestUserID(),
    }
    
    err = approveCommandHandler.Handle(ctx, approveCmd)
    assert.NoError(t, err)
    
    // Add investment
    investCmd := application.InvestLoanCommand{
        LoanID:     createCmd.LoanID,
        InvestorID: generateTestUserID(),
        Amount:     mustNewMoney(decimal.NewFromInt(10000), domain.USD),
    }
    
    err = investCommandHandler.Handle(ctx, investCmd)
    assert.NoError(t, err)
    
    // Verify final state
    loan, err := loanRepo.FindByID(createCmd.LoanID)
    assert.NoError(t, err)
    assert.Equal(t, domain.FullyFunded, loan.Status())
}
```

### 📊 Success Metrics and KPIs

#### Technical Quality Metrics
```go
// Automated metrics collection
type ArchitectureMetrics struct {
    DomainLogicCentralization  float64 // Target: >90%
    CyclomaticComplexity      int     // Target: <10 per method
    TestCoverage              float64 // Target: >95% domain layer
    DependencyDirection       bool    // Target: 100% compliance
}

// Quality gates
func ValidateArchitectureQuality() error {
    metrics := CollectArchitectureMetrics()
    
    if metrics.DomainLogicCentralization < 0.9 {
        return errors.New("domain logic centralization below 90%")
    }
    
    if metrics.TestCoverage < 0.95 {
        return errors.New("domain test coverage below 95%")
    }
    
    return nil
}
```

#### Business Value Metrics
- **Bug Fix Time**: Target 60% reduction through centralized business logic
- **Feature Development Speed**: Target 40% improvement through rich domain model  
- **Code Reusability**: Target 80% of business logic reusable across contexts
- **Business Rule Consistency**: Target 100% consistency through domain layer

---

## Project Structure Transformation

### 🏗️ Current vs Target Structure

#### Current Structure Issues
```
internal/loan/
├── internal/          ❌ Unnecessary nesting
│   ├── entity/        ❌ Anemic entities mixed with SQL concerns
│   ├── gateway/       ❌ Mixed HTTP and SQL in same layer
│   ├── interactor/    ❌ Business logic scattered
│   ├── usecase/       ❌ Just input/output DTOs
│   └── mocks/         ✅ Generated mocks (keep)
└── module.go          ✅ Dependency injection (enhance)
```

#### Target DDD Structure
```
internal/loan/                           # 📋 Loan Bounded Context
├── domain/                             # 🎯 Domain Layer (Pure Business Logic)
│   ├── entities/                       # Rich Aggregates & Entities
│   │   ├── loan.go                     # Loan aggregate root
│   │   ├── investment.go               # Investment entity
│   │   └── borrower.go                 # Borrower entity
│   ├── valueobjects/                   # Value Objects with Business Rules
│   │   ├── money.go                    # Currency-aware money
│   │   ├── interest_rate.go            # Interest rate with calculations
│   │   ├── credit_score.go             # Credit score with risk categories
│   │   └── loan_id.go                  # Strongly-typed identifiers
│   ├── services/                       # Domain Services
│   │   ├── risk_assessment_service.go  # Credit risk evaluation
│   │   ├── interest_calculation_service.go # Dynamic rate calculation
│   │   └── loan_matching_service.go    # Investor-loan matching
│   ├── events/                         # Domain Events
│   │   ├── loan_events.go              # Loan lifecycle events
│   │   └── investment_events.go        # Investment events
│   ├── repositories/                   # Repository Interfaces (Ports)
│   │   ├── loan_repository.go          # Loan aggregate repository
│   │   └── user_repository.go          # User repository
│   └── errors/                         # Domain-Specific Errors
│       ├── domain_errors.go            # Base domain errors
│       └── business_rule_errors.go     # Business rule violations
├── application/                        # 🔧 Application Layer (Use Case Orchestration)
│   ├── commands/                       # Command Handlers (CQRS Write)
│   │   ├── create_loan_command.go      # Create loan proposal
│   │   ├── approve_loan_command.go     # Approve loan application
│   │   └── invest_loan_command.go      # Make investment
│   ├── queries/                        # Query Handlers (CQRS Read)
│   │   ├── loan_queries.go             # Loan query operations
│   │   └── investment_queries.go       # Investment queries
│   ├── services/                       # Application Services
│   │   ├── loan_application_service.go # Cross-aggregate operations
│   │   └── notification_service.go     # User notifications
│   └── dto/                            # Data Transfer Objects
│       ├── loan_dto.go                 # Loan DTOs
│       └── investment_dto.go           # Investment DTOs
├── infrastructure/                     # 🔌 Infrastructure Layer (External Concerns)
│   ├── persistence/                    # Repository Implementations
│   │   ├── sql/                        # SQL Database Implementations
│   │   │   ├── loan_sql_repository.go  # SQL loan repository
│   │   │   └── entities/               # SQL mapping entities
│   │   │       └── loan_sql_entity.go  # SQL-specific entities
│   │   └── memory/                     # In-Memory (Testing)
│   │       └── loan_memory_repository.go
│   ├── http/                           # HTTP Adapters (Primary)
│   │   ├── handlers/                   # HTTP handlers
│   │   │   ├── loan_handler.go         # Loan API endpoints
│   │   │   └── investment_handler.go   # Investment endpoints
│   │   ├── middleware/                 # HTTP middleware
│   │   └── serializers/                # JSON serialization
│   ├── events/                         # Event Infrastructure
│   │   ├── publishers/                 # Event publishers
│   │   └── subscribers/                # Event subscribers
│   └── external/                       # External Service Adapters
│       ├── credit_bureau_adapter.go    # Credit bureau integration
│       └── payment_processor_adapter.go # Payment processing
├── test/                               # Testing
│   ├── integration/                    # Integration tests
│   ├── acceptance/                     # E2E acceptance tests
│   └── fixtures/                       # Test data
└── module.go                          # Enhanced DI container
```

### 🔄 Migration Strategy

#### Step 1: Create New Structure (Parallel Development)
```bash
# Create new directory structure without disturbing current code
mkdir -p internal/loan/domain/{entities,valueobjects,services,events,repositories,errors}
mkdir -p internal/loan/application/{commands,queries,services,dto}
mkdir -p internal/loan/infrastructure/{persistence/sql/entities,http/{handlers,middleware},events,external}
mkdir -p internal/loan/test/{integration,acceptance,fixtures}
```

#### Step 2: Value Objects Migration
```go
// Create new value objects in internal/loan/domain/valueobjects/
// Start with Money, InterestRate, LoanID, UserID
// Update existing DTOs to use new value objects gradually
```

#### Step 3: Domain Entities Enhancement  
```go
// Move and enhance entities from internal/loan/internal/entity/ 
// to internal/loan/domain/entities/
// Add business methods and domain logic
```

#### Step 4: Application Layer Refactoring
```go
// Transform interactors into command/query handlers
// Move from internal/loan/internal/interactor/
// to internal/loan/application/commands/ and queries/
```

#### Step 5: Infrastructure Reorganization
```go
// Separate SQL concerns from HTTP concerns
// Move HTTP handlers to infrastructure/http/handlers/
// Move SQL implementations to infrastructure/persistence/sql/
```

#### Step 6: Legacy Cleanup
```bash
# Remove old structure after migration is complete and tested
rm -rf internal/loan/internal/
```

### 📋 Migration Checklist

#### Domain Layer ✅
- [ ] Create value objects (Money, InterestRate, CreditScore, IDs)
- [ ] Enhance entities with business methods
- [ ] Implement domain services for complex logic
- [ ] Add domain events infrastructure
- [ ] Define repository interfaces
- [ ] Create domain-specific error types

#### Application Layer ✅
- [ ] Implement command handlers for write operations
- [ ] Implement query handlers for read operations  
- [ ] Create application services for cross-aggregate operations
- [ ] Design DTOs for external communication
- [ ] Implement application workflows

#### Infrastructure Layer ✅
- [ ] Separate HTTP handlers from business logic
- [ ] Create proper SQL repository implementations
- [ ] Implement event publishing infrastructure
- [ ] Create external service adapters
- [ ] Add monitoring and health checks

#### Testing Strategy ✅
- [ ] Unit tests for domain logic (entities, value objects, services)
- [ ] Integration tests for repository implementations
- [ ] API tests for HTTP handlers
- [ ] End-to-end acceptance tests for business workflows
- [ ] Performance tests for critical paths

---

## Migration Roadmap

### 🗓️ 16-Week Detailed Schedule

#### Weeks 1-4: Foundation and Domain Model
**Week 1: Value Objects Foundation**
- Day 1-2: Create Money value object with currency support
- Day 3-4: Implement InterestRate with business validation
- Day 5: Add strongly-typed IDs (LoanID, UserID, InvestmentID)

**Week 2: Domain Entities Enhancement**
- Day 1-3: Enhance Loan entity with business methods
- Day 4-5: Create Investment and Borrower entities
- Weekend: Add domain event placeholders

**Week 3: Domain Services (Part 1)**
- Day 1-3: Implement RiskAssessmentService
- Day 4-5: Create InterestCalculationService
- Weekend: Add integration tests

**Week 4: Domain Services (Part 2)**
- Day 1-3: Implement LoanMatchingService  
- Day 4-5: Create PaymentScheduleService
- Weekend: Domain layer integration testing

#### Weeks 5-8: Application and Events
**Week 5: Command Handlers**
- Day 1-2: CreateLoanCommandHandler
- Day 3-4: ApproveLoanCommandHandler
- Day 5: InvestLoanCommandHandler

**Week 6: Query Handlers and CQRS**
- Day 1-3: Implement query handlers for loan operations
- Day 4-5: Separate read/write models
- Weekend: CQRS pattern validation

**Week 7: Domain Events Infrastructure**
- Day 1-3: Event interfaces and base implementations
- Day 4-5: Loan lifecycle events
- Weekend: Event publishing integration

**Week 8: Event Integration**
- Day 1-3: Repository event publishing
- Day 4-5: Event subscribers and handlers
- Weekend: Event-driven workflow testing

#### Weeks 9-12: Repository and Infrastructure
**Week 9: Repository Pattern**
- Day 1-3: Design repository interfaces
- Day 4-5: Implement SQL repository for Loan aggregate
- Weekend: Repository testing with SQL mocks

**Week 10: Repository Implementation**
- Day 1-3: Complete SQL repository implementation
- Day 4-5: Add caching layer and optimization
- Weekend: Performance testing and optimization

**Week 11: Application Service Refactoring**
- Day 1-3: Refactor interactors to use new domain model
- Day 4-5: Implement feature flags for gradual rollout
- Weekend: A/B testing setup

**Week 12: Infrastructure Adapters**
- Day 1-3: HTTP handlers refactoring
- Day 4-5: External service adapters (credit bureau, payment processor)
- Weekend: Adapter integration testing

#### Weeks 13-16: Finalization and Migration
**Week 13: Monitoring and Observability**
- Day 1-3: Add domain-specific metrics
- Day 4-5: Implement health checks and logging
- Weekend: Observability dashboard setup

**Week 14: Project Structure Migration**
- Day 1-3: Execute directory restructuring
- Day 4-5: Update import paths and dependencies
- Weekend: Full regression testing

**Week 15: Integration and End-to-End Testing**
- Day 1-3: Comprehensive integration test suite
- Day 4-5: End-to-end acceptance tests
- Weekend: Performance benchmarking

**Week 16: Production Deployment and Validation**
- Day 1-3: Production deployment with feature flags
- Day 4-5: Monitor metrics and rollback capabilities
- Weekend: Success metrics validation

### 🎯 Success Milestones

#### Month 1 (Weeks 1-4): Domain Foundation
- ✅ All primitive types replaced with value objects
- ✅ Rich domain model with business methods
- ✅ Domain services extracting complex logic
- ✅ 70% domain test coverage

#### Month 2 (Weeks 5-8): Application and Events  
- ✅ CQRS pattern fully implemented
- ✅ Domain events for all state changes
- ✅ Application services orchestrate domain logic
- ✅ Event-driven integrations working

#### Month 3 (Weeks 9-12): Infrastructure and Integration
- ✅ Repository pattern with aggregate boundaries
- ✅ Infrastructure adapters standardized
- ✅ Feature flags enable gradual rollout
- ✅ Performance benchmarks met

#### Month 4 (Weeks 13-16): Production and Validation
- ✅ New structure deployed to production
- ✅ All legacy code removed
- ✅ Success metrics achieved
- ✅ Team trained on new architecture

---

## Success Metrics & Validation

### 📊 Quantitative Success Metrics

#### Code Quality Metrics
```go
// Automated architecture validation
type QualityMetrics struct {
    DomainLogicCentralization float64 // Target: >90%
    BusinessLogicInDomain     float64 // Target: >85%
    TestCoverageCore          float64 // Target: >95%
    CyclomaticComplexity      int     // Target: <10 avg
    DependencyRuleCompliance  float64 // Target: 100%
    CodeDuplication          float64 // Target: <5%
}

func ValidateArchitectureQuality() *QualityReport {
    return &QualityReport{
        DomainCentralization: measureDomainLogic(),
        TestCoverage:        measureTestCoverage(),
        Complexity:          measureComplexity(),
        Dependencies:        validateDependencyRules(),
        Duplication:         measureDuplication(),
        OverallGrade:        calculateOverallGrade(),
    }
}
```

#### Performance Metrics
```go
// Performance benchmarks
type PerformanceMetrics struct {
    LoanCreationLatency   time.Duration // Target: <100ms p95
    ApprovalProcessTime   time.Duration // Target: <500ms p95  
    DatabaseQueryTime     time.Duration // Target: <50ms p95
    EventProcessingTime   time.Duration // Target: <10ms p95
    ThroughputLoansPerSec int           // Target: >1000
}

func BenchmarkLoanOperations(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Benchmark loan creation with new domain model
        start := time.Now()
        loan := createTestLoan()
        duration := time.Since(start)
        
        if duration > 100*time.Millisecond {
            b.Errorf("Loan creation too slow: %v", duration)
        }
    }
}
```

#### Business Value Metrics
- **Bug Resolution Time**: Target 60% improvement
  - Current: 2-3 days average
  - Target: <1 day average
- **Feature Development Speed**: Target 40% improvement  
  - Current: 2 weeks per feature
  - Target: 1.2 weeks per feature
- **Business Rule Consistency**: Target 100%
  - Current: Manual validation, inconsistencies possible
  - Target: Automated domain validation, zero inconsistencies
- **Developer Onboarding**: Target 50% improvement
  - Current: 2 weeks to productive
  - Target: 1 week to productive

### ✅ Qualitative Success Indicators

#### Domain Expert Collaboration
- Business rules clearly expressed in ubiquitous language
- Domain experts can read and validate domain model code
- Business logic changes require minimal technical translation
- Domain model serves as living documentation

#### Developer Experience
- Clear separation of concerns reduces cognitive load
- Business logic concentrated in domain layer
- Infrastructure changes don't affect business logic
- Tests are focused and meaningful

#### System Maintainability  
- Business rule changes localized to domain layer
- New features integrate cleanly with existing model
- Technical debt concentrated in infrastructure, not domain
- Refactoring is safe and predictable

### 🔍 Validation Strategy

#### Pre-Migration Baseline
```bash
# Collect baseline metrics before starting migration
./scripts/collect_baseline_metrics.sh

# Results stored for comparison
metrics/
├── baseline_complexity.json
├── baseline_test_coverage.json  
├── baseline_performance.json
└── baseline_business_metrics.json
```

#### Continuous Validation
```yaml
# CI/CD pipeline validation
validation:
  domain_tests:
    command: "go test ./internal/loan/domain/... -v"
    coverage_threshold: 95%
    
  architecture_rules:
    command: "./scripts/validate_architecture.sh"
    rules:
      - "Domain layer has no infrastructure dependencies"
      - "Business logic percentage > 85%"
      - "Cyclomatic complexity < 10 average"
      
  performance_benchmarks:
    command: "go test -bench=. ./internal/loan/..."
    thresholds:
      loan_creation: "100ms"
      approval_process: "500ms"
```

#### Monthly Progress Reviews
```go
// Monthly architecture health check
type HealthCheckReport struct {
    Timestamp           time.Time
    QualityMetrics      QualityMetrics
    PerformanceMetrics  PerformanceMetrics
    BusinessMetrics     BusinessMetrics
    TechnicalDebt       TechnicalDebtScore
    TeamSatisfaction    int // 1-10 scale
    RecommendedActions  []string
}

func GenerateMonthlyReport() *HealthCheckReport {
    return &HealthCheckReport{
        Timestamp:          time.Now(),
        QualityMetrics:     collectQualityMetrics(),
        PerformanceMetrics: collectPerformanceMetrics(),
        BusinessMetrics:    collectBusinessMetrics(),
        TechnicalDebt:      assessTechnicalDebt(),
        TeamSatisfaction:   surveyTeam(),
        RecommendedActions: generateRecommendations(),
    }
}
```

---

## Future Architecture Considerations

### 🚀 Evolution Path (12-24 Months)

#### Event Sourcing Implementation
```go
// Future: Event sourcing for complete audit trail
type EventSourcedLoan struct {
    id     LoanID
    events []DomainEvent
}

func (esl *EventSourcedLoan) ReplayEvents(events []DomainEvent) *Loan {
    // Rebuild aggregate state from event stream
    loan := &Loan{}
    for _, event := range events {
        loan.Apply(event)
    }
    return loan
}

// Benefits:
// - Complete audit trail
// - Temporal queries (state at any point in time)
// - Business intelligence from event stream
// - Regulatory compliance
```

#### CQRS with Separate Read Models
```go
// Future: Optimized read models for different use cases
type LoanReadModel struct {
    ID              string
    BorrowerName    string
    Amount          decimal.Decimal
    Status          string
    RiskScore       string
    InvestorCount   int
    FundingProgress float64
}

type InvestmentReadModel struct {
    LoanID         string
    InvestorID     string
    Amount         decimal.Decimal
    ExpectedReturn decimal.Decimal
    Status         string
}

// Separate projections for different contexts:
// - Borrower dashboard
// - Investor marketplace  
// - Risk management
// - Regulatory reporting
```

#### Microservices Architecture
```go
// Future: Domain-based microservice boundaries
services/
├── loan-service/          # Core loan management
├── risk-service/          # Risk assessment and scoring
├── investment-service/    # Investment and matching
├── payment-service/       # Payment processing
├── notification-service/  # User communications
└── reporting-service/     # Analytics and compliance

// Each service:
// - Own database
// - Domain-driven boundaries
// - Event-driven communication
// - Independent deployment
```

#### Advanced Domain Patterns
```go
// Future: Saga pattern for complex workflows
type LoanApprovalSaga struct {
    sagaID         SagaID
    loanID         LoanID
    state          SagaState
    compensations  []CompensationAction
}

func (las *LoanApprovalSaga) Handle(event DomainEvent) error {
    switch e := event.(type) {
    case LoanProposedEvent:
        return las.startRiskAssessment(e.LoanID)
    case RiskAssessmentCompletedEvent:
        return las.processApprovalDecision(e)
    case LoanApprovedEvent:
        return las.notifyStakeholders(e)
    }
}

// Specification pattern for complex business rules
type LoanApprovalSpecification struct {
    riskService     *RiskAssessmentService
    complianceRules []ComplianceRule
}

func (las *LoanApprovalSpecification) IsSatisfiedBy(loan *Loan) bool {
    return las.checkRiskCriteria(loan) &&
           las.checkComplianceRules(loan) &&
           las.checkBusinessRules(loan)
}
```

### 📈 Scalability Considerations

#### Performance Optimization
- **Read Replicas**: Separate read/write database instances
- **Caching Strategy**: Redis for frequently accessed data
- **Event Sourcing**: Stream processing for real-time analytics
- **CQRS Projections**: Materialized views for complex queries

#### Compliance and Security
- **GDPR Compliance**: Right to erasure with event sourcing
- **SOX Compliance**: Complete audit trail through events
- **PCI Compliance**: Secure payment processing boundaries
- **Data Encryption**: End-to-end encryption for sensitive data

#### International Expansion
- **Multi-Currency Support**: Enhanced Money value object
- **Regulatory Variations**: Country-specific business rules
- **Localization**: Domain language and business terms
- **Tax Compliance**: Jurisdiction-specific calculations

### 🎯 Technology Evolution

#### Modern Go Patterns
```go
// Future: Generics for type-safe value objects
type ID[T any] struct {
    value uint64
}

type LoanID = ID[Loan]
type UserID = ID[User]

// Context-aware operations
func (r *Repository) FindWithContext[T Entity](
    ctx context.Context, 
    id ID[T],
) (*T, error) {
    // Type-safe repository operations
}
```

#### Cloud-Native Architecture  
- **Kubernetes Deployment**: Container orchestration
- **Service Mesh**: Inter-service communication
- **Cloud Events**: Standardized event format
- **Observability**: Distributed tracing and monitoring

This comprehensive guide provides a complete roadmap for transforming Daneizo into a truly domain-driven architecture that will scale with business growth while maintaining code quality and developer productivity.

<function_calls>
<invoke name="TodoWrite">
<parameter name="todos">[{"content": "Analyze all three documents to identify key themes and structure", "status": "completed", "activeForm": "Analyzing all three documents to identify key themes and structure"}, {"content": "Create unified document structure with logical flow", "status": "completed", "activeForm": "Creating unified document structure with logical flow"}, {"content": "Merge and enhance content with detailed examples", "status": "in_progress", "activeForm": "Merging and enhancing content with detailed examples"}, {"content": "Add comprehensive implementation roadmap", "status": "pending", "activeForm": "Adding comprehensive implementation roadmap"}]