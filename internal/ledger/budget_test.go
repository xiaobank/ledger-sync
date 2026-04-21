package ledger

import (
	"testing"
)

func TestBudgetIndex_Set_And_Get(t *testing.T) {
	idx := NewBudgetIndex()
	err := idx.Set(Budget{AccountID: "acc-1", Currency: "USD", Limit: 100000, Period: "2024-01"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, ok := idx.Get("acc-1", "2024-01")
	if !ok {
		t.Fatal("expected budget to be found")
	}
	if b.Limit != 100000 {
		t.Errorf("expected limit 100000, got %d", b.Limit)
	}
}

func TestBudgetIndex_Set_Invalid(t *testing.T) {
	idx := NewBudgetIndex()
	if err := idx.Set(Budget{AccountID: "", Currency: "USD", Limit: 100, Period: "2024-01"}); err == nil {
		t.Error("expected error for empty accountID")
	}
	if err := idx.Set(Budget{AccountID: "acc-1", Currency: "USD", Limit: 0, Period: "2024-01"}); err == nil {
		t.Error("expected error for zero limit")
	}
	if err := idx.Set(Budget{AccountID: "acc-1", Currency: "", Limit: 100, Period: "2024-01"}); err == nil {
		t.Error("expected error for empty currency")
	}
	if err := idx.Set(Budget{AccountID: "acc-1", Currency: "USD", Limit: 100, Period: ""}); err == nil {
		t.Error("expected error for empty period")
	}
}

func TestBudgetIndex_Get_NotFound(t *testing.T) {
	idx := NewBudgetIndex()
	_, ok := idx.Get("missing", "2024-01")
	if ok {
		t.Error("expected not found")
	}
}

func TestBudgetIndex_Evaluate_OK(t *testing.T) {
	idx := NewBudgetIndex()
	_ = idx.Set(Budget{AccountID: "acc-1", Currency: "USD", Limit: 50000, Period: "2024-01"})
	status, err := idx.Evaluate("acc-1", "2024-01", 30000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.Exceeded {
		t.Error("expected budget not exceeded")
	}
	if status.Remaining != 20000 {
		t.Errorf("expected remaining 20000, got %d", status.Remaining)
	}
}

func TestBudgetIndex_Evaluate_Exceeded(t *testing.T) {
	idx := NewBudgetIndex()
	_ = idx.Set(Budget{AccountID: "acc-1", Currency: "USD", Limit: 10000, Period: "2024-01"})
	status, err := idx.Evaluate("acc-1", "2024-01", 15000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !status.Exceeded {
		t.Error("expected budget to be exceeded")
	}
}

func TestBudgetIndex_Evaluate_NotFound(t *testing.T) {
	idx := NewBudgetIndex()
	_, err := idx.Evaluate("ghost", "2024-01", 0)
	if err == nil {
		t.Error("expected error for missing budget")
	}
}
