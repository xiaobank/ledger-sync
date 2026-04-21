package ledger

import (
	"fmt"
	"strings"
)

// RuleSummary is a flat representation of a rule for export.
type RuleSummary struct {
	ID        string
	Type      string
	AccountID string
	Currency  string
	Threshold int64
	TagKey    string
}

// ExportRulesCSV serialises all rules in ri as a CSV string.
func ExportRulesCSV(ri *RuleIndex) (string, error) {
	if ri == nil {
		return "", fmt.Errorf("export_rules: index must not be nil")
	}

	var sb strings.Builder
	sb.WriteString("id,type,account_id,currency,threshold,tag_key\n")

	for _, r := range ri.All() {
		line := fmt.Sprintf("%s,%s,%s,%s,%d,%s\n",
			escapeCSV(r.ID),
			escapeCSV(string(r.Type)),
			escapeCSV(r.AccountID),
			escapeCSV(r.Currency),
			r.Threshold,
			escapeCSV(r.TagKey),
		)
		sb.WriteString(line)
	}

	return sb.String(), nil
}

// ExportViolationsCSV serialises a slice of RuleViolation as a CSV string.
func ExportViolationsCSV(violations []RuleViolation) string {
	var sb strings.Builder
	sb.WriteString("rule_id,account_id,message\n")
	for _, v := range violations {
		sb.WriteString(fmt.Sprintf("%s,%s,%s\n",
			escapeCSV(v.RuleID),
			escapeCSV(v.AccountID),
			escapeCSV(v.Message),
		))
	}
	return sb.String()
}
