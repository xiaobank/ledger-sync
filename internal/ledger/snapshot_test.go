package ledger

import (
	"testing"
)

func makePopulatedBook(t *testing.T) *Book {
	t.Helper()
	b, _ := NewBook("test-book")

	cash, _ := NewAccount("cash-01", "Cash", AccountTypeAsset)
	rev, _ := NewAccount("rev-01", "Revenue", AccountTypeRevenue)
	_ = b.AddAccount(cash)
	_ = b.AddAccount(rev)

	tx, _ := NewTransaction([]Entry{
		{AccountID: "cash-01", Amount: 500, Currency: "USD", Type: EntryTypeDebit},
		{AccountID: "rev-01", Amount: 500, Currency: "USD", Type: EntryTypeCredit},
	})
	_ = b.Post(tx)
	return b
}

func TestTakeSnapshot_Valid(t *testing.T) {
	b := makePopulatedBook(t)
	snap, err := TakeSnapshot(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.BookName != "test-book" {
		t.Errorf("expected book name 'test-book', got %q", snap.BookName)
	}
	if snap.Balances["cash-01"]["USD"] != 500 {
		t.Errorf("expected cash balance 500, got %v", snap.Balances["cash-01"]["USD"])
	}
	if snap.Balances["rev-01"]["USD"] != 500 {
		t.Errorf("expected revenue balance 500, got %v", snap.Balances["rev-01"]["USD"])
	}
}

func TestTakeSnapshot_NilBook(t *testing.T) {
	_, err := TakeSnapshot(nil)
	if err == nil {
		t.Fatal("expected error for nil book")
	}
}

func TestDiffSnapshots_ChangedBalance(t *testing.T) {
	b := makePopulatedBook(t)

	before, _ := TakeSnapshot(b)

	tx, _ := NewTransaction([]Entry{
		{AccountID: "cash-01", Amount: 200, Currency: "USD", Type: EntryTypeDebit},
		{AccountID: "rev-01", Amount: 200, Currency: "USD", Type: EntryTypeCredit},
	})
	_ = b.Post(tx)

	after, _ := TakeSnapshot(b)

	diff := DiffSnapshots(before, after)
	if _, ok := diff["cash-01"]; !ok {
		t.Fatal("expected diff entry for cash-01")
	}
	pair := diff["cash-01"]["USD"]
	if pair[0] != 500 || pair[1] != 700 {
		t.Errorf("expected [500, 700], got %v", pair)
	}
}

func TestDiffSnapshots_NoChange(t *testing.T) {
	b := makePopulatedBook(t)
	snap, _ := TakeSnapshot(b)
	diff := DiffSnapshots(snap, snap)
	if len(diff) != 0 {
		t.Errorf("expected empty diff, got %v", diff)
	}
}
