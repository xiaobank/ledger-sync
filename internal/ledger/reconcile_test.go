package ledger

import (
	"testing"
)

func TestReconcile_AllMatch(t *testing.T) {
	b := makeBook(t)

	expected := map[string]int64{
		"cash":     0,
		"revenue":  0,
	}

	results, err := Reconcile(b, expected)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.Match {
			t.Errorf("expected match for account %q", r.AccountID)
		}
	}
	if HasDiscrepancies(results) {
		t.Error("expected no discrepancies")
	}
}

func TestReconcile_Mismatch(t *testing.T) {
	b := makeBook(t)

	expected := map[string]int64{
		"cash": 9999,
	}

	results, err := Reconcile(b, expected)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !HasDiscrepancies(results) {
		t.Error("expected discrepancies but found none")
	}
	if results[0].Match {
		t.Error("expected mismatch for cash account")
	}
}

func TestReconcile_NilBook(t *testing.T) {
	_, err := Reconcile(nil, map[string]int64{"cash": 0})
	if err == nil {
		t.Fatal("expected error for nil book")
	}
}

func TestReconcile_EmptyExpected(t *testing.T) {
	b := makeBook(t)
	_, err := Reconcile(b, map[string]int64{})
	if err == nil {
		t.Fatal("expected error for empty expected balances")
	}
}

func TestReconcile_UnknownAccount(t *testing.T) {
	b := makeBook(t)
	_, err := Reconcile(b, map[string]int64{"nonexistent": 0})
	if err == nil {
		t.Fatal("expected error for unknown account")
	}
}
