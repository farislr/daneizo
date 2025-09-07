package valueobjects

import "fmt"

// LoanID represents a strongly-typed loan identifier to prevent mixing with other ID types
type LoanID struct {
	id uint64
}

// NewLoanID creates a new LoanID with validation
func NewLoanID(id uint64) LoanID {
	if id == 0 {
		panic("LoanID cannot be zero")
	}
	return LoanID{id: id}
}

// Value returns the underlying uint64 value
func (l LoanID) Value() uint64 {
	return l.id
}

// String returns a formatted string representation
func (l LoanID) String() string {
	return fmt.Sprintf("loan-%d", l.id)
}

// Equal checks if this LoanID equals another
func (l LoanID) Equal(other LoanID) bool {
	return l.id == other.id
}

// IsZero returns true if the LoanID is zero (invalid)
func (l LoanID) IsZero() bool {
	return l.id == 0
}