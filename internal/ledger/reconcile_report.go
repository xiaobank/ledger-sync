package ledger

import (
	"fmt"
	"strings"
)

// ReconcileReport formats reconciliation results as a human-readable string.
type ReconcileReport struct {
	Results []ReconcileResult
}

// NewReconcileReport creates a ReconcileReport from a slice of results.
func NewReconcileReport(results []ReconcileResult) *ReconcileReport {
	return &ReconcileReport{Results: results}
}

// String renders the reconciliation report.
func (r *ReconcileReport) String() string {
	var sb strings.Builder
	title := "Reconciliation Report"
	width := 48
	sb.WriteString(repeatChar('=', width) + "\n")
	sb.WriteString(fmt.Sprintf("%-*s\n", width, title))
	sb.WriteString(repeatChar('=', width) + "\n")
	sb.WriteString(fmt.Sprintf("%-20s %10s %10s %6s\n", "Account", "Expected", "Actual", "OK?"))
	sb.WriteString(repeatChar('-', width) + "\n")

	for _, res := range r.Results {
		ok := "YES"
		if !res.Match {
			ok = "NO"
		}
		sb.WriteString(fmt.Sprintf("%-20s %10d %10d %6s\n",
			res.AccountID, res.Expected, res.Actual, ok))
	}

	sb.WriteString(repeatChar('=', width) + "\n")
	if HasDiscrepancies(r.Results) {
		sb.WriteString("STATUS: DISCREPANCIES FOUND\n")
	} else {
		sb.WriteString("STATUS: ALL BALANCED\n")
	}
	return sb.String()
}
