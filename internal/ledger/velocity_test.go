package ledger

import (
	"testing"
	"time"
)

func makeVelocityBook(t *testing.T) *Book {
	t.Helper()
	b := NewBook("vbook")
	acc, _ := NewAccount("acc-1", "Checking", AccountTypeAsset)
	cash, _ := NewAccount("cash", "Cash", AccountTypeAsset)
	_ = b.AddAccount(acc)
	_ = b.AddAccount(cash)
	return b
}

func TestVelocityIndex_Add_Valid(t *testing.T) {
	vi := NewVelocityIndex()
	err := vi.Add(VelocityRule{
		AccountID: "acc-1",
		Window:    24 * time.Hour,
		MaxCount:  5,
		Currency:  "USD",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVelocityIndex_Add_Invalid(t *testing.T) {
	vi := NewVelocityIndex()
	cases := []VelocityRule{
		{AccountID: "", Window: time.Hour, MaxCount: 1, Currency: "USD"},
		{AccountID: "a", Window: 0, MaxCount: 1, Currency: "USD"},
		{AccountID: "a", Window: time.Hour, MaxCount: 0, MaxAmount: 0, Currency: "USD"},
		{AccountID: "a", Window: time.Hour, MaxCount: -1, Currency: "USD"},
	}
	for _, r := range cases {
		if err := vi.Add(r); err == nil {
			t.Errorf("expected error for rule %+v", r)
		}
	}
}

func TestCheckVelocity_NoViolation(t *testing.T) {
	b := makeVelocityBook(t)
	now := time.Now()

	tx, _ := NewTransaction("tx1", now, []Leg{
		{AccountID: "acc-1", Amount: 100, Currency: "USD", Type: EntryTypeDebit},
		{AccountID: "cash", Amount: 100, Currency: "USD", Type: EntryTypeCredit},
	})
	_ = b.Post(tx)

	vi := NewVelocityIndex()
	_ = vi.Add(VelocityRule{AccountID: "acc-1", Window: time.Hour, MaxCount: 5, Currency: "USD"})

	violations, err := CheckVelocity(vi, b, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if HasVelocityViolations(violations) {
		t.Errorf("expected no violations, got %v", violations)
	}
}

func TestCheckVelocity_CountViolation(t *testing.T) {
	b := makeVelocityBook(t)
	now := time.Now()

	for i := 0; i < 4; i++ {
		tx, _ := NewTransaction(
			string(rune('a'+i)),
			now,
			[]Leg{
				{AccountID: "acc-1", Amount: 10, Currency: "USD", Type: EntryTypeDebit},
				{AccountID: "cash", Amount: 10, Currency: "USD", Type: EntryTypeCredit},
			},
		)
		_ = b.Post(tx)
	}

	vi := NewVelocityIndex()
	_ = vi.Add(VelocityRule{AccountID: "acc-1", Window: time.Hour, MaxCount: 2, Currency: "USD"})

	violations, err := CheckVelocity(vi, b, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !HasVelocityViolations(violations) {
		t.Error("expected count violation")
	}
}

func TestCheckVelocity_AmountViolation(t *testing.T) {
	b := makeVelocityBook(t)
	now := time.Now()

	tx, _ := NewTransaction("tx1", now, []Leg{
		{AccountID: "acc-1", Amount: 500, Currency: "USD", Type: EntryTypeDebit},
		{AccountID: "cash", Amount: 500, Currency: "USD", Type: EntryTypeCredit},
	})
	_ = b.Post(tx)

	vi := NewVelocityIndex()
	_ = vi.Add(VelocityRule{AccountID: "acc-1", Window: time.Hour, MaxAmount: 100, Currency: "USD"})

	violations, err := CheckVelocity(vi, b, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !HasVelocityViolations(violations) {
		t.Error("expected amount violation")
	}
	if violations[0].String() == "" {
		t.Error("expected non-empty violation string")
	}
}

func TestCheckVelocity_NilIndex(t *testing.T) {
	b := makeVelocityBook(t)
	_, err := CheckVelocity(nil, b, time.Now())
	if err == nil {
		t.Error("expected error for nil index")
	}
}

func TestCheckVelocity_NilBook(t *testing.T) {
	vi := NewVelocityIndex()
	_, err := CheckVelocity(vi, nil, time.Now())
	if err == nil {
		t.Error("expected error for nil book")
	}
}

func TestCheckVelocity_OutsideWindow(t *testing.T) {
	b := makeVelocityBook(t)
	now := time.Now()
	old := now.Add(-48 * time.Hour)

	tx, _ := NewTransaction("tx-old", old, []Leg{
		{AccountID: "acc-1", Amount: 999, Currency: "USD", Type: EntryTypeDebit},
		{AccountID: "cash", Amount: 999, Currency: "USD", Type: EntryTypeCredit},
	})
	_ = b.Post(tx)

	vi := NewVelocityIndex()
	_ = vi.Add(VelocityRule{AccountID: "acc-1", Window: time.Hour, MaxAmount: 100, Currency: "USD"})

	violations, err := CheckVelocity(vi, b, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if HasVelocityViolations(violations) {
		t.Error("expected no violations for out-of-window transaction")
	}
}
