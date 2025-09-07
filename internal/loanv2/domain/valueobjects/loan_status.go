package valueobjects

import "fmt"

// LoanStatus represents the status of a loan in its lifecycle
type LoanStatus string

const (
	Proposed       LoanStatus = "PROPOSED"        // Loan application submitted
	UnderReview    LoanStatus = "UNDER_REVIEW"    // Loan is being reviewed
	Approved       LoanStatus = "APPROVED"        // Loan approved, ready for investment
	PartiallyFunded LoanStatus = "PARTIALLY_FUNDED" // Some investments made, not fully funded
	FullyFunded    LoanStatus = "FULLY_FUNDED"    // Loan fully funded, ready for disbursement
	Disbursed      LoanStatus = "DISBURSED"       // Funds disbursed to borrower
	Active         LoanStatus = "ACTIVE"          // Loan is active with payment schedule
	PaidOff        LoanStatus = "PAID_OFF"        // Loan fully repaid
	Defaulted      LoanStatus = "DEFAULTED"       // Loan in default
	Rejected       LoanStatus = "REJECTED"        // Loan application rejected
	Cancelled      LoanStatus = "CANCELLED"       // Loan cancelled before disbursement
)

// String returns the string representation of the loan status
func (ls LoanStatus) String() string {
	return string(ls)
}

// IsValid checks if the loan status is valid
func (ls LoanStatus) IsValid() bool {
	switch ls {
	case Proposed, UnderReview, Approved, PartiallyFunded, FullyFunded, 
		 Disbursed, Active, PaidOff, Defaulted, Rejected, Cancelled:
		return true
	}
	return false
}

// CanTransitionTo checks if the current status can transition to the target status
func (ls LoanStatus) CanTransitionTo(target LoanStatus) bool {
	transitions := map[LoanStatus][]LoanStatus{
		Proposed: {UnderReview, Approved, Rejected, Cancelled},
		UnderReview: {Approved, Rejected, Cancelled},
		Approved: {PartiallyFunded, FullyFunded, Cancelled},
		PartiallyFunded: {FullyFunded, Cancelled},
		FullyFunded: {Disbursed, Cancelled},
		Disbursed: {Active},
		Active: {PaidOff, Defaulted},
		PaidOff: {}, // Terminal state
		Defaulted: {}, // Terminal state  
		Rejected: {}, // Terminal state
		Cancelled: {}, // Terminal state
	}
	
	allowedTransitions, exists := transitions[ls]
	if !exists {
		return false
	}
	
	for _, allowed := range allowedTransitions {
		if allowed == target {
			return true
		}
	}
	return false
}

// IsTerminal returns true if this is a terminal status (no further transitions possible)
func (ls LoanStatus) IsTerminal() bool {
	switch ls {
	case PaidOff, Defaulted, Rejected, Cancelled:
		return true
	}
	return false
}

// IsActive returns true if the loan is in an active state (accepting investments or payments)
func (ls LoanStatus) IsActive() bool {
	switch ls {
	case Approved, PartiallyFunded, FullyFunded, Active:
		return true
	}
	return false
}

// CanAcceptInvestments returns true if the loan can accept new investments
func (ls LoanStatus) CanAcceptInvestments() bool {
	switch ls {
	case Approved, PartiallyFunded:
		return true
	}
	return false
}

// CanBeDisbursed returns true if the loan can be disbursed
func (ls LoanStatus) CanBeDisbursed() bool {
	return ls == FullyFunded
}

// CanAcceptPayments returns true if the loan can accept payments
func (ls LoanStatus) CanAcceptPayments() bool {
	return ls == Active
}

// NewLoanStatus creates a new LoanStatus with validation
func NewLoanStatus(status string) (LoanStatus, error) {
	ls := LoanStatus(status)
	if !ls.IsValid() {
		return "", fmt.Errorf("invalid loan status: %s", status)
	}
	return ls, nil
}