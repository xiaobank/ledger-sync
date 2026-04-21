package ledger

import (
	"fmt"
	"strings"
)

// CurrencyRateRow represents a single row in a currency rate export.
type CurrencyRateRow struct {
	From string
	To   string
	Rate float64
}

// ExportRatesCSV returns a CSV string of all registered exchange rates
// from the provided CurrencyRegistry.
func ExportRatesCSV(reg *CurrencyRegistry) (string, error) {
	if reg == nil {
		return "", fmt.Errorf("currency registry is nil")
	}

	var sb strings.Builder
	sb.WriteString("from,to,rate\n")

	for key, rate := range reg.rates {
		parts := strings.SplitN(key, "->", 2)
		if len(parts) != 2 {
			continue
		}
		sb.WriteString(fmt.Sprintf("%s,%s,%.6f\n", parts[0], parts[1], rate))
	}

	return sb.String(), nil
}

// RateSummary returns a human-readable summary of all exchange rates
// in the registry.
func RateSummary(reg *CurrencyRegistry) string {
	if reg == nil || len(reg.rates) == 0 {
		return "No exchange rates registered.\n"
	}

	var sb strings.Builder
	sb.WriteString("Exchange Rates:\n")
	sb.WriteString(repeatChar('-', 30) + "\n")

	for key, rate := range reg.rates {
		parts := strings.SplitN(key, "->", 2)
		if len(parts) != 2 {
			continue
		}
		sb.WriteString(fmt.Sprintf("  %s -> %s : %.6f\n", parts[0], parts[1], rate))
	}

	return sb.String()
}
