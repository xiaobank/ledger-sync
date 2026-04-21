package ledger

import (
	"fmt"
	"strings"
)

// BudgetCheckResult holds the outcome of checking all budgets against a book.
type BudgetCheckResult struct {
	Statuses  []*BudgetStatus
	Violations []*BudgetStatus
}

// HasViolations returns true if any budget was exceeded.
func (r *BudgetCheckResult) HasViolations() bool {
	return len(r.Violations) > 0
}

// String returns a human-readable summary of the check result.
func (r *BudgetCheckResult) String() string {
	if len(r.Statuses) == 0 {
		return "BudgetCheck: no budgets evaluated"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("BudgetCheck: %d budget(s) evaluated, %d violation(s)\n",
		len(r.Statuses), len(r.Violations)))
	for _, s := range r.Statuses {
		status := "OK"
		if s.Exceeded {
			status = "EXCEEDED"
		}
		sb.WriteString(fmt.Sprintf("  [%s] account=%s period=%s spent=%s limit=%s remaining=%s\n",
			status,
			s.AccountID,
			s.Period,
			formatAmount(s.Spent, s.Currency),
			formatAmount(s.Limit, s.Currency),
			formatAmount(s.Remaining, s.Currency),
		))
	}
	return sb.String()
}

// CheckBudgets evaluates all budgets in the index against the actual balances
// in the provided book for the given period.
func CheckBudgets(index *BudgetIndex, book *Book, period string) (*BudgetCheckResult, error) {
	if index == nil {
		return nil, fmt.Errorf("budget check: index must not be nil")
	}
	if book == nil {
		return nil, fmt.Errorf("budget check: book must not be nil")
	}

	index.mu.RLock()
	defer index.mu.RUnlock()

	result := &BudgetCheckResult{}

	for key, budget := range index.budgets {
		_ = key
		if budget.Period != period {
			continue
		}
		acct, err := book.GetAccount(budget.AccountID)
		if err != nil {
			return nil, fmt.Errorf("budget check: %w", err)
		}
		spent := acct.Balance(budget.Currency)
		if spent < 0 {
			spent = -spent
		}
		status, err := index.Evaluate(budget.AccountID, budget.Period, spent)
		if err != nil {
			return nil, err
		}
		result.Statuses = append(result.Statuses, status)
		if status.Exceeded {
			result.Violations = append(result.Violations, status)
		}
	}
	return result, nil
}
