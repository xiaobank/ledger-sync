package ledger

import (
	"fmt"
	"strings"
)

// ExportBudgetReport renders a BudgetCheckResult as a plain-text table.
func ExportBudgetReport(result *BudgetCheckResult) string {
	if result == nil || len(result.Statuses) == 0 {
		return "No budget data to export.\n"
	}

	var sb strings.Builder
	header := fmt.Sprintf("%-12s %-10s %-12s %-12s %-12s %-10s %s\n",
		"Account", "Period", "Limit", "Spent", "Remaining", "Currency", "Status")
	sb.WriteString(header)
	sb.WriteString(repeatChar('-', len(header)-1) + "\n")

	for _, s := range result.Statuses {
		status := "OK"
		if s.Exceeded {
			status = "EXCEEDED"
		}
		line := fmt.Sprintf("%-12s %-10s %-12s %-12s %-12s %-10s %s\n",
			truncate(s.AccountID, 12),
			truncate(s.Period, 10),
			formatAmount(s.Limit, s.Currency),
			formatAmount(s.Spent, s.Currency),
			formatAmount(s.Remaining, s.Currency),
			s.Currency,
			status,
		)
		sb.WriteString(line)
	}
	return sb.String()
}

// ExportBudgetCSV renders a BudgetCheckResult as CSV.
func ExportBudgetCSV(result *BudgetCheckResult) string {
	if result == nil || len(result.Statuses) == 0 {
		return "account_id,period,currency,limit,spent,remaining,exceeded\n"
	}
	var sb strings.Builder
	sb.WriteString("account_id,period,currency,limit,spent,remaining,exceeded\n")
	for _, s := range result.Statuses {
		sb.WriteString(fmt.Sprintf("%s,%s,%s,%d,%d,%d,%t\n",
			s.AccountID, s.Period, s.Currency,
			s.Limit, s.Spent, s.Remaining, s.Exceeded,
		))
	}
	return sb.String()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
