package ledger

import (
	"errors"
	"fmt"
	"time"
)

// MarginRule defines a minimum margin (as a fraction) between two accounts.
// For example, assets / liabilities must remain above a threshold.
type MarginRule struct {
	ID          string
	Numerator   string  // account ID
	Denominator string  // account ID
	MinRatio    float64 // e.g. 1.5 means numerator must be >= 1.5x denominator
	Currency    string
}

// MarginViolation records a rule breach at evaluation time.
type MarginViolation struct {
	Rule      MarginRule
	Ratio     float64
	Evaluated time.Time
	Message   string
}

func (v MarginViolation) String() string {
	return fmt.Sprintf("margin violation [%s]: ratio=%.4f < min=%.4f (%s)",
		v.Rule.ID, v.Ratio, v.Rule.MinRatio, v.Evaluated.Format(time.DateOnly))
}

// MarginIndex holds a collection of margin rules.
type MarginIndex struct {
	rules map[string]MarginRule
}

// NewMarginIndex returns an initialised MarginIndex.
func NewMarginIndex() *MarginIndex {
	return &MarginIndex{rules: make(map[string]MarginRule)}
}

// Add registers a margin rule. Returns an error for invalid input.
func (m *MarginIndex) Add(rule MarginRule) error {
	if rule.ID == "" {
		return errors.New("margin rule ID must not be empty")
	}
	if rule.Numerator == "" || rule.Denominator == "" {
		return errors.New("margin rule numerator and denominator must not be empty")
	}
	if rule.MinRatio <= 0 {
		return errors.New("margin rule MinRatio must be positive")
	}
	if _, exists := m.rules[rule.ID]; exists {
		return fmt.Errorf("margin rule %q already exists", rule.ID)
	}
	m.rules[rule.ID] = rule
	return nil
}

// Rules returns a copy of all registered rules.
func (m *MarginIndex) Rules() []MarginRule {
	out := make([]MarginRule, 0, len(m.rules))
	for _, r := range m.rules {
		out = append(out, r)
	}
	return out
}

// CheckMargins evaluates all margin rules against the provided book.
// It returns any violations found. A nil book returns an error.
func CheckMargins(idx *MarginIndex, book *Book) ([]MarginViolation, error) {
	if book == nil {
		return nil, errors.New("book must not be nil")
	}
	if idx == nil {
		return nil, nil
	}

	now := time.Now().UTC()
	var violations []MarginViolation

	for _, rule := range idx.rules {
		num, err := book.GetAccount(rule.Numerator)
		if err != nil {
			return nil, fmt.Errorf("margin rule %q: numerator account %q not found", rule.ID, rule.Numerator)
		}
		den, err := book.GetAccount(rule.Denominator)
		if err != nil {
			return nil, fmt.Errorf("margin rule %q: denominator account %q not found", rule.ID, rule.Denominator)
		}

		numBal := num.Balance(rule.Currency)
		denBal := den.Balance(rule.Currency)

		if denBal == 0 {
			continue // avoid division by zero; treat as satisfied
		}

		ratio := float64(numBal) / float64(denBal)
		if ratio < rule.MinRatio {
			violations = append(violations, MarginViolation{
				Rule:      rule,
				Ratio:     ratio,
				Evaluated: now,
				Message: fmt.Sprintf("ratio %.4f is below minimum %.4f", ratio, rule.MinRatio),
			})
		}
	}
	return violations, nil
}

// HasMarginViolations returns true when violations is non-empty.
func HasMarginViolations(violations []MarginViolation) bool {
	return len(violations) > 0
}
