package ledger

import (
	"fmt"
	"strings"
)

// String returns a formatted text table of the journal for debugging/display.
func (j *Journal) String() string {
	if j == nil || len(j.Entries) == 0 {
		return "Journal: (empty)\n"
	}

	var sb strings.Builder
	header := fmt.Sprintf("%-12s %-10s %-20s %-14s %10s %10s\n",
		"Date", "TxID", "Account", "Currency", "Debit", "Credit")
	sb.WriteString(header)
	sb.WriteString(repeatChar('-', len(header)-1) + "\n")

	for _, e := range j.Entries {
		debit := ""
		credit := ""
		if e.Debit > 0 {
			debit = formatAmount(e.Debit)
		}
		if e.Credit > 0 {
			credit = formatAmount(e.Credit)
		}
		line := fmt.Sprintf("%-12s %-10s %-20s %-14s %10s %10s\n",
			e.Date.Format("2006-01-02"),
			truncate(e.TxID, 10),
			truncate(e.Account, 20),
			e.Currency,
			debit,
			credit,
		)
		sb.WriteString(line)
	}

	return sb.String()
}
