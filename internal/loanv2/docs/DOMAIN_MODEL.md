# Domain Model Documentation

## Loan Aggregate

The `Loan` aggregate is the central business entity in the lending domain. It encapsulates all business rules and invariants related to loan lifecycle management.

### Aggregate Root: Loan

#### Core Properties
- **ID**: Unique loan identifier (strongly-typed `LoanID`)
- **BorrowerID**: User who requested the loan (strongly-typed `UserID`)
- **PrincipalAmount**: Loan amount requested (`Money` value object)
- **InvestedAmount**: Amount currently invested (`Money` value object)
- **InterestRate**: Applied interest rate (`InterestRate` value object)
- **TermMonths**: Loan term in months (6-60 months)
- **Status**: Current loan status (`LoanStatus` enum)

#### Business Methods

##### Loan Creation
```go
func NewLoanProposal(
    id LoanID,
    borrowerID UserID,
    principalAmount Money,
    interestRate InterestRate,
    termMonths int,
    purpose string,
) (*Loan, error)
```
**Business Rules:**
- Minimum loan amount: $1,000
- Maximum loan amount: $1,000,000
- Term must be 6-60 months
- Interest rate must be 0-30%
- Automatically generates `LoanProposedEvent`

##### Loan Approval
```go
func (l *Loan) Approve(approverID UserID, conditions []string) error
```
**Business Rules:**
- Can only approve loans in `Proposed` or `UnderReview` status
- Borrower cannot approve their own loan
- Automatically transitions to `Approved` status
- Generates `LoanApprovedEvent`

##### Investment Management
```go
func (l *Loan) AddInvestment(
    investorID UserID,
    amount Money,
    expectedReturn Money,
) error
```
**Business Rules:**
- Only approved or partially funded loans can accept investments
- Borrower cannot invest in their own loan
- Investment cannot exceed remaining loan amount
- Automatically manages status transitions:
  - `Approved` → `PartiallyFunded` → `FullyFunded`
- Generates `InvestmentMadeEvent` and optionally `LoanFullyFundedEvent`

##### Loan Disbursement
```go
func (l *Loan) Disburse(agreementURL string, firstPaymentDue time.Time) error
```
**Business Rules:**
- Only fully funded loans can be disbursed
- Agreement document URL is required
- First payment must be at least 1 month from disbursement
- Automatically transitions to `Disbursed` status
- Generates `LoanDisbursedEvent`

#### Status Transitions

```
PROPOSED → APPROVED → PARTIALLY_FUNDED → FULLY_FUNDED → DISBURSED → ACTIVE → PAID_OFF
    ↓           ↓             ↓               ↓
REJECTED   CANCELLED    CANCELLED      CANCELLED
```

#### Aggregate Invariants
1. **Funding Constraint**: `InvestedAmount ≤ PrincipalAmount`
2. **Status Consistency**: Disbursed loans must be fully funded
3. **Agreement Requirement**: Disbursed loans must have agreement document
4. **Investment Separation**: Approved loans should not have investments yet

### Value Objects

#### Money
Currency-aware monetary value with business operations:
```go
type Money struct {
    amount   decimal.Decimal
    currency Currency
}
```

**Key Methods:**
- `Add(Money) (Money, error)` - Same currency addition
- `Subtract(Money) (Money, error)` - Same currency subtraction
- `MultiplyByRate(decimal.Decimal) Money` - Rate multiplication
- `GreaterThan(Money) (bool, error)` - Same currency comparison

#### InterestRate
Interest rate with business validation and calculations:
```go
type InterestRate struct {
    rate decimal.Decimal // as percentage (0-100)
}
```

**Key Methods:**
- `CalculateMonthlyPayment(Money, int) Money` - Monthly payment calculation
- `CalculateTotalInterest(Money, int) Money` - Total interest calculation
- `AsMonthlyRate() decimal.Decimal` - Convert to monthly decimal rate

#### Strongly-Typed IDs
- `LoanID` - Prevents mixing loan IDs with other entity IDs
- `UserID` - Prevents mixing user IDs with other entity IDs

#### CreditScore
Credit score with risk categorization:
```go
type CreditScore struct {
    score int // 300-850
}
```

**Risk Categories:**
- **Excellent**: 750+ (Low Risk)
- **Good**: 650-749 (Medium Risk)
- **Fair**: 600-649 (High Risk)
- **Poor**: <600 (Unacceptable Risk)

### Domain Events

#### LoanProposedEvent
Emitted when a new loan is created:
```go
type LoanProposedEvent struct {
    LoanID          LoanID
    BorrowerID      UserID
    PrincipalAmount Money
    InterestRate    InterestRate
    TermMonths      int
    LoanPurpose     string
}
```

#### LoanApprovedEvent
Emitted when a loan is approved:
```go
type LoanApprovedEvent struct {
    LoanID         LoanID
    ApproverID     UserID
    ApprovedAmount Money
    ApprovedRate   InterestRate
    Conditions     []string
}
```

#### InvestmentMadeEvent
Emitted when an investment is made:
```go
type InvestmentMadeEvent struct {
    LoanID            LoanID
    InvestorID        UserID
    Amount            Money
    ExpectedReturn    Money
    RemainingAmount   Money
    IsFullyFunded     bool
}
```

#### LoanFullyFundedEvent
Emitted when loan reaches full funding:
```go
type LoanFullyFundedEvent struct {
    LoanID              LoanID
    TotalAmount         Money
    NumberOfInvestors   int
    DisbursementEligible bool
}
```

#### LoanDisbursedEvent
Emitted when loan funds are disbursed:
```go
type LoanDisbursedEvent struct {
    LoanID              LoanID
    BorrowerID          UserID
    DisbursedAmount     Money
    AgreementDocumentURL string
    FirstPaymentDue     time.Time
}
```

### Domain Services

#### RiskAssessmentService
Performs comprehensive risk evaluation:

**Input:**
- Borrower information
- Requested loan amount
- Loan term

**Output:**
- Overall risk category (Low/Medium/High/Unacceptable)
- Credit score assessment
- Fraud risk evaluation
- Debt-to-income ratio calculation
- Recommended interest rate
- Maximum loan amount
- Human-readable reasoning

**Business Logic:**
- Credit bureau integration for credit scores
- Fraud detection service integration
- Debt-to-income ratio calculations
- Risk-based interest rate determination
- Maximum borrowing capacity assessment

### Repository Interfaces

#### LoanRepository
Primary repository for loan aggregate persistence:

**Key Methods:**
- `Save(context.Context, *Loan) error` - Persist aggregate
- `FindByID(context.Context, LoanID) (*Loan, error)` - Load by ID
- `FindByStatus(context.Context, LoanStatus) ([]*Loan, error)` - Query by status
- `FindRequiringApproval(context.Context) ([]*Loan, error)` - Business query
- `FindAvailableForInvestment(context.Context) ([]*Loan, error)` - Business query

#### UserRepository
Repository for user information needed by loan domain:

**Key Methods:**
- `FindByID(context.Context, UserID) (*User, error)` - Load user data
- `UpdateCreditScore(context.Context, UserID, CreditScore) error` - Update credit
- `IsEligibleForLoan(context.Context, UserID, Money) (bool, error)` - Eligibility check

## Business Rules Summary

### Loan Creation Rules
1. Minimum amount: $1,000
2. Maximum amount: $1,000,000
3. Term range: 6-60 months
4. Maximum interest rate: 30%
5. Risk assessment determines final rate

### Approval Rules
1. Only proposed loans can be approved
2. Self-approval not allowed
3. Risk assessment influences decision
4. Approval conditions can be specified

### Investment Rules
1. Self-investment not allowed
2. Investment amount validated against remaining principal
3. Status automatically progresses based on funding level
4. Over-investment prevented by business rules

### Disbursement Rules
1. Only fully funded loans eligible
2. Agreement document required
3. First payment date validation (minimum 1 month)
4. Automatic status progression to active

## Error Handling

### Domain Errors
- `DomainError` - Base domain error
- `ValidationError` - Field validation failures
- `NotFoundError` - Entity not found errors

### Business Rule Errors
- `LoanBusinessRuleError` - Loan-specific business violations
- `InvestmentBusinessRuleError` - Investment-specific violations
- `PaymentBusinessRuleError` - Payment-specific violations

### Common Error Scenarios
- `ErrLoanNotProposed` - Attempting operations on wrong status
- `ErrInvestmentExceedsLoan` - Investment amount validation
- `ErrSelfInvestment` - Business rule violation
- `ErrAgreementDocumentRequired` - Required field validation