package ledger

import "fmt"

// BalanceReport holds the balance summary for all accounts in a Book.
type BalanceReport struct {
	Entries []BalanceEntry
}

// BalanceEntry represents a single account's balance summary.
type BalanceEntry struct {
	AccountID   string
	AccountName string
	Currency    string
	Debit       int64
	Credit      int64
	Net         int64
}

// String returns a human-readable representation of the report.
func (r *BalanceReport) String() string {
	s := "Balance Report\n"
	s += fmt.Sprintf("%-20s %-20s %-6s %12s %12s %12s\n", "ID", "Name", "CCY", "Debit", "Credit", "Net")
	s += fmt.Sprintf("%s\n", repeatChar('-', 84))
	for _, e := range r.Entries {
		s += fmt.Sprintf("%-20s %-20s %-6s %12d %12d %12d\n",
			e.AccountID, e.AccountName, e.Currency, e.Debit, e.Credit, e.Net)
	}
	return s
}

func repeatChar(c rune, n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = c
	}
	return string(b)
}

// GenerateBalanceReport computes debit/credit totals per account from posted transactions.
func (b *Book) GenerateBalanceReport() (*BalanceReport, error) {
	type key struct {
		accountID string
		currency  string
	}

	totals := make(map[key]*BalanceEntry)

	for _, tx := range b.transactions {
		for _, leg := range tx.Legs {
			acc, err := b.GetAccount(leg.AccountID)
			if err != nil {
				return nil, fmt.Errorf("report: %w", err)
			}
			k := key{accountID: leg.AccountID, currency: leg.Currency}
			if _, ok := totals[k]; !ok {
				totals[k] = &BalanceEntry{
					AccountID:   acc.ID,
					AccountName: acc.Name,
					Currency:    leg.Currency,
				}
			}
			entry := totals[k]
			if leg.Type == Debit {
				entry.Debit += leg.Amount
			} else {
				entry.Credit += leg.Amount
			}
			entry.Net = entry.Debit - entry.Credit
		}
	}

	report := &BalanceReport{}
	for _, e := range totals {
		report.Entries = append(report.Entries, *e)
	}
	return report, nil
}
