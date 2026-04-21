package ledger

import (
	"strings"
	"testing"
)

func TestExportRulesCSV_Valid(t *testing.T) {
	ri := NewRuleIndex()
	_ = ri.Add(Rule{ID: "r1", Type: RuleTypeMinBalance, AccountID: "acc1", Currency: "USD", Threshold: 100})
	_ = ri.Add(Rule{ID: "r2", Type: RuleTypeRequireTag, AccountID: "acc2", TagKey: "invoice"})

	csv, err := ExportRulesCSV(ri)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(csv, "id,type,account_id") {
		t.Error("CSV missing header")
	}
	if !strings.Contains(csv, "r1") || !strings.Contains(csv, "r2") {
		t.Error("CSV missing rule rows")
	}
}

func TestExportRulesCSV_Empty(t *testing.T) {
	ri := NewRuleIndex()
	csv, err := ExportRulesCSV(ri)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(csv), "\n")
	if len(lines) != 1 {
		t.Errorf("expected header only, got %d lines", len(lines))
	}
}

func TestExportRulesCSV_NilIndex(t *testing.T) {
	_, err := ExportRulesCSV(nil)
	if err == nil {
		t.Error("expected error for nil index")
	}
}

func TestExportViolationsCSV(t *testing.T) {
	violations := []RuleViolation{
		{RuleID: "r1", AccountID: "acc1", Message: "balance too low"},
		{RuleID: "r2", AccountID: "acc2", Message: "no tagged transactions"},
	}
	csv := ExportViolationsCSV(violations)
	if !strings.HasPrefix(csv, "rule_id,account_id,message") {
		t.Error("missing header")
	}
	if !strings.Contains(csv, "r1") || !strings.Contains(csv, "r2") {
		t.Error("missing violation rows")
	}
}

func TestExportViolationsCSV_Empty(t *testing.T) {
	csv := ExportViolationsCSV(nil)
	lines := strings.Split(strings.TrimSpace(csv), "\n")
	if len(lines) != 1 {
		t.Errorf("expected header only, got %d lines", len(lines))
	}
}
