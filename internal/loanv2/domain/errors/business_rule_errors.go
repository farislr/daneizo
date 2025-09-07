package errors

// Business rule specific errors

// LoanBusinessRuleError represents loan-specific business rule violations
type LoanBusinessRuleError struct {
	*DomainError
}

// NewLoanBusinessRuleError creates a new loan business rule error
func NewLoanBusinessRuleError(message string) *LoanBusinessRuleError {
	return &LoanBusinessRuleError{
		DomainError: NewDomainErrorWithCode(message, "LOAN_BUSINESS_RULE_VIOLATION"),
	}
}

// Common loan business rule errors
var (
	ErrLoanNotProposed           = NewLoanBusinessRuleError("loan is not in proposed status")
	ErrLoanNotApproved           = NewLoanBusinessRuleError("loan is not in approved status")
	ErrLoanAlreadyApproved       = NewLoanBusinessRuleError("loan is already approved")
	ErrLoanNotFullyFunded        = NewLoanBusinessRuleError("loan is not fully funded")
	ErrLoanAlreadyDisbursed      = NewLoanBusinessRuleError("loan is already disbursed")
	ErrInvestmentExceedsLoan     = NewLoanBusinessRuleError("investment amount exceeds loan principal")
	ErrMinimumLoanAmount         = NewLoanBusinessRuleError("loan amount is below minimum required")
	ErrMaximumLoanAmount         = NewLoanBusinessRuleError("loan amount exceeds maximum allowed")
	ErrInvalidInterestRate       = NewLoanBusinessRuleError("interest rate is outside acceptable range")
	ErrInsufficientCreditScore   = NewLoanBusinessRuleError("credit score does not meet minimum requirements")
	ErrHighRiskProfile           = NewLoanBusinessRuleError("borrower risk profile exceeds acceptable limits")
	ErrInvalidLoanTerm           = NewLoanBusinessRuleError("loan term is outside acceptable range")
	ErrBorrowerDebtRatioTooHigh  = NewLoanBusinessRuleError("borrower debt-to-income ratio exceeds maximum allowed")
	ErrFraudRiskDetected         = NewLoanBusinessRuleError("fraud risk detected - application rejected")
	ErrAgreementDocumentRequired = NewLoanBusinessRuleError("loan agreement document is required")
)

// InvestmentBusinessRuleError represents investment-specific business rule violations
type InvestmentBusinessRuleError struct {
	*DomainError
}

// NewInvestmentBusinessRuleError creates a new investment business rule error
func NewInvestmentBusinessRuleError(message string) *InvestmentBusinessRuleError {
	return &InvestmentBusinessRuleError{
		DomainError: NewDomainErrorWithCode(message, "INVESTMENT_BUSINESS_RULE_VIOLATION"),
	}
}

// Common investment business rule errors
var (
	ErrMinimumInvestmentAmount = NewInvestmentBusinessRuleError("investment amount is below minimum required")
	ErrMaximumInvestmentAmount = NewInvestmentBusinessRuleError("investment amount exceeds maximum allowed")
	ErrInvestorNotEligible     = NewInvestmentBusinessRuleError("investor is not eligible for this loan")
	ErrSelfInvestment          = NewInvestmentBusinessRuleError("borrower cannot invest in their own loan")
	ErrInsufficientFunds       = NewInvestmentBusinessRuleError("investor has insufficient funds")
	ErrLoanNotInvestable       = NewInvestmentBusinessRuleError("loan is not available for investment")
)

// PaymentBusinessRuleError represents payment-specific business rule violations
type PaymentBusinessRuleError struct {
	*DomainError
}

// NewPaymentBusinessRuleError creates a new payment business rule error
func NewPaymentBusinessRuleError(message string) *PaymentBusinessRuleError {
	return &PaymentBusinessRuleError{
		DomainError: NewDomainErrorWithCode(message, "PAYMENT_BUSINESS_RULE_VIOLATION"),
	}
}

// Common payment business rule errors
var (
	ErrPaymentAmountInvalid     = NewPaymentBusinessRuleError("payment amount is invalid")
	ErrPaymentAlreadyProcessed  = NewPaymentBusinessRuleError("payment has already been processed")
	ErrLatePaymentFeeRequired   = NewPaymentBusinessRuleError("late payment fee is required")
	ErrInsufficientPayment      = NewPaymentBusinessRuleError("payment amount is insufficient")
	ErrPaymentMethodNotAllowed  = NewPaymentBusinessRuleError("payment method is not allowed")
	ErrPaymentProcessingFailed  = NewPaymentBusinessRuleError("payment processing failed")
)