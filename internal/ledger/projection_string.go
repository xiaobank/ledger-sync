package ledger

import (
	"fmt"
	"strings"
)

// String returns a formatted table of projection entries.
func (r *ProjectionResult) String() string {
	if r == nil || len(r.Entries) == 0 {
		return "ProjectionResult: (empty)\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-12s  %-22s  %-6s  %12s\n", "AccountID", "Timestamp", "Cur", "Balance"))
	sb.WriteString(strings.Repeat("-", 58) + "\n")
	for _, e := range r.Entries {
		for cur, bal := range e.Balance {
			sb.WriteString(fmt.Sprintf("%-12s  %-22s  %-6s  %12.4f\n",
				e.AccountID,
				e.At.Format("2006-01-02 15:04:05"),
				cur,
				bal,
			))
		}
	}
	return sb.String()
}

// String returns a one-line summary of a ProjectionSummary.
func (s *ProjectionSummary) String() string {
	if s == nil {
		return "ProjectionSummary: (nil)\n"
	}
	return fmt.Sprintf("Projection[%s %s] steps=%d final=%.4f peak=%.4f trough=%.4f\n",
		s.AccountID, s.Currency, s.Steps, s.Final, s.Peak, s.Trough)
}
