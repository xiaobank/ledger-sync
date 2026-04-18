package ledger

import (
	"testing"
)

func makeBook(t *testing.T) (*Book, *Account, *Account) {
	t.Helper()
	b := NewBook()
	cash, _ := NewAccount("cash", "Cash", Asset)
	revenue, _ := NewAccount("revenue", "Revenue", Revenue)
	_ = b.AddAccount(cash)
	_ = b.AddAccount(revenue)
	return b, cash, revenue
}

func TestBook_AddAccount_Duplicate(t *testing.T) {
	b := NewBook()
	a, _ := NewAccount("acc", "Test", Asset)
	_ = b.AddAccount(a)
	if err := b.AddAccount(a); err == nil {
		t.Fatal("expected error for duplicate account")
	}
}

func TestBook_GetAccount_NotFound(t *testing.T) {
	b := NewBook()
	_, err := b.GetAccount("missing")
	if err == nil {
		t.Fatal("expected error for missing account")
	}
}

func TestBook_Post_Valid(t *testing.T) {
	b, cash, revenue := makeBook(t)
	tx, err := NewTransaction("tx-1", []Entry{
		{AccountID: "cash", Type: Debit, Amount: 1000, Currency: "USD"},
		{AccountID: "revenue", Type: Credit, Amount: 1000, Currency: "USD"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := b.Post(tx); err != nil {
		t.Fatalf("post failed: %v", err)
	}
	if cash.Balance("USD") != 1000 {
		t.Fatalf("expected cash 1000, got %d", cash.Balance("USD"))
	}
	if revenue.Balance("USD") != -1000 {
		t.Fatalf("expected revenue -1000, got %d", revenue.Balance("USD"))
	}
}

func TestBook_Post_MissingAccount(t *testing.T) {
	b, _, _ := makeBook(t)
	tx, _ := NewTransaction("tx-2", []Entry{
		{AccountID: "ghost", Type: Debit, Amount: 500, Currency: "USD"},
		{AccountID: "cash", Type: Credit, Amount: 500, Currency: "USD"},
	})
	if err := b.Post(tx); err == nil {
		t.Fatal("expected error for missing account")
	}
}
