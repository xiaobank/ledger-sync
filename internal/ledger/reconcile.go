package ledger

import (
	"errors"
	"fmt"
)

// ReconcileResult holds the outcome of a reconciliation check.
type ReconcileResult struct {
	AccountID string
	Expected  int64
	Actual    int64
	Match     bool
}

// Reconcile compares expected balances against actual account balances in a Book.
// expectedBalances is a map of accountID -> expected balance in minor units.
func Reconcile(b *Book, expectedBalances map[string]int64) ([]ReconcileResult, error) {
	if b == nil {
		return nil, errors.New("reconcile: book must not be nil")
	}
	if len(expectedBalances) == 0 {
		return nil, errors.New("reconcile: expected balances must not be empty")
	}

	results := make([]ReconcileResult, 0, len(expectedBalances))

	for id, expected := range expectedBalances {
		acc, err := b.GetAccount(id)
		if err != nil {
			return nil, fmt.Errorf("reconcile: account %q not found in book", id)
		}
		actual := acc.Balance()
		results = append(results, ReconcileResult{
			AccountID: id,
			Expected:  expected,
			Actual:    actual,
			Match:     expected == actual,
		})
	}

	return results, nil
}

// HasDiscrepancies returns true if any result in the slice does not match.
func HasDiscrepancies(results []ReconcileResult) bool {
	for _, r := range results {
		if !r.Match {
			return true
		}
	}
	return false
}
