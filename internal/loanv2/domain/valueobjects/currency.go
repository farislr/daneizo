package valueobjects

import "fmt"

// Currency represents different currencies supported by the system
type Currency string

const (
	USD Currency = "USD"
	IDR Currency = "IDR"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
	JPY Currency = "JPY"
)

// String returns the string representation of the currency
func (c Currency) String() string {
	return string(c)
}

// IsValid checks if the currency is supported
func (c Currency) IsValid() bool {
	switch c {
	case USD, IDR, EUR, GBP, JPY:
		return true
	}
	return false
}

// Symbol returns the currency symbol
func (c Currency) Symbol() string {
	switch c {
	case USD:
		return "$"
	case IDR:
		return "Rp"
	case EUR:
		return "€"
	case GBP:
		return "£"
	case JPY:
		return "¥"
	}
	return string(c)
}

// NewCurrency creates a new Currency with validation
func NewCurrency(c string) (Currency, error) {
	currency := Currency(c)
	if !currency.IsValid() {
		return "", fmt.Errorf("invalid currency: %s", c)
	}
	return currency, nil
}