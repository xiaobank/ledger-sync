package ledger

import (
	"testing"
)

func makeMarginBook(t *testing.T) *Book {
	t.Helper()
	b := makeBook(t)

	assets, err := NewAccount("assets", "Assets", AccountTypeAsset)
	if err != nil {
		t.Fatalf("NewAccount: %v", err)
	}
	liabilities, err := NewAccount("liabilities", "Liabilities", AccountTypeLiability)
	if err != nil {
		t.Fatalf("NewAccount: %v", err)
	}

	if err := b.AddAccount(assets); err != nil {
		t.Fatalf("AddAccount assets: %v", err)
	}
	if err := b.AddAccount(liabilities); err != nil {
		t.Fatalf("AddAccount liabilities: %v", err)
	}

	// assets = 300, liabilities = 100  => ratio 3.0
	tx, _ := NewTransaction("tx-margin-1", []Entry{
		{AccountID: "assets", Amount: 300, Currency: "USD", Type: EntryTypeDebit},
		{AccountID: "liabilities", Amount: 300, Currency: "USD", Type: EntryTypeCredit},
	})
	if err := b.Post(tx); err != nil {
		t.Fatalf("Post: %v", err)
	}
	return b
}

func TestMarginIndex_Add_Valid(t *testing.T) {
	idx := NewMarginIndex()
	err := idx.Add(MarginRule{ID: "r1", Numerator: "assets", Denominator: "liabilities", MinRatio: 1.5, Currency: "USD"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(idx.Rules()) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(idx.Rules()))
	}
}

func TestMarginIndex_Add_Invalid(t *testing.T) {
	idx := NewMarginIndex()
	cases := []MarginRule{
		{ID: "", Numerator: "a", Denominator: "b", MinRatio: 1.0},
		{ID: "r1", Numerator: "", Denominator: "b", MinRatio: 1.0},
		{ID: "r2", Numerator: "a", Denominator: "", MinRatio: 1.0},
		{ID: "r3", Numerator: "a", Denominator: "b", MinRatio: 0},
	}
	for _, c := range cases {
		if err := idx.Add(c); err == nil {
			t.Errorf("expected error for rule %+v", c)
		}
	}
}

func TestMarginIndex_Add_Duplicate(t *testing.T) {
	idx := NewMarginIndex()
	rule := MarginRule{ID: "r1", Numerator: "a", Denominator: "b", MinRatio: 1.0, Currency: "USD"}
	_ = idx.Add(rule)
	if err := idx.Add(rule); err == nil {
		t.Fatal("expected error for duplicate rule")
	}
}

func TestCheckMargins_NoViolation(t *testing.T) {
	b := makeMarginBook(t)
	idx := NewMarginIndex()
	_ = idx.Add(MarginRule{ID: "r1", Numerator: "assets", Denominator: "liabilities", MinRatio: 2.0, Currency: "USD"})

	violations, err := CheckMargins(idx, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if HasMarginViolations(violations) {
		t.Errorf("expected no violations, got %d", len(violations))
	}
}

func TestCheckMargins_Violation(t *testing.T) {
	b := makeMarginBook(t)
	idx := NewMarginIndex()
	// ratio is 3.0; require 5.0 => violation
	_ = idx.Add(MarginRule{ID: "r1", Numerator: "assets", Denominator: "liabilities", MinRatio: 5.0, Currency: "USD"})

	violations, err := CheckMargins(idx, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !HasMarginViolations(violations) {
		t.Fatal("expected a violation")
	}
	if violations[0].Rule.ID != "r1" {
		t.Errorf("expected rule r1, got %s", violations[0].Rule.ID)
	}
}

func TestCheckMargins_NilBook(t *testing.T) {
	idx := NewMarginIndex()
	_ = idx.Add(MarginRule{ID: "r1", Numerator: "assets", Denominator: "liabilities", MinRatio: 1.5, Currency: "USD"})

	_, err := CheckMargins(idx, nil)
	if err == nil {
		t.Fatal("expected error when book is nil")
	}
}
