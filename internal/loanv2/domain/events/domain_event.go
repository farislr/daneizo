package events

import (
	"time"
	"github.com/google/uuid"
	"github.com/farislr/daneizo/internal/loanv2/domain/valueobjects"
)

// DomainEvent represents a domain event that can be published
type DomainEvent interface {
	EventName() string
	EventID() string
	AggregateID() string
	AggregateType() string
	EventVersion() int
	Timestamp() time.Time
	UserID() *valueobjects.UserID
	CorrelationID() string
	CausationID() *string
	Payload() interface{}
	Metadata() map[string]interface{}
}

// BaseEvent provides common functionality for all domain events
type BaseEvent struct {
	eventID       string
	aggregateID   string
	aggregateType string
	eventVersion  int
	timestamp     time.Time
	userID        *valueobjects.UserID
	correlationID string
	causationID   *string
	metadata      map[string]interface{}
}

// NewBaseEvent creates a new base event
func NewBaseEvent(
	aggregateID string,
	aggregateType string,
	userID *valueobjects.UserID,
	correlationID string,
	causationID *string,
) BaseEvent {
	return BaseEvent{
		eventID:       uuid.New().String(),
		aggregateID:   aggregateID,
		aggregateType: aggregateType,
		eventVersion:  1,
		timestamp:     time.Now(),
		userID:        userID,
		correlationID: correlationID,
		causationID:   causationID,
		metadata:      make(map[string]interface{}),
	}
}

// EventID returns the unique event identifier
func (e BaseEvent) EventID() string {
	return e.eventID
}

// AggregateID returns the aggregate identifier
func (e BaseEvent) AggregateID() string {
	return e.aggregateID
}

// AggregateType returns the aggregate type
func (e BaseEvent) AggregateType() string {
	return e.aggregateType
}

// EventVersion returns the event version
func (e BaseEvent) EventVersion() int {
	return e.eventVersion
}

// Timestamp returns the event timestamp
func (e BaseEvent) Timestamp() time.Time {
	return e.timestamp
}

// UserID returns the user who triggered the event
func (e BaseEvent) UserID() *valueobjects.UserID {
	return e.userID
}

// CorrelationID returns the correlation identifier
func (e BaseEvent) CorrelationID() string {
	return e.correlationID
}

// CausationID returns the causation identifier
func (e BaseEvent) CausationID() *string {
	return e.causationID
}

// Metadata returns the event metadata
func (e BaseEvent) Metadata() map[string]interface{} {
	return e.metadata
}

// AddMetadata adds metadata to the event
func (e *BaseEvent) AddMetadata(key string, value interface{}) {
	e.metadata[key] = value
}