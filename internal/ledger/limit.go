package ledger

import (
	"errors"
	"fmt"
	"sync"
)

// LimitType defines whether a limit applies to debit or credit activity.
type LimitType string

const (
	LimitDebit  LimitType = "debit"
	LimitCredit LimitType = "credit"
)

// AccountLimit defines a maximum allowed cumulative amount for an account
// over a given currency and direction.
type AccountLimit struct {
	AccountID string
	Currency  string
	Type      LimitType
	Max       int64 // in minor units
}

// LimitViolation describes a single breach of an account limit.
type LimitViolation struct {
	AccountID string
	Currency  string
	Type      LimitType
	Limit     int64
	Actual    int64
}

func (v LimitViolation) String() string {
	return fmt.Sprintf("account %s exceeded %s limit for %s: limit=%d actual=%d",
		v.AccountID, v.Type, v.Currency, v.Limit, v.Actual)
}

// LimitIndex stores account limits and evaluates them against a Book.
type LimitIndex struct {
	mu     sync.RWMutex
	limits []AccountLimit
}

// NewLimitIndex creates an empty LimitIndex.
func NewLimitIndex() *LimitIndex {
	return &LimitIndex{}
}

// Add registers an account limit. Returns an error for invalid input.
func (li *LimitIndex) Add(l AccountLimit) error {
	if l.AccountID == "" {
		return errors.New("limit: account ID must not be empty")
	}
	if l.Currency == "" {
		return errors.New("limit: currency must not be empty")
	}
	if l.Type != LimitDebit && l.Type != LimitCredit {
		return fmt.Errorf("limit: unknown limit type %q", l.Type)
	}
	if l.Max <= 0 {
		return errors.New("limit: max must be positive")
	}
	li.mu.Lock()
	defer li.mu.Unlock()
	li.limits = append(li.limits, l)
	return nil
}

// All returns a copy of all registered limits.
func (li *LimitIndex) All() []AccountLimit {
	li.mu.RLock()
	defer li.mu.RUnlock()
	out := make([]AccountLimit, len(li.limits))
	copy(out, li.limits)
	return out
}
