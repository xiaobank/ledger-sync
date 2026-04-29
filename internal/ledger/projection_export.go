package ledger

import (
	"fmt"
	"strings"
)

// ProjectionSummary holds a human-readable summary of a projection result.
type ProjectionSummary struct {
	AccountID string
	Currency  string
	Steps     int
	Final     float64
	Peak      float64
	Trough    float64
}

// SummariseProjection returns a summary of the projection result.
func SummariseProjection(r *ProjectionResult, currency string) *ProjectionSummary {
	if r == nil || len(r.Entries) == 0 {
		return nil
	}
	s := &ProjectionSummary{
		AccountID: r.Entries[0].AccountID,
		Currency:  currency,
		Steps:     len(r.Entries),
		Peak:      r.Entries[0].Balance[currency],
		Trough:    r.Entries[0].Balance[currency],
	}
	for _, e := range r.Entries {
		v := e.Balance[currency]
		if v > s.Peak {
			s.Peak = v
		}
		if v < s.Trough {
			s.Trough = v
		}
		s.Final = v
	}
	return s
}

// ExportProjectionCSV serialises a ProjectionResult to CSV format.
func ExportProjectionCSV(r *ProjectionResult) (string, error) {
	if r == nil {
		return "", fmt.Errorf("projection export: result must not be nil")
	}
	var sb strings.Builder
	sb.WriteString("account_id,timestamp,currency,balance\n")
	for _, e := range r.Entries {
		for cur, bal := range e.Balance {
			fmt.Fprintf(&sb, "%s,%s,%s,%.4f\n",
				e.AccountID,
				e.At.Format("2006-01-02T15:04:05Z"),
				cur,
				bal,
			)
		}
	}
	return sb.String(), nil
}
