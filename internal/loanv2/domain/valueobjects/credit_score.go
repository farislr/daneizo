package valueobjects

import "fmt"

// CreditScore represents a credit score with business validation and risk categorization
type CreditScore struct {
	score int
}

// RiskCategory represents different risk levels based on credit score
type RiskCategory int

const (
	LowRisk RiskCategory = iota + 1
	MediumRisk
	HighRisk
	UnacceptableRisk
)

func (r RiskCategory) String() string {
	return [...]string{"", "Low", "Medium", "High", "Unacceptable"}[r]
}

// NewCreditScore creates a new CreditScore with validation
func NewCreditScore(score int) (CreditScore, error) {
	if score < 300 || score > 850 {
		return CreditScore{}, fmt.Errorf("credit score must be between 300 and 850, got: %d", score)
	}
	return CreditScore{score: score}, nil
}

// Value returns the credit score value
func (cs CreditScore) Value() int {
	return cs.score
}

// IsExcellent returns true if credit score is excellent (750+)
func (cs CreditScore) IsExcellent() bool {
	return cs.score >= 750
}

// IsGood returns true if credit score is good (650+)
func (cs CreditScore) IsGood() bool {
	return cs.score >= 650
}

// IsFair returns true if credit score is fair (600+)
func (cs CreditScore) IsFair() bool {
	return cs.score >= 600
}

// IsPoor returns true if credit score is poor (<600)
func (cs CreditScore) IsPoor() bool {
	return cs.score < 600
}

// RiskCategory returns the risk category based on the credit score
func (cs CreditScore) RiskCategory() RiskCategory {
	if cs.IsExcellent() {
		return LowRisk
	} else if cs.IsGood() {
		return MediumRisk
	} else if cs.IsFair() {
		return HighRisk
	}
	return UnacceptableRisk
}

// String returns a formatted string representation
func (cs CreditScore) String() string {
	return fmt.Sprintf("%d", cs.score)
}

// StringWithCategory returns the score with its risk category
func (cs CreditScore) StringWithCategory() string {
	return fmt.Sprintf("%d (%s Risk)", cs.score, cs.RiskCategory().String())
}

// Equal checks if this credit score equals another
func (cs CreditScore) Equal(other CreditScore) bool {
	return cs.score == other.score
}

// GreaterThan checks if this credit score is greater than another
func (cs CreditScore) GreaterThan(other CreditScore) bool {
	return cs.score > other.score
}

// LessThan checks if this credit score is less than another
func (cs CreditScore) LessThan(other CreditScore) bool {
	return cs.score < other.score
}