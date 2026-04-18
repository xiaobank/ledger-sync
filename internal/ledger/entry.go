package ledger

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// EntryType represents the side of a ledger entry.
type EntryType string

const (
	Debit  EntryType = "debit"
	Credit EntryType = "credit"
)

// Entry represents a single line in a double-entry ledger.
type Entry struct {
	ID          string    `json:"id"`
	TransactionID string  `json:"transaction_id"`
	Account     string    `json:"account"`
	Type        EntryType `json:"type"`
	AmountCents int64     `json:"amount_cents"`
	Currency    string    `json:"currency"`
	CreatedAt   time.Time `json:"created_at"`
}

// Transaction groups a balanced set of entries.
type Transaction struct {
	ID        string    `json:"id"`
	Reference string    `json:"reference"`
	Entries   []Entry   `json:"entries"`
	CreatedAt time.Time `json:"created_at"`
}

// NewTransaction creates a validated double-entry transaction.
func NewTransaction(reference string, entries []Entry) (*Transaction, error) {
	if err := validateBalance(entries); err != nil {
		return nil, err
	}
	txID := uuid.NewString()
	now := time.Now().UTC()
	for i := range entries {
		if entries[i].ID == "" {
			entries[i].ID = uuid.NewString()
		}
		entries[i].TransactionID = txID
		entries[i].CreatedAt = now
	}
	return &Transaction{
		ID:        txID,
		Reference: reference,
		Entries:   entries,
		CreatedAt: now,
	}, nil
}

// validateBalance ensures debits equal credits within the same currency.
func validateBalance(entries []Entry) error {
	balance := make(map[string]int64)
	for _, e := range entries {
		if e.AmountCents <= 0 {
			return errors.New("amount_cents must be positive")
		}
		switch e.Type {
		case Debit:
			balance[e.Currency] += e.AmountCents
		case Credit:
			balance[e.Currency] -= e.AmountCents
		default:
			return errors.New("invalid entry type: " + string(e.Type))
		}
	}
	for currency, net := range balance {
		if net != 0 {
			return errors.New("transaction not balanced for currency: " + currency)
		}
	}
	return nil
}
