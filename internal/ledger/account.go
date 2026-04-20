package ledger

import (
	"errors"
	"sync"
)

// AccountType represents whether an account is debit-normal or credit-normal.
type AccountType string

const (
	Asset     AccountType = "asset"
	Liability AccountType = "liability"
	Equity    AccountType = "equity"
	Revenue   AccountType = "revenue"
	Expense   AccountType = "expense"
)

// Account holds a named ledger account with a running balance per currency.
type Account struct {
	mu       sync.RWMutex
	ID       string
	Name     string
	Type     AccountType
	Balances map[string]int64 // currency -> minor units
}

// NewAccount creates a new Account with the given id, name and type.
func NewAccount(id, name string, t AccountType) (*Account, error) {
	if id == "" {
		return nil, errors.New("account id must not be empty")
	}
	if name == "" {
		return nil, errors.New("account name must not be empty")
	}
	switch t {
	case Asset, Liability, Equity, Revenue, Expense:
	default:
		return nil, errors.New("unknown account type")
	}
	return &Account{
		ID:       id,
		Name:     name,
		Type:     t,
		Balances: make(map[string]int64),
	}, nil
}

// Balance returns the balance for the given currency.
func (a *Account) Balance(currency string) int64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.Balances[currency]
}

// Currencies returns a snapshot of all currency codes that have a recorded
// balance on this account.
func (a *Account) Currencies() []string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	currencies := make([]string, 0, len(a.Balances))
	for c := range a.Balances {
		currencies = append(currencies, c)
	}
	return currencies
}

// apply adds delta (positive or negative minor units) to the account balance.
func (a *Account) apply(currency string, delta int64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Balances[currency] += delta
}
