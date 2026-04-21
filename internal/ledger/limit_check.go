package ledger

import (
	"errors"
	"fmt"
)

// CheckLimits evaluates all limits in the index against the transactions
// recorded in book. It returns a slice of violations (empty means no breach).
func CheckLimits(idx *LimitIndex, book *Book) ([]LimitViolation, error) {
	if idx == nil {
		return nil, errors.New("CheckLimits: limit index must not be nil")
	}
	if book == nil {
		return nil, errors.New("CheckLimits: book must not be nil")
	}

	// Accumulate totals: map[accountID+":"+currency+":"+type] -> sum
	totals := make(map[string]int64)

	book.mu.RLock()
	for _, tx := range book.transactions {
		for _, leg := range tx.Legs {
			key := fmt.Sprintf("%s:%s:%s", leg.AccountID, leg.Currency, string(leg.Type))
			totals[key] += leg.Amount
		}
	}
	book.mu.RUnlock()

	var violations []LimitViolation
	for _, lim := range idx.All() {
		key := fmt.Sprintf("%s:%s:%s", lim.AccountID, lim.Currency, string(lim.Type))
		actual := totals[key]
		if actual > lim.Max {
			violations = append(violations, LimitViolation{
				AccountID: lim.AccountID,
				Currency:  lim.Currency,
				Type:      lim.Type,
				Limit:     lim.Max,
				Actual:    actual,
			})
		}
	}
	return violations, nil
}

// HasLimitViolations is a convenience wrapper that returns true when at least
// one violation exists.
func HasLimitViolations(idx *LimitIndex, book *Book) (bool, error) {
	v, err := CheckLimits(idx, book)
	if err != nil {
		return false, err
	}
	return len(v) > 0, nil
}
