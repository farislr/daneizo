package valueobjects

import (
	"fmt"
	"regexp"
	"strings"
)

// PhoneNumber represents a validated and formatted phone number
type PhoneNumber struct {
	number string // stored in E.164 format (+country_code_number)
}

// phoneRegex matches various phone number formats
var phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{6,14}$`)

// NewPhoneNumber creates a new PhoneNumber with validation and formatting
func NewPhoneNumber(number string) (PhoneNumber, error) {
	// Clean the input
	cleaned := strings.ReplaceAll(number, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")
	cleaned = strings.TrimSpace(cleaned)
	
	if cleaned == "" {
		return PhoneNumber{}, fmt.Errorf("phone number cannot be empty")
	}
	
	// Add + if not present and number doesn't start with +
	if !strings.HasPrefix(cleaned, "+") {
		// Assume Indonesian number if no country code
		if len(cleaned) >= 10 && (strings.HasPrefix(cleaned, "08") || strings.HasPrefix(cleaned, "628")) {
			if strings.HasPrefix(cleaned, "08") {
				// Convert 08xx to +628xx
				cleaned = "+628" + cleaned[2:]
			} else if strings.HasPrefix(cleaned, "628") {
				// Add + to make it +628xx
				cleaned = "+" + cleaned
			}
		} else {
			// Default to Indonesian country code if number looks local
			if len(cleaned) >= 9 && !strings.HasPrefix(cleaned, "62") {
				cleaned = "+62" + cleaned
			} else {
				cleaned = "+" + cleaned
			}
		}
	}
	
	if !phoneRegex.MatchString(cleaned) {
		return PhoneNumber{}, fmt.Errorf("invalid phone number format: %s", number)
	}
	
	if len(cleaned) > 16 { // E.164 max length is 15 digits + plus sign
		return PhoneNumber{}, fmt.Errorf("phone number too long: %s", cleaned)
	}
	
	return PhoneNumber{number: cleaned}, nil
}

// Value returns the phone number in E.164 format
func (p PhoneNumber) Value() string {
	return p.number
}

// String returns the phone number in E.164 format
func (p PhoneNumber) String() string {
	return p.number
}

// FormattedString returns the phone number in a human-readable format
func (p PhoneNumber) FormattedString() string {
	if len(p.number) == 0 {
		return ""
	}
	
	// Simple formatting for Indonesian numbers (+628xxx)
	if strings.HasPrefix(p.number, "+628") && len(p.number) >= 12 {
		// +628123456789 -> +62 812-3456-789
		return fmt.Sprintf("%s %s-%s-%s", 
			p.number[:3], 
			p.number[3:6], 
			p.number[6:10], 
			p.number[10:])
	}
	
	// Default formatting
	if len(p.number) > 4 {
		return fmt.Sprintf("%s %s", p.number[:4], p.number[4:])
	}
	
	return p.number
}

// CountryCode returns the country code (without +)
func (p PhoneNumber) CountryCode() string {
	if !strings.HasPrefix(p.number, "+") {
		return ""
	}
	
	// Indonesian numbers
	if strings.HasPrefix(p.number, "+62") {
		return "62"
	}
	
	// US/Canada
	if strings.HasPrefix(p.number, "+1") {
		return "1"
	}
	
	// Extract country code (1-4 digits after +)
	for i := 2; i <= 5 && i <= len(p.number); i++ {
		return p.number[1:i]
	}
	
	return ""
}

// IsIndonesian returns true if this is an Indonesian phone number
func (p PhoneNumber) IsIndonesian() bool {
	return strings.HasPrefix(p.number, "+62")
}

// Equal checks if this phone number equals another
func (p PhoneNumber) Equal(other PhoneNumber) bool {
	return p.number == other.number
}

// IsEmpty returns true if the phone number is empty
func (p PhoneNumber) IsEmpty() bool {
	return p.number == ""
}