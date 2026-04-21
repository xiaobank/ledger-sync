package ledger

import (
	"fmt"
	"strings"
)

// ExportForecastCSV serialises a ForecastResult to a CSV string.
// Columns: account_id, currency, projected, period_date
func ExportForecastCSV(r *ForecastResult) (string, error) {
	if r == nil {
		return "", fmt.Errorf("export forecast csv: result must not be nil")
	}

	var sb strings.Builder
	sb.WriteString("account_id,currency,projected,period_date\n")

	for _, e := range r.Entries {
		line := fmt.Sprintf("%s,%s,%s,%s\n",
			e.AccountID,
			e.Currency,
			formatAmount(e.Projected),
			e.At.Format("2006-01-02"),
		)
		sb.WriteString(line)
	}

	return sb.String(), nil
}

// ForecastSummary returns a human-readable summary table of the forecast.
func ForecastSummary(r *ForecastResult) string {
	if r == nil || len(r.Entries) == 0 {
		return "No forecast data available.\n"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Forecast generated at: %s\n", r.GeneratedAt.Format("2006-01-02 15:04:05")))
	sb.WriteString(repeatChar('-', 60) + "\n")
	sb.WriteString(fmt.Sprintf("%-20s %-8s %14s  %s\n", "Account", "CCY", "Projected", "Period Date"))
	sb.WriteString(repeatChar('-', 60) + "\n")

	for _, e := range r.Entries {
		sb.WriteString(fmt.Sprintf("%-20s %-8s %14s  %s\n",
			e.AccountID,
			e.Currency,
			formatAmount(e.Projected),
			e.At.Format("2006-01-02"),
		))
	}

	sb.WriteString(repeatChar('-', 60) + "\n")
	return sb.String()
}
