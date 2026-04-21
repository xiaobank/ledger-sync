package ledger

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// RuleType identifies the kind of posting rule.
type RuleType string

const (
	RuleTypeMinBalance RuleType = "min_balance"
	RuleTypeMaxBalance RuleType = "max_balance"
	RuleTypeRequireTag RuleType = "require_tag"
)

// Rule defines a constraint that can be evaluated against a book.
type Rule struct {
	ID        string
	Type      RuleType
	AccountID string
	Currency  string
	Threshold int64
	TagKey    string
	CreatedAt time.Time
}

// RuleIndex stores and evaluates posting rules.
type RuleIndex struct {
	rules map[string]*Rule
}

// NewRuleIndex creates an empty RuleIndex.
func NewRuleIndex() *RuleIndex {
	return &RuleIndex{rules: make(map[string]*Rule)}
}

// Add registers a rule. Returns an error if the rule is invalid or duplicate.
func (ri *RuleIndex) Add(r Rule) error {
	if strings.TrimSpace(r.ID) == "" {
		return errors.New("rule: ID must not be empty")
	}
	if strings.TrimSpace(r.AccountID) == "" {
		return errors.New("rule: AccountID must not be empty")
	}
	switch r.Type {
	case RuleTypeMinBalance, RuleTypeMaxBalance:
		if strings.TrimSpace(r.Currency) == "" {
			return errors.New("rule: Currency required for balance rule")
		}
	case RuleTypeRequireTag:
		if strings.TrimSpace(r.TagKey) == "" {
			return errors.New("rule: TagKey required for require_tag rule")
		}
	default:
		return fmt.Errorf("rule: unknown type %q", r.Type)
	}
	if _, exists := ri.rules[r.ID]; exists {
		return fmt.Errorf("rule: duplicate ID %q", r.ID)
	}
	r.CreatedAt = time.Now()
	ri.rules[r.ID] = &r
	return nil
}

// Get retrieves a rule by ID.
func (ri *RuleIndex) Get(id string) (*Rule, bool) {
	r, ok := ri.rules[id]
	return r, ok
}

// All returns a slice of all registered rules.
func (ri *RuleIndex) All() []Rule {
	out := make([]Rule, 0, len(ri.rules))
	for _, r := range ri.rules {
		out = append(out, *r)
	}
	return out
}
