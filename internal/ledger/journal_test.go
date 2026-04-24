package ledger

import (
	"testing"
	"time"
)

func makeJournalBook(t *testing.T) *Book {
	t.Helper()
	b := makeBook(t)
	_ = b.AddAccount(NewAccount("acc-a", "Account A", AccountAsset))
	_ = b.AddAccount(NewAccount("acc-b", "Account B", AccountLiability))

	tx, err := NewTransaction("tx-1", "first transfer",
		[]Leg{{AccountID: "acc-a", Amount: 500, Currency: "USD", Type: EntryDebit},
			{AccountID: "acc-b", Amount: 500, Currency: "USD", Type: EntryCredit}},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tx.Timestamp = time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	_ = b.Post(tx)
	return b
}

func TestBuildJournal_Valid(t *testing.T) {
	b := makeJournalBook(t)
	j, err := BuildJournal(b)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(j.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(j.Entries))
	}
}

func TestBuildJournal_NilBook(t *testing.T) {
	_, err := BuildJournal(nil)
	if err == nil {
		t.Fatal("expected error for nil book")
	}
}

func TestJournal_FilterByAccount(t *testing.T) {
	b := makeJournalBook(t)
	j, _ := BuildJournal(b)
	filtered := j.FilterByAccount("acc-a")
	if len(filtered.Entries) != 1 {
		t.Fatalf("expected 1 entry for acc-a, got %d", len(filtered.Entries))
	}
	if filtered.Entries[0].Account != "acc-a" {
		t.Errorf("unexpected account: %s", filtered.Entries[0].Account)
	}
}

func TestJournal_FilterByDateRange(t *testing.T) {
	b := makeJournalBook(t)
	j, _ := BuildJournal(b)

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	filtered := j.FilterByDateRange(from, to)
	if len(filtered.Entries) != 2 {
		t.Fatalf("expected 2 entries in range, got %d", len(filtered.Entries))
	}

	outside := j.FilterByDateRange(
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
	)
	if len(outside.Entries) != 0 {
		t.Fatalf("expected 0 entries outside range, got %d", len(outside.Entries))
	}
}

func TestJournal_DebitCreditAssignment(t *testing.T) {
	b := makeJournalBook(t)
	j, _ := BuildJournal(b)

	for _, e := range j.Entries {
		if e.Account == "acc-a" && e.Debit != 500 {
			t.Errorf("acc-a expected debit 500, got %d", e.Debit)
		}
		if e.Account == "acc-b" && e.Credit != 500 {
			t.Errorf("acc-b expected credit 500, got %d", e.Credit)
		}
	}
}
