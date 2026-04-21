package ledger

import (
	"fmt"
	"strings"
	"time"
)

// AlertSeverity represents the urgency level of an alert.
type AlertSeverity string

const (
	SeverityInfo    AlertSeverity = "INFO"
	SeverityWarning AlertSeverity = "WARNING"
	SeverityCritical AlertSeverity = "CRITICAL"
)

// AlertRule defines a threshold-based rule evaluated against account balances.
type AlertRule struct {
	AccountID string
	Currency  string
	Threshold int64 // in minor units
	Below     bool  // true = alert when balance < Threshold; false = alert when balance > Threshold
	Severity  AlertSeverity
	Message   string
}

// Alert is a triggered alert produced by evaluating rules against a book.
type Alert struct {
	Rule      AlertRule
	Balance   int64
	Triggered time.Time
}

// String returns a human-readable description of the alert.
func (a Alert) String() string {
	direction := "above"
	if a.Rule.Below {
		direction = "below"
	}
	return fmt.Sprintf("[%s] %s — account %s balance %s %s threshold %s: %s",
		a.Rule.Severity,
		a.Triggered.Format(time.RFC3339),
		a.Rule.AccountID,
		formatAmount(a.Balance, a.Rule.Currency),
		direction,
		formatAmount(a.Rule.Threshold, a.Rule.Currency),
		a.Rule.Message,
	)
}

// AlertEngine evaluates a set of rules against a Book.
type AlertEngine struct {
	rules []AlertRule
}

// NewAlertEngine creates an AlertEngine with no rules.
func NewAlertEngine() *AlertEngine {
	return &AlertEngine{}
}

// AddRule registers a rule with the engine. Returns an error if the rule is invalid.
func (e *AlertEngine) AddRule(r AlertRule) error {
	if strings.TrimSpace(r.AccountID) == "" {
		return fmt.Errorf("alert rule: accountID must not be empty")
	}
	if strings.TrimSpace(r.Currency) == "" {
		return fmt.Errorf("alert rule: currency must not be empty")
	}
	if r.Severity == "" {
		return fmt.Errorf("alert rule: severity must not be empty")
	}
	e.rules = append(e.rules, r)
	return nil
}

// Evaluate checks all rules against the provided Book and returns triggered alerts.
func (e *AlertEngine) Evaluate(b *Book) ([]Alert, error) {
	if b == nil {
		return nil, fmt.Errorf("alert engine: book must not be nil")
	}
	now := time.Now().UTC()
	var alerts []Alert
	for _, rule := range e.rules {
		acc, err := b.GetAccount(rule.AccountID)
		if err != nil {
			continue // unknown account; skip silently
		}
		bal := acc.Balance(rule.Currency)
		triggered := (rule.Below && bal < rule.Threshold) || (!rule.Below && bal > rule.Threshold)
		if triggered {
			alerts = append(alerts, Alert{
				Rule:      rule,
				Balance:   bal,
				Triggered: now,
			})
		}
	}
	return alerts, nil
}
