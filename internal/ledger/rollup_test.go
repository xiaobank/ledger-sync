package ledger

import (
	"testing"
	"time"
)

func makeRollupBook(t *testing.T) *Book {
	t.Helper()
	b := makeBook(t)

	leg1 := Leg{AccountID: "acc-1", Amount: 5000, Currency: "USD", Type: LegCredit}
	leg2 := Leg{AccountID: "acc-2", Amount: 5000, Currency: "USD", Type: LegDebit}

	tx1 := Transaction{
		ID:       "tx-1",
		PostedAt: time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		Legs:     []Leg{leg1, leg2},
	}
	tx2 := Transaction{
		ID:       "tx-2",
		PostedAt: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
		Legs:     []Leg{{AccountID: "acc-1", Amount: 3000, Currency: "USD", Type: LegCredit}, {AccountID: "acc-2", Amount: 3000, Currency: "USD", Type: LegDebit}},
	}
	tx3 := Transaction{
		ID:       "tx-3",
		PostedAt: time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC),
		Legs:     []Leg{{AccountID: "acc-1", Amount: 2000, Currency: "EUR", Type: LegCredit}, {AccountID: "acc-2", Amount: 2000, Currency: "EUR", Type: LegDebit}},
	}

	if err := b.Post(tx1); err != nil {
		t.Fatalf("post tx1: %v", err)
	}
	if err := b.Post(tx2); err != nil {
		t.Fatalf("post tx2: %v", err)
	}
	if err := b.Post(tx3); err != nil {
		t.Fatalf("post tx3: %v", err)
	}
	return b
}

func TestRollup_Monthly(t *testing.T) {
	b := makeRollupBook(t)
	res, err := RollupTransactions(b, "acc-1", RollupMonthly)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.AccountID != "acc-1" {
		t.Errorf("expected acc-1, got %s", res.AccountID)
	}
	if len(res.Buckets) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(res.Buckets))
	}
	if res.Buckets[0].Label != "2024-01" || res.Buckets[0].Total != 8000 {
		t.Errorf("bucket 0: got label=%s total=%d", res.Buckets[0].Label, res.Buckets[0].Total)
	}
	if res.Buckets[1].Label != "2024-02" || res.Buckets[1].Currency != "EUR" {
		t.Errorf("bucket 1: got label=%s currency=%s", res.Buckets[1].Label, res.Buckets[1].Currency)
	}
}

func TestRollup_Daily(t *testing.T) {
	b := makeRollupBook(t)
	res, err := RollupTransactions(b, "acc-1", RollupDaily)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Buckets) != 3 {
		t.Fatalf("expected 3 daily buckets, got %d", len(res.Buckets))
	}
}

func TestRollup_NilBook(t *testing.T) {
	_, err := RollupTransactions(nil, "acc-1", RollupMonthly)
	if err == nil {
		t.Error("expected error for nil book")
	}
}

func TestRollup_EmptyAccountID(t *testing.T) {
	b := makeRollupBook(t)
	_, err := RollupTransactions(b, "", RollupMonthly)
	if err == nil {
		t.Error("expected error for empty accountID")
	}
}

func TestRollup_UnknownPeriod(t *testing.T) {
	b := makeRollupBook(t)
	_, err := RollupTransactions(b, "acc-1", RollupPeriod("yearly"))
	if err == nil {
		t.Error("expected error for unknown period")
	}
}
