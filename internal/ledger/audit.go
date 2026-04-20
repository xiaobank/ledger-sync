package ledger

import (
	"fmt"
	"time"
)

// AuditEventType classifies the kind of audit event.
type AuditEventType string

const (
	AuditEventPost          AuditEventType = "POST"
	AuditEventAddAccount    AuditEventType = "ADD_ACCOUNT"
	AuditEventReconcile     AuditEventType = "RECONCILE"
	AuditEventSnapshot      AuditEventType = "SNAPSHOT"
)

// AuditEvent represents a single recorded action in the ledger.
type AuditEvent struct {
	Timestamp time.Time
	EventType AuditEventType
	ActorID   string
	Details   string
}

// AuditLog holds an ordered list of audit events.
type AuditLog struct {
	events []AuditEvent
}

// NewAuditLog creates an empty AuditLog.
func NewAuditLog() *AuditLog {
	return &AuditLog{}
}

// Record appends a new event to the audit log.
func (a *AuditLog) Record(eventType AuditEventType, actorID, details string) {
	a.events = append(a.events, AuditEvent{
		Timestamp: time.Now().UTC(),
		EventType: eventType,
		ActorID:   actorID,
		Details:   details,
	})
}

// Events returns a copy of all recorded audit events.
func (a *AuditLog) Events() []AuditEvent {
	result := make([]AuditEvent, len(a.events))
	copy(result, a.events)
	return result
}

// FilterByType returns all events matching the given type.
func (a *AuditLog) FilterByType(eventType AuditEventType) []AuditEvent {
	var filtered []AuditEvent
	for _, e := range a.events {
		if e.EventType == eventType {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

// FilterByActor returns all events recorded for a given actorID.
func (a *AuditLog) FilterByActor(actorID string) []AuditEvent {
	var filtered []AuditEvent
	for _, e := range a.events {
		if e.ActorID == actorID {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

// String returns a human-readable summary of the audit log.
func (a *AuditLog) String() string {
	if len(a.events) == 0 {
		return "audit log: no events recorded\n"
	}
	out := fmt.Sprintf("audit log: %d event(s)\n", len(a.events))
	for _, e := range a.events {
		out += fmt.Sprintf("  [%s] %s actor=%s %s\n",
			e.Timestamp.Format(time.RFC3339),
			e.EventType,
			e.ActorID,
			e.Details,
		)
	}
	return out
}
