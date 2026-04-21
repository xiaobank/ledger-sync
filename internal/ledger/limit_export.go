package ledger

import (
	"bytes"
	"fmt"
)

// ExportLimitViolationsCSV formats a slice of LimitViolation values as a
// CSV string with a header row.
func ExportLimitViolationsCSV(violations []LimitViolation) string {
	var buf bytes.Buffer
	buf.WriteString("account_id,currency,type,limit,actual,excess\n")
	for _, v := range violations {
		excess := v.Actual - v.Limit
		fmt.Fprintf(&buf, "%s,%s,%s,%s,%s,%s\n",
			v.AccountID,
			v.Currency,
			string(v.Type),
			formatAmount(v.Limit),
			formatAmount(v.Actual),
			formatAmount(excess),
		)
	}
	return buf.String()
}

// LimitSummary returns a human-readable multi-line summary of violations.
func LimitSummary(violations []LimitViolation) string {
	if len(violations) == 0 {
		return "No limit violations.\n"
	}
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%d limit violation(s):\n", len(violations))
	for _, v := range violations {
		fmt.Fprintf(&buf, "  - %s\n", v.String())
	}
	return buf.String()
}
