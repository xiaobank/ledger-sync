package ledger

import (
	"errors"
	"sync"
)

// Book is an in-memory collection of accounts that can post transactions.
type Book struct {
	mu       sync.RWMutex
	accounts map[string]*Account
}

// NewBook returns an empty Book.
func NewBook() *Book {
	return &Book{accounts: make(map[string]*Account)}
}

// AddAccount registers an account with the book.
func (b *Book) AddAccount(a *Account) error {
	if a == nil {
		return errors.New("account must not be nil")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, exists := b.accounts[a.ID]; exists {
		return errors.New("account already exists: " + a.ID)
	}
	b.accounts[a.ID] = a
	return nil
}

// GetAccount retrieves an account by id.
func (b *Book) GetAccount(id string) (*Account, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	a, ok := b.accounts[id]
	if !ok {
		return nil, errors.New("account not found: " + id)
	}
	return a, nil
}

// Accounts returns a snapshot of all accounts currently registered in the book.
func (b *Book) Accounts() []*Account {
	b.mu.RLock()
	defer b.mu.RUnlock()
	list := make([]*Account, 0, len(b.accounts))
	for _, a := range b.accounts {
		list = append(list, a)
	}
	return list
}

// Post validates and applies a Transaction to the relevant accounts.
func (b *Book) Post(tx *Transaction) error {
	if tx == nil {
		return errors.New("transaction must not be nil")
	}
	// Resolve all accounts first to fail fast before mutating state.
	b.mu.RLock()
	accounts := make([]*Account, len(tx.Entries))
	for i, e := range tx.Entries {
		a, ok := b.accounts[e.AccountID]
		if !ok {
			b.mu.RUnlock()
			return errors.New("account not found: " + e.AccountID)
		}
		accounts[i] = a
	}
	b.mu.RUnlock()

	for i, e := range tx.Entries {
		delta := e.Amount
		if e.Type == Credit {
			delta = -delta
		}
		accounts[i].apply(e.Currency, delta)
	}
	return nil
}
