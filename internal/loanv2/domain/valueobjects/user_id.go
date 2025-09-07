package valueobjects

import "fmt"

// UserID represents a strongly-typed user identifier to prevent mixing with other ID types
type UserID struct {
	id uint64
}

// NewUserID creates a new UserID with validation
func NewUserID(id uint64) UserID {
	if id == 0 {
		panic("UserID cannot be zero")
	}
	return UserID{id: id}
}

// Value returns the underlying uint64 value
func (u UserID) Value() uint64 {
	return u.id
}

// String returns a formatted string representation
func (u UserID) String() string {
	return fmt.Sprintf("user-%d", u.id)
}

// Equal checks if this UserID equals another
func (u UserID) Equal(other UserID) bool {
	return u.id == other.id
}

// IsZero returns true if the UserID is zero (invalid)
func (u UserID) IsZero() bool {
	return u.id == 0
}