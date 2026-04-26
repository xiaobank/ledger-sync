package ledger

import (
	"testing"
	"time"
)

func makeAllocationBook(t *testing.T) *Book {
	t.Helper()
	b, _ := NewBook("alloc-book")
	acc1, _ := NewAccount("src", "Source", AssetAccount)
	acc2, _ := NewAccount("tgt1", "Target One", ExpenseAccount)
	acc3, _ := NewAccount("tgt2", "Target Two", ExpenseAccount)
	_ = b.AddAccount(acc1)
	_ = b.AddAccount(acc2)
	_ = b.AddAccount(acc3)
	tx, _ := NewTransaction("tx-1", time.Now(), []Leg{
		{AccountID: "src", Amount: 1000, Currency: "USD", Type: DebitEntry},
		{AccountID: "tgt1", Amount: 1000, Currency: "USD", Type: CreditEntry},
	})
	_ = b.Post(tx)
	return b
}

func TestNewAllocationRule_Valid(t *testing.T) {
	targets := []AllocationTarget{
		{AccountID: "tgt1", Share: 0.6},
		{AccountID: "tgt2", Share: 0.4},
	}
	_, err := NewAllocationRule("src", targets)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNewAllocationRule_EmptySource(t *testing.T) {
	_, err := NewAllocationRule("", []AllocationTarget{{AccountID: "tgt1", Share: 1.0}})
	if err == nil {
		t.Fatal("expected error for empty source")
	}
}

func TestNewAllocationRule_SharesNotSumToOne(t *testing.T) {
	targets := []AllocationTarget{
		{AccountID: "tgt1", Share: 0.3},
		{AccountID: "tgt2", Share: 0.3},
	}
	_, err := NewAllocationRule("src", targets)
	if err == nil {
		t.Fatal("expected error when shares do not sum to 1")
	}
}

func TestNewAllocationRule_InvalidShare(t *testing.T) {
	_, err := NewAllocationRule("src", []AllocationTarget{{AccountID: "tgt1", Share: -0.5}})
	if err == nil {
		t.Fatal("expected error for negative share")
	}
}

func TestNewAllocationRule_NoTargets(t *testing.T) {
	_, err := NewAllocationRule("src", nil)
	if err == nil {
		t.Fatal("expected error for empty targets")
	}
}

func TestApplyAllocation_Valid(t *testing.T) {
	b := makeAllocationBook(t)
	rule, _ := NewAllocationRule("src", []AllocationTarget{
		{AccountID: "tgt1", Share: 0.5},
		{AccountID: "tgt2", Share: 0.5},
	})
	results, err := ApplyAllocation(b, rule, "USD", time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if len(results[0].Allocations) != 2 {
		t.Fatalf("expected 2 allocations, got %d", len(results[0].Allocations))
	}
	for _, a := range results[0].Allocations {
		if a.Amount != 500 {
			t.Errorf("expected amount 500, got %d for account %s", a.Amount, a.AccountID)
		}
	}
}

func TestApplyAllocation_NilBook(t *testing.T) {
	rule, _ := NewAllocationRule("src", []AllocationTarget{{AccountID: "tgt1", Share: 1.0}})
	_, err := ApplyAllocation(nil, rule, "USD", time.Now())
	if err == nil {
		t.Fatal("expected error for nil book")
	}
}

func TestApplyAllocation_NoCurrencyMatch(t *testing.T) {
	b := makeAllocationBook(t)
	rule, _ := NewAllocationRule("src", []AllocationTarget{{AccountID: "tgt1", Share: 1.0}})
	results, err := ApplyAllocation(b, rule, "EUR", time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results for non-matching currency, got %d", len(results))
	}
}
