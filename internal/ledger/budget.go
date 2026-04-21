package ledger

import (
	"errors"
	"fmt"
	"sync"
)

// Budget represents a spending or allocation limit for an account over a period.
type Budget struct {
	AccountID string
	Currency  string
	Limit     int64 // in minor units (e.g. cents)
	Period    string // e.g. "2024-01", "2024-Q1", "2024"
}

// BudgetStatus holds the evaluated state of a budget.
type BudgetStatus struct {
	Budget
	Spent     int64
	Remaining int64
	Exceeded  bool
}

// BudgetIndex manages a collection of budgets.
type BudgetIndex struct {
	mu      sync.RWMutex
	budgets map[string]*Budget // key: accountID+":"+period
}

// NewBudgetIndex creates an empty BudgetIndex.
func NewBudgetIndex() *BudgetIndex {
	return &BudgetIndex{
		budgets: make(map[string]*Budget),
	}
}

func budgetKey(accountID, period string) string {
	return accountID + ":" + period
}

// Set adds or replaces a budget entry.
func (b *BudgetIndex) Set(budget Budget) error {
	if budget.AccountID == "" {
		return errors.New("budget: accountID must not be empty")
	}
	if budget.Period == "" {
		return errors.New("budget: period must not be empty")
	}
	if budget.Limit <= 0 {
		return errors.New("budget: limit must be positive")
	}
	if budget.Currency == "" {
		return errors.New("budget: currency must not be empty")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	copy := budget
	b.budgets[budgetKey(budget.AccountID, budget.Period)] = &copy
	return nil
}

// Get retrieves a budget for a given account and period.
func (b *BudgetIndex) Get(accountID, period string) (*Budget, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	budget, ok := b.budgets[budgetKey(accountID, period)]
	if !ok {
		return nil, false
	}
	copy := *budget
	return &copy, true
}

// Evaluate computes the BudgetStatus for an account/period given actual spending.
func (b *BudgetIndex) Evaluate(accountID, period string, spent int64) (*BudgetStatus, error) {
	budget, ok := b.Get(accountID, period)
	if !ok {
		return nil, fmt.Errorf("budget: no budget found for account %q period %q", accountID, period)
	}
	remaining := budget.Limit - spent
	return &BudgetStatus{
		Budget:    *budget,
		Spent:     spent,
		Remaining: remaining,
		Exceeded:  spent > budget.Limit,
	}, nil
}
