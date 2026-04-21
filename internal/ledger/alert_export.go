package ledger

import (
	"fmt"
	"strings"
)

// ExportAlertsCSV formats a slice of Alert values as a CSV string.
// Columns: triggered_at, severity, account_id, currency, balance, threshold, direction, message
func ExportAlertsCSV(alerts []Alert) string {
	if len(alerts) == 0 {
		return "triggered_at,severity,account_id,currency,balance,threshold,direction,message\n"
	}
	var sb strings.Builder
	sb.WriteString("triggered_at,severity,account_id,currency,balance,threshold,direction,message\n")
	for _, a := range alerts {
		direction := "above"
		if a.Rule.Below {
			direction = "below"
		}
		line := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s\n",
			a.Triggered.UTC().Format("2006-01-02T15:04:05Z"),
			escapeCSV(string(a.Rule.Severity)),
			escapeCSV(a.Rule.AccountID),
			escapeCSV(a.Rule.Currency),
			formatAmount(a.Balance, a.Rule.Currency),
			formatAmount(a.Rule.Threshold, a.Rule.Currency),
			direction,
			escapeCSV(a.Rule.Message),
		)
		sb.WriteString(line)
	}
	return sb.String()
}

// escapeCSV wraps a field in double-quotes if it contains a comma, quote, or newline.
func escapeCSV(s string) string {
	if strings.ContainsAny(s, ",\"\n") {
		s = strings.ReplaceAll(s, "\"", "\"\"")
		return "\"" + s + "\""
	}
	return s
}
