package ledger

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// VelocityRule defines a maximum transaction volume (count and/or total amount)
// allowed for a given account within a rolling time window.
type VelocityRule struct {
	AccountID  string
	Window     time.Duration
	MaxCount   int     // 0 means no count limit
	MaxAmount  float64 // 0 means no amount limit
	Currency   string
}

// VelocityViolation describes a rule that was breached.
type VelocityViolation struct {
	Rule        VelocityRule
	ActualCount int
	ActualAmount float64
	Reason      string
}

func (v VelocityViolation) String() string {
	return fmt.Sprintf("velocity violation [%s]: %s (count=%d, amount=%.2f %s)",
		v.Rule.AccountID, v.Reason, v.ActualCount, v.ActualAmount, v.Rule.Currency)
}

// VelocityIndex stores velocity rules keyed by account ID.
type VelocityIndex struct {
	mu    sync.RWMutex
	rules map[string][]VelocityRule
}

// NewVelocityIndex creates an empty VelocityIndex.
func NewVelocityIndex() *VelocityIndex {
	return &VelocityIndex{rules: make(map[string][]VelocityRule)}
}

// Add registers a VelocityRule. Returns an error for invalid inputs.
func (vi *VelocityIndex) Add(r VelocityRule) error {
	if r.AccountID == "" {
		return errors.New("velocity: account ID must not be empty")
	}
	if r.Window <= 0 {
		return errors.New("velocity: window must be positive")
	}
	if r.MaxCount < 0 || r.MaxAmount < 0 {
		return errors.New("velocity: limits must be non-negative")
	}
	if r.MaxCount == 0 && r.MaxAmount == 0 {
		return errors.New("velocity: at least one of MaxCount or MaxAmount must be set")
	}
	vi.mu.Lock()
	defer vi.mu.Unlock()
	vi.rules[r.AccountID] = append(vi.rules[r.AccountID], r)
	return nil
}

// Rules returns a copy of all rules for the given account.
func (vi *VelocityIndex) Rules(accountID string) []VelocityRule {
	vi.mu.RLock()
	defer vi.mu.RUnlock()
	src := vi.rules[accountID]
	out := make([]VelocityRule, len(src))
	copy(out, src)
	return out
}

// CheckVelocity evaluates velocity rules against the transactions in a Book.
// Only transactions whose Timestamp falls within the rule window (relative to now)
// are considered. Returns all violations found.
func CheckVelocity(vi *VelocityIndex, b *Book, now time.Time) ([]VelocityViolation, error) {
	if vi == nil {
		return nil, errors.New("velocity: index must not be nil")
	}
	if b == nil {
		return nil, errors.New("velocity: book must not be nil")
	}

	vi.mu.RLock()
	defer vi.mu.RUnlock()

	var violations []VelocityViolation

	for accountID, rules := range vi.rules {
		for _, rule := range rules {
			windowStart := now.Add(-rule.Window)
			var count int
			var total float64

			for _, tx := range b.Transactions() {
				if tx.Timestamp.Before(windowStart) {
					continue
				}
				for _, leg := range tx.Legs {
					if leg.AccountID == accountID && leg.Currency == rule.Currency {
						count++
						total += leg.Amount
						break
					}
				}
			}

			if rule.MaxCount > 0 && count > rule.MaxCount {
				violations = append(violations, VelocityViolation{
					Rule:         rule,
					ActualCount:  count,
					ActualAmount: total,
					Reason:       fmt.Sprintf("count %d exceeds max %d", count, rule.MaxCount),
				})
			} else if rule.MaxAmount > 0 && total > rule.MaxAmount {
				violations = append(violations, VelocityViolation{
					Rule:         rule,
					ActualCount:  count,
					ActualAmount: total,
					Reason:       fmt.Sprintf("amount %.2f exceeds max %.2f", total, rule.MaxAmount),
				})
			}
		}
	}
	return violations, nil
}

// HasVelocityViolations returns true if any violations are present.
func HasVelocityViolations(violations []VelocityViolation) bool {
	return len(violations) > 0
}
