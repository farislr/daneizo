# DDD Project Structure Guide for Daneizo

## Overview

This guide provides a comprehensive project structure redesign to support the Domain-Driven Design improvements outlined in the DDD Improvement Guide. The new structure follows DDD principles, Clean Architecture layers, and Go conventions.

---

## 🏗️ Current vs Target Structure

### Current Structure (Anemic)
```
internal/loan/
├── internal/
│   ├── entity/sqlentity/     # Data containers only
│   ├── gateway/              # Mixed concerns (HTTP + SQL)
│   ├── interactor/           # Business logic scattered
│   ├── usecase/              # Input/Output structs
│   └── mocks/                # Generated mocks
└── module.go                 # Dependency injection
```

**Problems with Current Structure:**
- ❌ Business logic scattered across interactors
- ❌ SQL entities mixed with domain concepts
- ❌ No clear domain layer separation
- ❌ Gateways handling multiple concerns
- ❌ Missing value objects and domain services

### Target Structure (Rich DDD)
```
internal/
├── loan/                     # Bounded Context
│   ├── domain/              # 🎯 Domain Layer (Core Business Logic)
│   │   ├── entities/        # Aggregates & Rich Entities
│   │   │   ├── loan.go              # Loan aggregate root
│   │   │   ├── investment.go        # Investment entity
│   │   │   └── borrower.go          # Borrower entity
│   │   ├── valueobjects/    # Value Objects
│   │   │   ├── money.go             # Money with currency
│   │   │   ├── interest_rate.go     # Interest rate with validation
│   │   │   ├── loan_id.go           # Strongly-typed loan ID
│   │   │   ├── user_id.go           # Strongly-typed user ID
│   │   │   └── credit_score.go      # Credit score with business rules
│   │   ├── services/        # Domain Services
│   │   │   ├── risk_assessment.go   # Risk calculation logic
│   │   │   ├── interest_calculation.go # Interest rate calculations
│   │   │   └── loan_matching.go     # Investor-loan matching
│   │   ├── events/          # Domain Events
│   │   │   ├── loan_events.go       # Loan lifecycle events
│   │   │   └── investment_events.go # Investment events
│   │   ├── specifications/  # Business Rules as Specifications
│   │   │   └── loan_specifications.go
│   │   └── repositories/    # Repository Interfaces (Ports)
│   │       ├── loan_repository.go   # Loan aggregate repository
│   │       └── user_repository.go   # User repository
│   ├── application/         # 🔧 Application Layer (Use Cases)
│   │   ├── commands/        # Command Handlers (Write Operations)
│   │   │   ├── create_loan_command.go
│   │   │   ├── approve_loan_command.go
│   │   │   ├── invest_loan_command.go
│   │   │   └── disburse_loan_command.go
│   │   ├── queries/         # Query Handlers (Read Operations)
│   │   │   ├── loan_queries.go
│   │   │   ├── investment_queries.go
│   │   │   └── borrower_queries.go
│   │   ├── services/        # Application Services (Orchestration)
│   │   │   ├── loan_application_service.go
│   │   │   └── investment_application_service.go
│   │   └── dto/            # Data Transfer Objects
│   │       ├── loan_dto.go
│   │       └── investment_dto.go
│   ├── infrastructure/      # 🔌 Infrastructure Layer (External Concerns)
│   │   ├── persistence/     # Repository Implementations
│   │   │   ├── sql/
│   │   │   │   ├── loan_sql_repository.go
│   │   │   │   ├── user_sql_repository.go
│   │   │   │   └── entities/        # SQL mapping entities
│   │   │   │       ├── loan_sql_entity.go
│   │   │   │       └── user_sql_entity.go
│   │   │   └── memory/     # In-memory implementations (testing)
│   │   │       └── loan_memory_repository.go
│   │   ├── http/           # HTTP Adapters
│   │   │   ├── handlers/
│   │   │   │   ├── loan_handler.go
│   │   │   │   └── investment_handler.go
│   │   │   ├── middleware/
│   │   │   │   └── validation_middleware.go
│   │   │   └── routes.go
│   │   ├── events/         # Event Infrastructure
│   │   │   ├── publishers/
│   │   │   │   └── domain_event_publisher.go
│   │   │   └── subscribers/
│   │   │       └── loan_event_subscriber.go
│   │   └── external/       # External Service Adapters
│   │       ├── credit_bureau_adapter.go
│   │       └── payment_processor_adapter.go
│   ├── test/               # Integration & Acceptance Tests
│   │   ├── integration/
│   │   │   └── loan_integration_test.go
│   │   └── acceptance/
│   │       └── loan_acceptance_test.go
│   └── module.go           # Dependency Injection & Module Setup
├── shared/                 # 🌐 Shared Kernel
│   ├── domain/            # Common Domain Concepts
│   │   ├── events/
│   │   │   ├── domain_event.go      # Base domain event interface
│   │   │   └── event_store.go       # Event store interface
│   │   ├── errors/
│   │   │   ├── domain_error.go      # Domain-specific errors
│   │   │   └── business_error.go    # Business rule violations
│   │   ├── valueobjects/
│   │   │   ├── currency.go          # Currency enumeration
│   │   │   └── timestamp.go         # Business timestamp
│   │   └── specifications/
│   │       └── specification.go     # Specification pattern base
│   └── infrastructure/    # Shared Infrastructure
│       ├── persistence/
│       │   ├── transaction.go       # Transaction management
│       │   └── connection.go        # Database connection
│       ├── events/
│       │   ├── event_bus.go         # Event bus implementation
│       │   └── event_dispatcher.go  # Event dispatcher
│       ├── logging/
│       │   └── structured_logger.go # Structured logging
│       └── validation/
│           └── validator.go         # Domain validation utilities
└── app/                    # 🚀 Application Bootstrap
    ├── config/
    │   └── config.go
    ├── bootstrap/
    │   └── app.go
    └── server/
        └── http_server.go
```

---

## 📋 Layer Responsibilities

### 1. Domain Layer (`internal/loan/domain/`)
**Purpose**: Contains core business logic, rules, and domain concepts

**Key Principles**:
- ✅ **Framework Independent**: No external dependencies
- ✅ **Rich Behavior**: Entities contain business methods
- ✅ **Invariant Enforcement**: Aggregates maintain consistency
- ✅ **Business Language**: Uses ubiquitous language

**Package Contents**:
```go
// entities/loan.go - Rich aggregate root
type Loan struct {
    id              LoanID
    borrowerID      UserID
    principalAmount Money
    // ... other fields
}

func (l *Loan) Approve(approverID UserID) error {
    if l.status != Proposed {
        return NewDomainError("can only approve proposed loans")
    }
    // Business logic here...
}

// valueobjects/money.go - Value object with validation
type Money struct {
    amount   decimal.Decimal
    currency Currency
}

func NewMoney(amount decimal.Decimal, currency Currency) (Money, error) {
    if amount.IsNegative() {
        return Money{}, errors.New("amount cannot be negative")
    }
    return Money{amount: amount, currency: currency}, nil
}

// services/risk_assessment.go - Domain service
type RiskAssessmentService struct {
    creditBureauPort CreditBureauPort
}

func (ras *RiskAssessmentService) AssessRisk(borrower Borrower, amount Money) RiskScore {
    // Complex business logic for risk assessment
}
```

### 2. Application Layer (`internal/loan/application/`)
**Purpose**: Orchestrates domain objects to fulfill use cases

**Key Principles**:
- ✅ **Use Case Coordination**: Orchestrates domain objects
- ✅ **Transaction Management**: Handles cross-aggregate transactions
- ✅ **DTO Transformation**: Converts between domain and external models
- ✅ **Security Enforcement**: Applies authorization rules

**Package Contents**:
```go
// commands/create_loan_command.go - Command handler
type CreateLoanCommandHandler struct {
    loanRepo      domain.LoanRepository
    userRepo      domain.UserRepository
    eventPublisher shared.EventPublisher
}

func (h *CreateLoanCommandHandler) Handle(cmd CreateLoanCommand) error {
    // 1. Load domain objects
    // 2. Execute business logic
    // 3. Save changes
    // 4. Publish events
}

// services/loan_application_service.go - Application service
type LoanApplicationService struct {
    loanRepo         domain.LoanRepository
    riskService      *domain.RiskAssessmentService
    interestService  *domain.InterestCalculationService
}

func (las *LoanApplicationService) ProcessLoanApplication(input CreateLoanInput) error {
    // Orchestrate multiple domain services and repositories
}
```

### 3. Infrastructure Layer (`internal/loan/infrastructure/`)
**Purpose**: Implements technical concerns and external integrations

**Key Principles**:
- ✅ **Dependency Implementation**: Implements domain interfaces
- ✅ **External Integration**: Handles databases, HTTP, events
- ✅ **Framework Specific**: Can use external libraries
- ✅ **Adapter Pattern**: Adapts external services to domain interfaces

**Package Contents**:
```go
// persistence/sql/loan_sql_repository.go - Repository implementation
type LoanSQLRepository struct {
    db *sql.DB
    eventPublisher domain.EventPublisher
}

func (r *LoanSQLRepository) Save(loan *domain.Loan) error {
    // 1. Map domain aggregate to SQL entities
    // 2. Save to database in transaction
    // 3. Publish domain events
    // 4. Clear events from aggregate
}

// http/handlers/loan_handler.go - HTTP adapter
type LoanHandler struct {
    loanService *application.LoanApplicationService
}

func (h *LoanHandler) CreateLoan(w http.ResponseWriter, r *http.Request) {
    // 1. Parse HTTP request
    // 2. Convert to command/query
    // 3. Call application service
    // 4. Return HTTP response
}
```

---

## 🚀 Migration Strategy

### Phase 1: Foundation Setup (Week 1-2)
**Goal**: Create new structure alongside existing code

1. **Create New Directory Structure**
```bash
mkdir -p internal/loan/domain/{entities,valueobjects,services,events,repositories}
mkdir -p internal/loan/application/{commands,queries,services,dto}
mkdir -p internal/loan/infrastructure/{persistence/sql,http/handlers,events}
mkdir -p internal/shared/{domain,infrastructure}
```

2. **Implement Value Objects First**
   - Start with `Money`, `InterestRate`, `LoanID`
   - Add comprehensive tests
   - No breaking changes to existing code

3. **Create Domain Events Infrastructure**
   - Base event interfaces
   - Event publisher interface
   - In-memory event bus for testing

### Phase 2: Rich Domain Model (Week 3-4)
**Goal**: Migrate business logic to domain layer

1. **Create Rich Entities**
   - Transform anemic `Loan` to rich aggregate
   - Move business logic from interactors
   - Add business method tests

2. **Extract Domain Services**
   - Risk assessment service
   - Interest calculation service
   - Create interfaces and implementations

3. **Maintain Backward Compatibility**
```go
// Legacy adapter during transition
type LegacyCreateLoanInteractor struct {
    newDomainService *domain.LoanService
}

func (l *LegacyCreateLoanInteractor) Execute(input usecase.CreateProposedLoanInput) error {
    // Convert legacy input to domain objects
    // Call new domain service
    // Return in legacy format
}
```

### Phase 3: Application Layer (Week 5-6)
**Goal**: Implement CQRS and proper use case orchestration

1. **Command/Query Handlers**
   - Separate read and write operations
   - Implement command handlers
   - Create query handlers for read models

2. **Application Services**
   - Orchestrate multiple domain services
   - Handle cross-cutting concerns
   - Transaction management

### Phase 4: Infrastructure Completion (Week 7-8)
**Goal**: Complete infrastructure layer and remove legacy code

1. **Repository Pattern**
   - Implement proper aggregate repositories
   - Add optimistic locking
   - Event sourcing preparation

2. **HTTP Layer Redesign**
   - RESTful API design
   - Proper error handling
   - OpenAPI documentation

3. **Legacy Code Removal**
   - Remove old interactors
   - Clean up SQL entities
   - Update dependency injection

---

## 📁 File Organization Rules

### Naming Conventions
```go
// Domain entities - simple names
loan.go              // Loan aggregate root
investment.go        // Investment entity

// Value objects - descriptive names  
money.go            // Money value object
interest_rate.go    // InterestRate value object
loan_id.go         // LoanID value object

// Domain services - with _service suffix
risk_assessment_service.go     // RiskAssessmentService
interest_calculation_service.go // InterestCalculationService

// Events - grouped by aggregate
loan_events.go         // All loan-related events
investment_events.go   // All investment-related events

// Repositories - interface vs implementation
repositories/loan_repository.go           // Interface
persistence/sql/loan_sql_repository.go   // Implementation

// Commands/Queries - with _command/_query suffix
create_loan_command.go    // CreateLoanCommand and handler
loan_queries.go          // Loan query handlers
```

### Package Import Rules
```go
// ✅ Allowed dependencies
internal/loan/domain          → internal/shared/domain
internal/loan/application     → internal/loan/domain
internal/loan/infrastructure  → internal/loan/application
internal/loan/infrastructure  → internal/loan/domain

// ❌ Forbidden dependencies  
internal/loan/domain         → internal/loan/application     // Domain can't depend on application
internal/loan/domain         → internal/loan/infrastructure  // Domain can't depend on infrastructure
internal/shared/domain       → internal/loan/domain          // Shared can't depend on specific domains
```

### Testing Structure
```
internal/loan/
├── domain/
│   ├── entities/
│   │   ├── loan.go
│   │   └── loan_test.go              # Unit tests alongside code
│   └── services/
│       ├── risk_assessment_service.go  
│       └── risk_assessment_service_test.go
├── application/
│   └── commands/
│       ├── create_loan_command.go
│       └── create_loan_command_test.go # Integration tests
└── test/                              # Separate test package
    ├── integration/
    │   └── loan_integration_test.go   # Full integration tests
    └── acceptance/
        └── loan_acceptance_test.go    # End-to-end tests
```

---

## 🔧 Dependency Injection Setup

### Module Pattern
```go
// internal/loan/module.go
package loan

import (
    "database/sql"
    
    "github.com/farislr/daneizo/internal/loan/application"
    "github.com/farislr/daneizo/internal/loan/domain"
    "github.com/farislr/daneizo/internal/loan/infrastructure"
)

type Module struct {
    // Domain services
    RiskAssessmentService    *domain.RiskAssessmentService
    InterestCalculationService *domain.InterestCalculationService
    
    // Application services  
    LoanApplicationService   *application.LoanApplicationService
    
    // Repositories
    LoanRepository          domain.LoanRepository
    UserRepository          domain.UserRepository
}

func NewModule(db *sql.DB, eventPublisher shared.EventPublisher) *Module {
    // Create repositories
    loanRepo := infrastructure.NewLoanSQLRepository(db, eventPublisher)
    userRepo := infrastructure.NewUserSQLRepository(db)
    
    // Create domain services
    riskService := domain.NewRiskAssessmentService(/* external ports */)
    interestService := domain.NewInterestCalculationService()
    
    // Create application services
    loanAppService := application.NewLoanApplicationService(
        loanRepo, 
        userRepo, 
        riskService, 
        interestService,
    )
    
    return &Module{
        RiskAssessmentService:      riskService,
        InterestCalculationService: interestService,
        LoanApplicationService:     loanAppService,
        LoanRepository:            loanRepo,
        UserRepository:            userRepo,
    }
}
```

### Wire Generation (Optional)
```go
//go:build wireinject
// +build wireinject

package main

import (
    "github.com/google/wire"
    "github.com/farislr/daneizo/internal/loan"
)

func InitializeLoanModule(db *sql.DB) (*loan.Module, error) {
    wire.Build(
        loan.NewModule,
        infrastructure.NewLoanSQLRepository,
        domain.NewRiskAssessmentService,
        // ... other dependencies
    )
    return &loan.Module{}, nil
}
```

---

## 📚 Documentation Structure

### Per-Bounded Context Documentation
```
internal/loan/
├── README.md                    # Loan domain overview
├── DOMAIN_MODEL.md             # Domain model documentation
├── UBIQUITOUS_LANGUAGE.md      # Business terminology
└── docs/
    ├── adr/                    # Architecture Decision Records
    │   ├── 001-aggregate-boundaries.md
    │   ├── 002-event-sourcing.md
    │   └── 003-cqrs-implementation.md
    ├── api/
    │   └── loan-api.md         # API documentation
    └── examples/
        └── loan-lifecycle.md   # Usage examples
```

### Root Documentation
```
docs/
├── ARCHITECTURE_OVERVIEW.md    # High-level architecture
├── DDD_IMPROVEMENT_GUIDE.md    # DDD implementation guide
├── DDD_PROJECT_STRUCTURE.md   # This document
└── DEVELOPMENT_GUIDE.md       # Development practices
```

---

## 🎯 Success Metrics

### Code Quality Metrics
- **Domain Logic Concentration**: 80%+ of business logic in domain layer
- **Cyclomatic Complexity**: <10 for domain methods
- **Test Coverage**: 95%+ for domain layer, 85%+ overall
- **Dependency Violations**: 0 violations of layer dependency rules

### Development Velocity Metrics  
- **Feature Implementation Time**: 40% reduction through rich domain model
- **Bug Fix Time**: 60% reduction through centralized business logic
- **New Developer Onboarding**: 50% faster with clear structure

### Business Value Metrics
- **Business Rule Consistency**: 0 inconsistencies across use cases
- **Requirement Traceability**: 100% business rules mapped to domain code
- **Domain Expert Collaboration**: Improved through ubiquitous language

---

## 🔮 Future Considerations

### Multi-Bounded Context Support
When the system grows, additional bounded contexts can be added:
```
internal/
├── loan/           # Loan management bounded context
├── user/           # User management bounded context  
├── payment/        # Payment processing bounded context
├── notification/   # Notification bounded context
├── reporting/      # Reporting bounded context
└── shared/         # Shared kernel
```

### Advanced Patterns
- **Event Sourcing**: Complete audit trail for regulatory compliance
- **CQRS**: Separate read/write models for performance optimization
- **Saga Pattern**: Distributed transaction management across contexts
- **Domain Events**: Integration between bounded contexts

This structure provides a solid foundation for implementing rich domain-driven design while maintaining the flexibility to evolve as the P2P lending platform grows in complexity and scale.