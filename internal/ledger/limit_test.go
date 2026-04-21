package ledger

import (
	"strings"
	"testing"
)

func makeLimitBook(t *testing.T) *Book {
	t.Helper()
	b := makeBook(t)
	tx, err := NewTransaction("tx-limit-1", []Leg{
		{AccountID: "acc-a", Currency: "USD", Type: EntryDebit, Amount: 500_00},
		{AccountID: "acc-b", Currency: "USD", Type: EntryCredit, Amount: 500_00},
	})
	if err != nil {
		t.Fatalf("NewTransaction: %v", err)
	}
	if err := b.Post(tx); err != nil {
		t.Fatalf("Post: %v", err)
	}
	return b
}

func TestLimitIndex_Add_Valid(t *testing.T) {
	idx := NewLimitIndex()
	err := idx.Add(AccountLimit{AccountID: "acc-a", Currency: "USD", Type: LimitDebit, Max: 1000_00})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := len(idx.All()); got != 1 {
		t.Fatalf("expected 1 limit, got %d", got)
	}
}

func TestLimitIndex_Add_Invalid(t *testing.T) {
	idx := NewLimitIndex()
	cases := []AccountLimit{
		{AccountID: "", Currency: "USD", Type: LimitDebit, Max: 100},
		{AccountID: "acc-a", Currency: "", Type: LimitDebit, Max: 100},
		{AccountID: "acc-a", Currency: "USD", Type: "unknown", Max: 100},
		{AccountID: "acc-a", Currency: "USD", Type: LimitDebit, Max: 0},
	}
	for _, c := range cases {
		if err := idx.Add(c); err == nil {
			t.Errorf("expected error for limit %+v", c)
		}
	}
}

func TestCheckLimits_NoViolation(t *testing.T) {
	b := makeLimitBook(t)
	idx := NewLimitIndex()
	_ = idx.Add(AccountLimit{AccountID: "acc-a", Currency: "USD", Type: LimitDebit, Max: 1000_00})
	v, err := CheckLimits(idx, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(v) != 0 {
		t.Fatalf("expected no violations, got %d", len(v))
	}
}

func TestCheckLimits_Violation(t *testing.T) {
	b := makeLimitBook(t)
	idx := NewLimitIndex()
	_ = idx.Add(AccountLimit{AccountID: "acc-a", Currency: "USD", Type: LimitDebit, Max: 100_00})
	v, err := CheckLimits(idx, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].AccountID != "acc-a" {
		t.Errorf("expected acc-a, got %s", v[0].AccountID)
	}
}

func TestCheckLimits_NilInputs(t *testing.T) {
	idx := NewLimitIndex()
	b := &Book{}
	if _, err := CheckLimits(nil, b); err == nil {
		t.Error("expected error for nil index")
	}
	if _, err := CheckLimits(idx, nil); err == nil {
		t.Error("expected error for nil book")
	}
}

func TestExportLimitViolationsCSV(t *testing.T) {
	v := []LimitViolation{
		{AccountID: "acc-a", Currency: "USD", Type: LimitDebit, Limit: 100_00, Actual: 500_00},
	}
	out := ExportLimitViolationsCSV(v)
	if !strings.Contains(out, "acc-a") {
		t.Error("CSV missing account ID")
	}
	if !strings.HasPrefix(out, "account_id,") {
		t.Error("CSV missing header")
	}
}

func TestLimitSummary_Empty(t *testing.T) {
	out := LimitSummary(nil)
	if !strings.Contains(out, "No limit violations") {
		t.Errorf("unexpected summary: %s", out)
	}
}
