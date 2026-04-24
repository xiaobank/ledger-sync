package ledger

import (
	"fmt"
	"strings"
	"time"
)

// JournalEntry represents a single line in the journal ledger view.
type JournalEntry struct {
	Date        time.Time
	TxID        string
	Description string
	Account     string
	Debit       int64
	Credit      int64
	Currency    string
}

// Journal holds an ordered list of journal entries derived from a Book.
type Journal struct {
	Entries []JournalEntry
}

// BuildJournal constructs a Journal from all transactions in the given Book.
// Entries are ordered by transaction timestamp.
func BuildJournal(b *Book) (*Journal, error) {
	if b == nil {
		return nil, fmt.Errorf("journal: book must not be nil")
	}

	var entries []JournalEntry

	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, tx := range b.transactions {
		for _, leg := range tx.Legs {
			je := JournalEntry{
				Date:        tx.Timestamp,
				TxID:        tx.ID,
				Description: tx.Description,
				Account:     leg.AccountID,
				Currency:    leg.Currency,
			}
			if leg.Type == EntryDebit {
				je.Debit = leg.Amount
			} else {
				je.Credit = leg.Amount
			}
			entries = append(entries, je)
		}
	}

	// stable sort by date then txID
	sortJournalEntries(entries)

	return &Journal{Entries: entries}, nil
}

// FilterByAccount returns a new Journal containing only entries for the given account.
func (j *Journal) FilterByAccount(accountID string) *Journal {
	var out []JournalEntry
	for _, e := range j.Entries {
		if strings.EqualFold(e.Account, accountID) {
			out = append(out, e)
		}
	}
	return &Journal{Entries: out}
}

// FilterByDateRange returns a new Journal with entries within [from, to] inclusive.
func (j *Journal) FilterByDateRange(from, to time.Time) *Journal {
	var out []JournalEntry
	for _, e := range j.Entries {
		if !e.Date.Before(from) && !e.Date.After(to) {
			out = append(out, e)
		}
	}
	return &Journal{Entries: out}
}

// sortJournalEntries sorts in-place by Date asc, then TxID asc.
func sortJournalEntries(entries []JournalEntry) {
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0; j-- {
			a, b := entries[j-1], entries[j]
			if a.Date.After(b.Date) || (a.Date.Equal(b.Date) && a.TxID > b.TxID) {
				entries[j-1], entries[j] = entries[j], entries[j-1]
			} else {
				break
			}
		}
	}
}
