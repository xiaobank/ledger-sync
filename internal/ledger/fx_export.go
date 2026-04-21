package ledger

import (
	"fmt"
	"strings"
)

// FXConversionSummary holds a human-readable summary of an FXConversion.
type FXConversionSummary struct {
	From      string
	To        string
	Rate      float64
	AppliedAt string
}

// SummariseFXConversions converts a slice of FXConversion into FXConversionSummary records.
func SummariseFXConversions(convs []FXConversion) []FXConversionSummary {
	out := make([]FXConversionSummary, 0, len(convs))
	for _, c := range convs {
		out = append(out, FXConversionSummary{
			From:      strings.ToUpper(c.FromCurrency),
			To:        strings.ToUpper(c.ToCurrency),
			Rate:      c.Rate,
			AppliedAt: c.AppliedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	return out
}

// ExportFXConversionsCSV renders a slice of FXConversion as a CSV string.
func ExportFXConversionsCSV(convs []FXConversion) string {
	if len(convs) == 0 {
		return "from,to,rate,applied_at\n"
	}
	var sb strings.Builder
	sb.WriteString("from,to,rate,applied_at\n")
	for _, c := range convs {
		sb.WriteString(fmt.Sprintf("%s,%s,%.6f,%s\n",
			strings.ToUpper(c.FromCurrency),
			strings.ToUpper(c.ToCurrency),
			c.Rate,
			c.AppliedAt.Format("2006-01-02T15:04:05Z"),
		))
	}
	return sb.String()
}
