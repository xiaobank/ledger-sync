package ledger

import (
	"errors"
	"time"
)

// EntryType distinguishes debit from credit entries.
type EntryType string

const (
	Debit  EntryType = "debit"
	Credit EntryType = "credit"
)

// Entry is a single line in a double-entry transaction.
type Entry struct {
	AccountID string
	Type      EntryType
	Amount    int64  // minor units (e.g. cents)
	Currency  string // ISO 4217
}

// Transaction groups balanced entries under a single ID.
type Transaction struct {
	ID        string
	CreatedAt time.Time
	Entries   []Entry
}

// NewTransaction validates and returns a Transaction.
func NewTransaction(id string, entries []Entry) (*Transaction, error) {
	if id == "" {
		return nil, errors.New("transaction id must not be empty")
	}
	if len(entries) < 2 {
		return nil, errors.New("transaction must have at least two entries")
	}
	for _, e := range entries {
		if e.Amount <= 0 {
			return nil, errors.New("entry amount must be positive")
		}
		if e.Currency == "" {
			return nil, errors.New("entry currency must not be empty")
		}
		if e.AccountID == "" {
			return nil, errors.New("entry account id must not be empty")
		}
		switch e.Type {
		case Debit, Credit:
		default:
			return nil, errors.New("unknown entry type: " + string(e.Type))
		}
	}
	if err := validateBalance(entries); err != nil {
		return nil, err
	}
	return &Transaction{
		ID:        id,
		CreatedAt: time.Now().UTC(),
		Entries:   entries,
	}, nil
}

// validateBalance ensures debits equal credits for each currency.
func validateBalance(entries []Entry) error {
	type sums struct{ debit, credit int64 }
	totals := make(map[string]*sums)
	for _, e := range entries {
		if _, ok := totals[e.Currency]; !ok {
			totals[e.Currency] = &sums{}
		}
		if e.Type == Debit {
			totals[e.Currency].debit += e.Amount
		} else {
			totals[e.Currency].credit += e.Amount
		}
	}
	for currency, s := range totals {
		if s.debit != s.credit {
			return errors.New("transaction not balanced for currency: " + currency)
		}
	}
	return nil
}
