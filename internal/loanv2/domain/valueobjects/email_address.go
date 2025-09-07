package valueobjects

import (
	"fmt"
	"regexp"
	"strings"
)

// EmailAddress represents a validated email address
type EmailAddress struct {
	address string
}

// emailRegex is a basic email validation pattern
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewEmailAddress creates a new EmailAddress with validation
func NewEmailAddress(address string) (EmailAddress, error) {
	address = strings.TrimSpace(strings.ToLower(address))
	
	if address == "" {
		return EmailAddress{}, fmt.Errorf("email address cannot be empty")
	}
	
	if len(address) > 254 {
		return EmailAddress{}, fmt.Errorf("email address too long: %d characters", len(address))
	}
	
	if !emailRegex.MatchString(address) {
		return EmailAddress{}, fmt.Errorf("invalid email address format: %s", address)
	}
	
	return EmailAddress{address: address}, nil
}

// Value returns the email address as a string
func (e EmailAddress) Value() string {
	return e.address
}

// String returns the email address as a string
func (e EmailAddress) String() string {
	return e.address
}

// Domain returns the domain part of the email address
func (e EmailAddress) Domain() string {
	parts := strings.Split(e.address, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// LocalPart returns the local part of the email address (before @)
func (e EmailAddress) LocalPart() string {
	parts := strings.Split(e.address, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[0]
}

// Equal checks if this email address equals another
func (e EmailAddress) Equal(other EmailAddress) bool {
	return e.address == other.address
}

// IsEmpty returns true if the email address is empty
func (e EmailAddress) IsEmpty() bool {
	return e.address == ""
}