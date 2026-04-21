package ledger

import (
	"strings"
	"testing"
)

func makeBudgetBook(t *testing.T) *Book {
	t.Helper()
	book := NewBook("budget-test")
	acct, err := NewAccount("exp-1", "Marketing", AccountTypeExpense)
	if err != nil {
		t.Fatalf("NewAccount: %v", err)
	}
	if err := book.AddAccount(acct); err != nil {
		t.Fatalf("AddAccount: %v", err)
	}
	return book
}

func TestCheckBudgets_NoViolation(t *testing.T) {
	book := makeBudgetBook(t)
	idx := NewBudgetIndex()
	_ = idx.Set(Budget{AccountID: "exp-1", Currency: "USD", Limit: 500000, Period: "2024-01"})

	result, err := CheckBudgets(idx, book, "2024-01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HasViolations() {
		t.Error("expected no violations")
	}
}

func TestCheckBudgets_NilIndex(t *testing.T) {
	book := makeBudgetBook(t)
	_, err := CheckBudgets(nil, book, "2024-01")
	if err == nil {
		t.Error("expected error for nil index")
	}
}

func TestCheckBudgets_NilBook(t *testing.T) {
	idx := NewBudgetIndex()
	_, err := CheckBudgets(idx, nil, "2024-01")
	if err == nil {
		t.Error("expected error for nil book")
	}
}

func TestCheckBudgets_PeriodFilter(t *testing.T) {
	book := makeBudgetBook(t)
	idx := NewBudgetIndex()
	_ = idx.Set(Budget{AccountID: "exp-1", Currency: "USD", Limit: 100, Period: "2024-02"})

	result, err := CheckBudgets(idx, book, "2024-01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Statuses) != 0 {
		t.Errorf("expected 0 statuses for different period, got %d", len(result.Statuses))
	}
}

func TestBudgetCheckResult_String(t *testing.T) {
	result := &BudgetCheckResult{
		Statuses: []*BudgetStatus{
			{Budget: Budget{AccountID: "acc-1", Currency: "USD", Limit: 10000, Period: "2024-01"}, Spent: 12000, Remaining: -2000, Exceeded: true},
		},
		Violations: []*BudgetStatus{},
	}
	result.Violations = result.Statuses
	out := result.String()
	if !strings.Contains(out, "EXCEEDED") {
		t.Errorf("expected EXCEEDED in output, got: %s", out)
	}
}

func TestBudgetCheckResult_String_Empty(t *testing.T) {
	result := &BudgetCheckResult{}
	out := result.String()
	if !strings.Contains(out, "no budgets") {
		t.Errorf("expected 'no budgets' in output, got: %s", out)
	}
}
