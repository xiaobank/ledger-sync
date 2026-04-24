package ledger

import (
	"fmt"
	"strings"
)

// JournalSummary holds a human-readable summary line for a JournalEntry.
type JournalSummary struct {
	Date    string
	TxID    string
	Account string
	Debit   string
	Credit  string
}

// ExportJournalCSV renders the journal as a CSV string.
func ExportJournalCSV(j *Journal) (string, error) {
	if j == nil {
		return "", fmt.Errorf("journal export: journal must not be nil")
	}

	var sb strings.Builder
	sb.WriteString("date,tx_id,description,account,debit,credit,currency\n")

	for _, e := range j.Entries {
		line := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s\n",
			e.Date.Format("2006-01-02"),
			escapeCSV(e.TxID),
			escapeCSV(e.Description),
			escapeCSV(e.Account),
			formatAmount(e.Debit),
			formatAmount(e.Credit),
			escapeCSV(e.Currency),
		)
		sb.WriteString(line)
	}

	return sb.String(), nil
}

// SummariseJournal returns a slice of JournalSummary for display purposes.
func SummariseJournal(j *Journal) []JournalSummary {
	if j == nil {
		return nil
	}
	out := make([]JournalSummary, 0, len(j.Entries))
	for _, e := range j.Entries {
		out = append(out, JournalSummary{
			Date:    e.Date.Format("2006-01-02"),
			TxID:    e.TxID,
			Account: e.Account,
			Debit:   formatAmount(e.Debit),
			Credit:  formatAmount(e.Credit),
		})
	}
	return out
}
