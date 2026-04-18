package ledger

import (
	"testing"
)

func TestNewAccount_Valid(t *testing.T) {
	a, err := NewAccount("acc-1", "Cash", Asset)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.ID != "acc-1" || a.Name != "Cash" || a.Type != Asset {
		t.Fatal("account fields not set correctly")
	}
}

func TestNewAccount_EmptyID(t *testing.T) {
	_, err := NewAccount("", "Cash", Asset)
	if err == nil {
		t.Fatal("expected error for empty id")
	}
}

func TestNewAccount_EmptyName(t *testing.T) {
	_, err := NewAccount("acc-1", "", Asset)
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestNewAccount_InvalidType(t *testing.T) {
	_, err := NewAccount("acc-1", "Cash", AccountType("unknown"))
	if err == nil {
		t.Fatal("expected error for invalid type")
	}
}

func TestAccount_Balance(t *testing.T) {
	a, _ := NewAccount("acc-1", "Cash", Asset)
	if a.Balance("USD") != 0 {
		t.Fatal("expected zero balance")
	}
	a.apply("USD", 500)
	if a.Balance("USD") != 500 {
		t.Fatalf("expected 500, got %d", a.Balance("USD"))
	}
	a.apply("USD", -200)
	if a.Balance("USD") != 300 {
		t.Fatalf("expected 300, got %d", a.Balance("USD"))
	}
}

func TestAccount_MultiCurrency(t *testing.T) {
	a, _ := NewAccount("acc-2", "Bank", Asset)
	a.apply("USD", 1000)
	a.apply("EUR", 800)
	if a.Balance("USD") != 1000 || a.Balance("EUR") != 800 {
		t.Fatal("multi-currency balances incorrect")
	}
}
