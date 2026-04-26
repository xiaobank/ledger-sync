package ledger

import (
	"errors"
	"fmt"
	"time"
)

// AllocationRule defines how a transaction amount should be split across accounts.
type AllocationRule struct {
	SourceAccount string
	Targets       []AllocationTarget
}

// AllocationTarget specifies a destination account and its share (0 < Share <= 1).
type AllocationTarget struct {
	AccountID string
	Share     float64 // fraction of total, e.g. 0.5 for 50%
}

// AllocationResult holds the derived transactions produced by an allocation.
type AllocationResult struct {
	SourceTxID  string
	Allocations []AllocatedEntry
}

// AllocatedEntry represents a single split entry.
type AllocatedEntry struct {
	AccountID string
	Amount    int64
	Currency  string
}

// NewAllocationRule validates and returns an AllocationRule.
func NewAllocationRule(source string, targets []AllocationTarget) (AllocationRule, error) {
	if source == "" {
		return AllocationRule{}, errors.New("allocation: source account must not be empty")
	}
	if len(targets) == 0 {
		return AllocationRule{}, errors.New("allocation: at least one target required")
	}
	var total float64
	for _, t := range targets {
		if t.AccountID == "" {
			return AllocationRule{}, errors.New("allocation: target account ID must not be empty")
		}
		if t.Share <= 0 || t.Share > 1 {
			return AllocationRule{}, fmt.Errorf("allocation: share for %q must be in (0,1], got %v", t.AccountID, t.Share)
		}
		total += t.Share
	}
	if total < 0.9999 || total > 1.0001 {
		return AllocationRule{}, fmt.Errorf("allocation: shares must sum to 1.0, got %v", total)
	}
	return AllocationRule{SourceAccount: source, Targets: targets}, nil
}

// ApplyAllocation distributes the debit legs of matching transactions in the book
// according to the rule, posting synthetic credit/debit pairs for each target.
func ApplyAllocation(b *Book, rule AllocationRule, currency string, asOf time.Time) ([]AllocationResult, error) {
	if b == nil {
		return nil, errors.New("allocation: book must not be nil")
	}
	var results []AllocationResult
	for _, tx := range b.Transactions() {
		for _, leg := range tx.Legs {
			if leg.AccountID != rule.SourceAccount || leg.Currency != currency || leg.Type != DebitEntry {
				continue
			}
			result := AllocationResult{SourceTxID: tx.ID}
			for _, target := range rule.Targets {
				share := int64(float64(leg.Amount) * target.Share)
				result.Allocations = append(result.Allocations, AllocatedEntry{
					AccountID: target.AccountID,
					Amount:    share,
					Currency:  currency,
				})
			}
			results = append(results, result)
		}
	}
	return results, nil
}
