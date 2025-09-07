package errors

import (
	"fmt"
)

// DomainError represents a domain-specific error
type DomainError struct {
	message string
	code    string
}

// NewDomainError creates a new domain error
func NewDomainError(message string) *DomainError {
	return &DomainError{
		message: message,
		code:    "DOMAIN_ERROR",
	}
}

// NewDomainErrorWithCode creates a new domain error with a specific code
func NewDomainErrorWithCode(message, code string) *DomainError {
	return &DomainError{
		message: message,
		code:    code,
	}
}

// Error implements the error interface
func (e *DomainError) Error() string {
	return e.message
}

// Code returns the error code
func (e *DomainError) Code() string {
	return e.code
}

// IsDomainError checks if an error is a domain error
func IsDomainError(err error) bool {
	_, ok := err.(*DomainError)
	return ok
}

// ValidationError represents a validation error
type ValidationError struct {
	field   string
	message string
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		field:   field,
		message: message,
	}
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.field, e.message)
}

// Field returns the field that failed validation
func (e *ValidationError) Field() string {
	return e.field
}

// Message returns the validation message
func (e *ValidationError) Message() string {
	return e.message
}

// NotFoundError represents an entity not found error
type NotFoundError struct {
	entity string
	id     string
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(entity, id string) *NotFoundError {
	return &NotFoundError{
		entity: entity,
		id:     id,
	}
}

// Error implements the error interface
func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with id '%s' not found", e.entity, e.id)
}

// Entity returns the entity type that was not found
func (e *NotFoundError) Entity() string {
	return e.entity
}

// ID returns the ID that was not found
func (e *NotFoundError) ID() string {
	return e.id
}