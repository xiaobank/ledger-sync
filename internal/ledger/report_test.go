package ledger

import (
	"strings"
	"testing"
)

func TestGenerateBalanceReport_Valid(t *testing.T) {
	b := makeBook(t)

	tx, err := NewTransaction("tx1", []Leg{
		{AccountID: "acc-cash", Type: Debit, Amount: 1000, Currency: "USD"},
		{AccountID: "acc-revenue", Type: Credit, Amount: 1000, Currency: "USD"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := b.Post(tx); err != nil {
		t.Fatalf("unexpected error posting: %v", err)
	}

	report, err := b.GenerateBalanceReport()
	if err != nil {
		t.Fatalf("unexpected error generating report: %v", err)
	}
	if len(report.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(report.Entries))
	}

	for _, e := range report.Entries {
		switch e.AccountID {
		case "acc-cash":
			if e.Debit != 1000 || e.Credit != 0 || e.Net != 1000 {
				t.Errorf("unexpected cash entry: %+v", e)
			}
		case "acc-revenue":
			if e.Credit != 1000 || e.Debit != 0 || e.Net != -1000 {
				t.Errorf("unexpected revenue entry: %+v", e)
			}
		default:
			t.Errorf("unexpected account in report: %s", e.AccountID)
		}
	}
}

func TestGenerateBalanceReport_Empty(t *testing.T) {
	b := makeBook(t)
	report, err := b.GenerateBalanceReport()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(report.Entries) != 0 {
		t.Errorf("expected empty report, got %d entries", len(report.Entries))
	}
}

func TestBalanceReport_String(t *testing.T) {
	r := &BalanceReport{
		Entries: []BalanceEntry{
			{AccountID: "acc-cash", AccountName: "Cash", Currency: "USD", Debit: 500, Credit: 0, Net: 500},
		},
	}
	s := r.String()
	if s == "" {
		t.Error("expected non-empty report string")
	}
	if !strings.Contains(s, "acc-cash") {
		t.Errorf("expected report string to contain account ID 'acc-cash', got: %s", s)
	}
	if !strings.Contains(s, "Cash") {
		t.Errorf("expected report string to contain account name 'Cash', got: %s", s)
	}
}

func TestBalanceReport_String_Empty(t *testing.T) {
	r := &BalanceReport{
		Entries: []BalanceEntry{},
	}
	s := r.String()
	if s == "" {
		t.Error("expected non-empty report string even for empty entries")
	}
	if strings.Contains(s, "acc-") {
		t.Errorf("expected no account entries in empty report string, got: %s", s)
	}
}
